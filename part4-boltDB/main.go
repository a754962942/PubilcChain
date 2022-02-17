package main

import (
	"fmt"

	"github.com/boltdb/bolt"
)

func main() {
	db, err := bolt.Open("./my.db", 0600, nil)
	if err != nil {
		fmt.Printf("Open db failed ,err:\n", err)
		return
	}
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		//创建bucket
		bucket, _ := tx.CreateBucket([]byte("Bucket1"))
		//设置key-val
		bucket.Put([]byte("name"), []byte("abc"))
		return nil
	})
	//查询数据库数据
	db.View(func(tx *bolt.Tx) error {
		//获取bucket
		bucket := tx.Bucket([]byte("Bucket1"))
		val := bucket.Get([]byte("name"))
		fmt.Println(string(val))
		return nil
	})
}
