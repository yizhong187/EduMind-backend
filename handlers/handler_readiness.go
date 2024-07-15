package handlers

import (
	"net/http"

	"github.com/yizhong187/EduMind-backend/contextKeys"
	"github.com/yizhong187/EduMind-backend/internal/config"
	"github.com/yizhong187/EduMind-backend/internal/util"
)

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {

	apiCfg, ok := r.Context().Value(contextKeys.ConfigKey).(*config.ApiConfig)
	if !ok || apiCfg == nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Configuration not found")
		return
	}

	util.RespondWithJSON(w, http.StatusOK, "Service ready")
}
