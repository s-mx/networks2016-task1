package users

import (
	"database/sql"
	"networks2016-task1/configs"
	"log"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB = nil

func getDataSourceName(postgresConfig configs.PostgresConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s", postgresConfig.UserDb,
														 postgresConfig.PasswordDb,
														 postgresConfig.Host,
														 postgresConfig.NameDb)
}

func initPostgres(postgresConfig configs.PostgresConfig) {
	var err error
	db, err = sql.Open("postgres", getDataSourceName(postgresConfig))
	if err != nil {
		log.Fatal(err)
	}
}

type usersInPostgres struct {
}

func newUsersInPostgres() *usersInPostgres {
	return &usersInPostgres{
	}
}

func (usersInPostgres) AddUser(login, password string) (*User, error) {
	user := NewUser(login, password) // TODO: проверить SQL инъекции
	_, err := db.Exec("INSERT INTO users_table(login, password_hash, method_hash) VALUES($1, $2, $3)",
							user.Login, user.passwordHashBytes, method)

	if err != nil {
		// TODO: handle it
		log.Fatal(err)
	}

	return user, nil
}

func (usersInPostgres) CheckExistance(login, password string) (bool, error) {
	user := NewUser(login, password)
	user2 := &User{}
	err := db.QueryRow("SELECT login, password_hash FROM users_table WHERE login = $1 AND password_hash = $2",
						user.Login, user.passwordHashBytes).Scan(&user2.Login, &user2.passwordHashBytes)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		// TODO: handle it
		log.Fatal(err)
	default:
		log.Println("User with login: ", user2.Login, " password: ", user2.passwordHashBytes)
	}

	return true, nil
}

func (usersInPostgres) GetUser(login, password string) (*User, error) {
	user := NewUser(login, password)
	user2 := &User{}
	err := db.QueryRow("SELECT login, password_hash FROM users_table WHERE login = $1 AND password_hash = $2",
		user.Login, user.passwordHashBytes).Scan(&user2.Login, &user2.passwordHashBytes)

		switch {
	case err == sql.ErrNoRows:
		return nil, err
	case err != nil:
		// TODO: handle it
		log.Fatal(err)
	default:
		log.Println("User with login: ", user2.Login, " password: ", user2.passwordHashBytes)
	}

	return user2, nil
}
