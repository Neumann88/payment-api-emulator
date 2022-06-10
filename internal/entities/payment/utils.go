package payment

import (
	"net/http"
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
