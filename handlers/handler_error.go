package handlers

import (
	"net/http"

	"github.com/yizhong187/EduMind-backend/internal/util"
)

func HandlerError(w http.ResponseWriter, r *http.Request) {
	util.RespondWithError(w, 400, "Something went wrong :(")
}
