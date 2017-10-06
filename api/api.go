package api

import (
	"net/http"
	"fmt"
	"networks2016-task1/session"
	"strconv"
	"github.com/gorilla/mux"
)

func factorize(number int64) string {
	listPrimes := make([]int64, 0)

	currentPrime := int64(2)
	for currentPrime * currentPrime <= number {
		for number % currentPrime == 0 {
			listPrimes = append(listPrimes, currentPrime)
			number /= currentPrime
		}

		currentPrime++
	}

	if number != 1 {
		listPrimes = append(listPrimes, number)
	}

	return fmt.Sprint(listPrimes)
}

func validate(val string) (bool, error) {
	return true, nil
}

func Factorize(responseWriter http.ResponseWriter, request *http.Request) {
	if session.Manager.IsLogged(responseWriter, request) == false {
		responseWriter.WriteHeader(http.StatusUnauthorized)
		responseWriter.Write([]byte("401 - You should be authorized"))
		return
	}

	vars := mux.Vars(request)
	data := vars["number"]

	numberString := string(data)
	number, _ := strconv.ParseInt(numberString, 10, 64)
	result := factorize(number)
	responseWriter.Header().Set("Content-Type", "text/html")
	fmt.Fprint(responseWriter, result)
}
