package utils

import (
	"net/http"
	"net/mail"
	"strconv"

	"github.com/gorilla/mux"
)

func GetQueryId(r *http.Request) (int64, error) {
	id := mux.Vars(r)["id"]
	convertedID, err := ConverteIDtoI64(id)

	if err != nil {
		return 0, err
	}

	return convertedID, nil
}

func ConverteIDtoI64(id string) (int64, error) {
	convertedID, err := strconv.Atoi(id)

	if err != nil {
		return 0, err
	}

	return int64(convertedID), nil
}

func IsEmail(address string) bool {
	_, err := mail.ParseAddress(address)

	return err == nil
}
