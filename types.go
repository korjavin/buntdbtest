package main

import (
	"time"
)

type Order struct {
	Oid      int
	Pid      int
	Account  string
	Amount   float32
	Statuses []Orderstatus
	Request  Request
	Bill     []Bill
}
type Orderstatus struct {
	Sid   int
	Regdt time.Time
}
type Request struct {
	Command  string
	Account  string
	Txn_id   string
	Md5      string
	Txn_date string
	Sum      float32
	Ka       string
	Regdt    time.Time
}
type Bill struct {
	Req   string
	Res   string
	Regdt time.Time
}
