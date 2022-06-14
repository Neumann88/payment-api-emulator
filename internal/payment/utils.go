package payment

import (
	"errors"
	"net/http"
	"net/mail"
	"strconv"

	"github.com/gorilla/mux"
)

func getQueryId(r *http.Request) (int64, error) {
	id := mux.Vars(r)["id"]
	convertedID, err := converteIDtoI64(id)

	if err != nil {
		return 0, err
	}

	return convertedID, nil
}

func converteIDtoI64(id string) (int64, error) {
	convertedID, err := strconv.Atoi(id)

	if err != nil {
		return 0, err
	}

	return int64(convertedID), nil
}

func checkTerminalStatusRow(row int64) error {
	if row == 0 {
		return errors.New("terminal status")
	}

	return nil
}

func isEmail(address string) bool {
	_, err := mail.ParseAddress(address)

	return err == nil
}
