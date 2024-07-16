package util

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ErrResponse struct {
	Error string `json:"error"`
}

// RespondWithError sends an HTTP response with an error message.
// Used to standardize the way errors are returned to the client, also logs any 5XX server errors for server-side troubleshooting.
func RespondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5XX error: ", msg)
	}

	RespondWithJSON(w, code, ErrResponse{Error: msg})
}

func RespondWithInternalServerError(w http.ResponseWriter) {
	RespondWithJSON(w, 500, ErrResponse{Error: "Internal server error."})
}

func RespondWithMissingParametersError(w http.ResponseWriter) {
	RespondWithJSON(w, 400, ErrResponse{Error: "Missing one or more required parameters."})
}

func RespondWithBadRequest(w http.ResponseWriter, msg string) {
	RespondWithJSON(w, http.StatusBadRequest, ErrResponse{Error: msg})
}

// RespondWithJSON sends an HTTP response in JSON format.
// It is a generic function that can be used to return any payload as a JSON response.
func RespondWithJSON(w http.ResponseWriter, code int, payload any) {

	dat, err := json.Marshal(payload)

	if err != nil {
		fmt.Println(err)
		log.Printf("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/JSON")
	w.WriteHeader(code)
	w.Write(dat)
}
