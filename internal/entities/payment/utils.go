package payment

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func getQueryId(r *http.Request) (int64, error) {
	id := mux.Vars(r)["id"]
	paymentID, err := strconv.Atoi(id)

	if err != nil {
		return 0, err
	}

	return int64(paymentID), nil
}
