package users

import (
	"sync"
	"errors"
	"reflect"
)

type usersInMemoryDb struct {
	table	map[string]*User
	mutex	sync.Mutex
}

func newUsersInMemoryDb() *usersInMemoryDb {
	return &usersInMemoryDb{
		table: make(map[string]*User),
	}
}

func (db *usersInMemoryDb) AddUser(login, password string) (*User, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	_, dup := db.table[login]
	if dup {
		return nil, errors.New("Error: there is a duplicate login in auth")
	}

	user := NewUser(login, password)
	db.table[login] = user
	return user, nil
}

func (db *usersInMemoryDb) CheckExistance(login, password string) (bool, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	user, ok := db.table[login]
	if !ok {
		return false, nil
	}

	return reflect.DeepEqual(user.passwordHashBytes, getHash(password)), nil
}

func (db *usersInMemoryDb) GetUser(login, password string) (*User, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	user, ok := db.table[login]
	if !ok || !reflect.DeepEqual(user.passwordHashBytes, getHash(password)) {
		return nil, errors.New("Wrong password")
	}

	return user, nil
}