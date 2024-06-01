package handlers

import (
	"net/http"

	"github.com/yizhong187/EduMind-backend/internal/config"
	"github.com/yizhong187/EduMind-backend/internal/util"
	"github.com/yizhong187/EduMind-backend/middlewares"
)

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {

	apiCfg := r.Context().Value(middlewares.ConfigKey).(config.ApiConfig)

	util.RespondWithJSON(w, 200, apiCfg)
}
