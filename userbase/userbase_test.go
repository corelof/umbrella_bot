package userbase

import (
	"fmt"
	"github.com/boltdb/bolt"
	"testing"
)

func clearBucket() {
	db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte("users"))
		if err != nil {
			return fmt.Errorf("delete bucket: %s", err)
		}
		return nil
	})
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func TestAddUser(t *testing.T) {
	clearBucket()
	AddUser(1)
	AddUser(2)
	AddUser(3)
	AddUser(4)
	AddUser(1)
	users, err := AllUsers()
	if len(users) != 4 || err != nil {
		t.Errorf("expected 4 users, found %d, error is %s", len(users), err)
	}
}

func TestRemoveUser(t *testing.T) {
	RemoveUser(1)
	RemoveUser(1)
	RemoveUser(2)
	users, err := AllUsers()
	if len(users) != 2 || err != nil {
		t.Errorf("expected 2 users, found %d, error is %s", len(users), err)
	}
}

func TestAllUsers(t *testing.T) {
	users, err := AllUsers()
	if len(users) != 2 || err != nil {
		t.Errorf("expected 2 users, found %d, error is %s", len(users), err)
	}
	if users[0] > users[1] {
		users[0], users[1] = users[1], users[0]
	}
	if users[0] != 3 || users[1] != 4 {
		t.Errorf("expected 3 and 4, found %d and %d", users[0], users[1])
	}
}
