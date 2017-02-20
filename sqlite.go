package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"time"
)

var (
	db *sql.DB
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./sokol.db")
	checkErr(err)
}

func getData(ochan chan<- Order, stopchan chan<- bool) {
	rows, err := db.Query("select oid,cid,pid,account,amount from orders")
	checkErr(err)
	var oid int
	var cid string
	var pid int
	var account string
	var amount float32

	for rows.Next() {
		err = rows.Scan(&oid, &cid, &pid, &account, &amount)
		checkErr(err)
		order := Order{Oid: oid, Pid: pid, Account: account, Amount: amount, Statuses: []Orderstatus(nil)}
		fillStatuses(&order)
		fillBills(&order)
		fillReq(&order)
		// log.Prinln("send")
		ochan <- order
	}
	stopchan <- true
}
func fillStatuses(o *Order) {
	oid := o.Oid
	st, err := db.Query("select sid,regdt from orderstatus where oid=?", oid)
	checkErr(err)
	var sid int
	var regdt time.Time
	for st.Next() {
		err = st.Scan(&sid, &regdt)
		checkErr(err)
		st1 := Orderstatus{Sid: sid, Regdt: regdt}
		o.Statuses = append(o.Statuses, st1)
	}
}
func fillBills(o *Order) {
	oid := o.Oid
	st, err := db.Query("select req,res,regdt from bills where oid=?", oid)
	checkErr(err)
	var req string
	var res string
	var regdt time.Time
	for st.Next() {
		err = st.Scan(&req, &res, &regdt)
		checkErr(err)
		b1 := Bill{Res: res, Req: req, Regdt: regdt}
		o.Bill = append(o.Bill, b1)
	}
}
func fillReq(o *Order) {
	oid := o.Oid
	st, err := db.Query("select    command, account , txn_id , md5 ,  ka  , regdt from requests where oid=?", oid)
	checkErr(err)
	var command string
	var account string
	var txn_id string
	var md5 string
	var ka string
	var regdt time.Time
	for st.Next() {
		err = st.Scan(&command, &account, &txn_id, &md5, &ka, &regdt)
		checkErr(err)
		r1 := Request{Command: command, Account: account, Txn_id: txn_id, Md5: md5, Txn_date: "", Ka: ka, Regdt: regdt}
		o.Request = r1
	}
}
