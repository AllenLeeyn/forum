package utils

import "net/http"

func IsError(w http.ResponseWriter, r *http.Request, shouldBe interface{}, errMessage string, errStatus int, check string) {
	if check == "method" {
		if r.Method != shouldBe {
			http.Error(w, errMessage, errStatus)
		}
	}
}
