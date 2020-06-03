package main

import (
	"fmt"
	"github.com/zusux/go-mysql-db/Db"
	"math/rand"
	"time"
)

func main()  {
	Db.Connect("127.0.0.1",3306,"zuusx","root","123456","","utf8mb4",100,10)
	for i := 0;i<100000;i++{
		go dotest(i)
	}
	for {
		time.Sleep(1 * time.Second)
	}
}

func dotest(l int)  {

	db := Db.NewDb()
	for i:=1;i<100;i++{
		db.Debug(true).Table("book").Where("book_id","=",i).Select()
		db.Debug(true).Table("book").Where("book_name","=","A").Select()
		fmt.Println(l,",",i)
	}
}


func  GetRandomString(l int) string {
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}