package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/buntdb"
	"github.com/tidwall/gjson"
	"log"
	"strconv"
)

var (
	bdb *buntdb.DB
)

func writer(ochan <-chan Order) {
	for {
		o := <-ochan
		buf, err := json.Marshal(o)
		if err != nil {
			panic(err)
		}
		log.Println(string(buf))
		err = bdb.Update(func(tx *buntdb.Tx) error {
			_, _, err := tx.Set(strconv.Itoa(o.Oid), string(buf), nil)
			return err
		})
	}
}
func init() {
	bdb, _ = buntdb.Open("my.db")
}
func main() {

	err := bdb.Update(func(tx *buntdb.Tx) error {
		err := tx.CreateIndex("sid", "*", buntdb.IndexJSON("Status.#[Sid=6].Sid"))
		return err
	})
	checkErr(err)
	bdb.View(func(tx *buntdb.Tx) error {
		tx.AscendGreaterOrEqual("sid", `{"Status.Sid":6}`, func(key, value string) bool {
			var order Order
			err := json.Unmarshal([]byte(value), &order)
			if err != nil {
				fmt.Println("error:", err)
			}
			name := gjson.Get(value, `Statuses.#[Sid=6].Sid`)

			fmt.Printf("key=%s, sid=%s\n", key, name)
			return true
		})
		return nil
	})
}

func copy() {
	ochan := make(chan Order)
	stopchan := make(chan bool)
	go getData(ochan, stopchan)
	go writer(ochan)
	<-stopchan
}
