package configs

import (
	"io/ioutil"
	"encoding/json"
	"log"
)

type PostgresConfig struct {
	UserDb		string `json:"userDb"`
	PasswordDb	string `json:"password"`
	Host 		string `json:""`
	NameDb		string `json:""`
	Salt		string `json:""`
}

func ReadPostgresConfig(configName string) PostgresConfig {
	valueBytes, _ := ioutil.ReadFile(configName)

	var config PostgresConfig
	err := json.Unmarshal(valueBytes, &config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}
