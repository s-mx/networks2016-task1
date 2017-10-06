package users

import (
	"crypto/sha1"
	"crypto/sha256"
	"golang.org/x/crypto/pbkdf2"
	"log"
)

const (
	SHA256 = iota
	PBKDF2
)

var SALT_STRING string
var SALT_BYTES	[]byte
var PBKDF2_LENGTH int

func sha256HashDriver(password string) []byte {
	sum := sha256.Sum256([]byte(password + SALT_STRING))
	return sum[:]
}

func pbkdf2HashDriver(password string) []byte {
	passwordBytes := []byte(password)
	return pbkdf2.Key(passwordBytes, SALT_BYTES, 4096, PBKDF2_LENGTH, sha1.New)
}

var hashFunc func(password string) []byte

func initHash(method int, pbkdf2_length int, salt string) {
	if method == SHA256 {
		hashFunc = sha256HashDriver
	} else if method == PBKDF2 {
		hashFunc = pbkdf2HashDriver
	} else {
		log.Fatal("Unknown hash method")
	}

	if pbkdf2_length > 0 {
		PBKDF2_LENGTH = pbkdf2_length
	} else {
		PBKDF2_LENGTH = 32
	}

	SALT_STRING = salt
	SALT_BYTES = []byte(salt)
}

func getHash(password string) []byte {
	return hashFunc(password)
}

// TODO: add getHashByMethod
