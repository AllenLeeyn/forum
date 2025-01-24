package handlers

import "net/http"

//	Place for custom errors and error checking functions

// Made this for taking some http value (like method) and comparing it to "shouldBe"
// and if they don't match then throw an error
func IsHttpError(w http.ResponseWriter, r *http.Request, shouldBe interface{}, errMessage string, errStatus int, check string) {
	if check == "method" {
		if r.Method != shouldBe {
			ExecuteError(w, errMessage, errStatus)
		}
	}
}
