package main

import (
	"encoding/json"
	"github.com/tidwall/buntdb"
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

func main() {
	ochan := make(chan Order)
	stopchan := make(chan bool)
	bdb, _ = buntdb.Open("my.db")
	go getData(ochan, stopchan)
	go writer(ochan)
	<-stopchan
}
