package userbase

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

func init() {
	var err error
	db, err = bolt.Open("users.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func AddUser(chatID int64) error {
	key := []byte(strconv.FormatInt(chatID, 10))
	return db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte("users")).Put(key, []byte{'1'})
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func RemoveUser(chatID int64) error {
	key := []byte(strconv.FormatInt(chatID, 10))
	return db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte("users")).Delete(key)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func AllUsers() ([]int64, error) {
	var res []int64
	err := db.View(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("users")).ForEach(func(k []byte, v []byte) error {
			if v[0] == '1' {
				id, err := strconv.ParseInt(string(k), 10, 64)
				if err != nil {
					return nil
				}
				res = append(res, id)
			}
			return nil
		})
	})
	if err != nil {
		return []int64{}, err
	}
	return res, nil
}
