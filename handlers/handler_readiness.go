package handlers

import (
	"net/http"

	"github.com/yizhong187/EduMind-backend/internal/util"
)

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {

	util.RespondWithJSON(w, 200, struct{}{})
}
