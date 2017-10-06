package users

import (
	"networks2016-task1/configs"
)

type usersDbInterface interface {
	AddUser(login, password string) (*User, error)
	CheckExistance(login, password string) (bool, error)
	GetUser(login, password string) (*User, error)
}

const method = SHA256

func InitAuth(postgresConfigPath string) {
	postgresConfig := configs.ReadPostgresConfig(postgresConfigPath)
	salt := postgresConfig.Salt
	initHash(method, 32, salt)
	initPostgres(postgresConfig)
}

type User struct {
	Login             string
	passwordHashBytes []byte
}

func NewUser(login, password string) *User {
	return &User{
		Login:             login,
		passwordHashBytes: getHash(password),
	}
}

type Users struct {
	Db	usersDbInterface
}

func NewUsers() (*Users, error) {
	return &Users{
		Db: newUsersInPostgres(),
		//Db: newUsersInMemoryDb(),
	}, nil
}
