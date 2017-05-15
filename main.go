package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	_ "net/http/pprof"

	"github.com/tidwall/buntdb"
	"github.com/tidwall/gjson"
	// "github.com/tidwall/gjson"
)

var (
	bdb        *buntdb.DB
	oneTimeGet int
	wg         sync.WaitGroup
)

func writer(number int, ochan <-chan Order) {
	for {
		o := <-ochan
		// log.Printf("Got %v", o)
		buf, err := json.Marshal(o)
		if err != nil {
			panic(err)
		}
		// log.Println(string(buf))
		err = bdb.Update(func(tx *buntdb.Tx) error {
			_, _, err := tx.Set(strconv.Itoa(o.Oid), string(buf), nil)
			return err
		})
		checkErr(err)
		// fmt.Printf("%d", number)
	}
}
func init() {
	var err error
	bdb, err = buntdb.Open("my.db")
	checkErr(err)
}
func compareDatesBuilder(path string) func(string, string) bool {
	return func(a, b string) bool {
		astr := gjson.Get(a, path).String()
		bstr := gjson.Get(b, path).String()
		start, err := time.Parse("2006-01-02T15:04:05Z", astr)
		if err != nil {
			return true
		}
		end, err := time.Parse("2006-01-02T15:04:05Z", bstr)
		if err != nil {
			return false
		}
		return end.After(start)
	}

}
func compareDates(a, b string) bool {
	astr := gjson.Get(a, `Statuses.#[Sid=6].Regdt`).String()
	bstr := gjson.Get(b, `Statuses.#[Sid=6].Regdt`).String()
	start, err := time.Parse("2006-01-02T15:04:05Z", astr)
	if err != nil {
		return true
	}
	end, err := time.Parse("2006-01-02T15:04:05Z", bstr)
	if err != nil {
		return false
	}
	return end.After(start)

}
func main() {
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()
	oneTimeGet = 5000

	// copy()
	err := bdb.Update(func(tx *buntdb.Tx) error {
		tx.CreateIndex("sid6", "*", compareDatesBuilder(`Statuses.#[Sid=6].Regdt`))
		tx.CreateIndex("sid3", "*", compareDatesBuilder(`Statuses.#[Sid=3].Regdt`))
		tx.CreateIndex("sid1", "*", compareDatesBuilder(`Statuses.#[Sid=1].Regdt`))
		return nil
	})
	checkErr(err)
	bdb.View(func(tx *buntdb.Tx) error {
		tx.AscendGreaterOrEqual("sid1", `{"Statuses":[{"Sid":1,"Regdt":"2017-05-10T05:04:00Z"}]}`, func(key, value string) bool {
			var order Order
			err := json.Unmarshal([]byte(value), &order)
			if err != nil {
				fmt.Println("error:", err)
			}

			regdtstr := gjson.Get(value, `Statuses.#[Sid=6].Regdt`)
			fmt.Println("time:", regdtstr)
			_, err = time.Parse("2006-01-02T15:04:05Z", regdtstr.String())
			if err != nil {
				return true
			}

			fmt.Printf("key=%s, o=%v \n", key, value)
			return false
		})
		return nil
	})
}

func copy() {
	ochan := make(chan Order, oneTimeGet)
	min, max := getCount()
	log.Println(min, max)
	for i := min; i < max; i += oneTimeGet {
		wg.Add(1)
		go getData(ochan, i)
	}

	for i := 0; i < 10; i++ {
		go writer(i, ochan)
	}
	time.Sleep(1 * time.Second)
	wg.Wait()
}
