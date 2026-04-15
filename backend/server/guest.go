package server

import (
	"data_base_project/service_logic"
	"encoding/json"
	"net/http"
	"strconv"
)

func SetupGuestRouter(
	repetitorService service_logic.IRepetitorService,
) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc(GUEST_GET_REPETITORS, GuestGetRepetitorsHandler(repetitorService))

	return router
}

func GuestGetRepetitorsHandler(repetitorService service_logic.IRepetitorService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		repetitorsOffsetStr := r.URL.Query().Get("repetitors_offset")
		repetitorsOffset, err := strconv.ParseInt(repetitorsOffsetStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid offset", http.StatusBadRequest)
			return
		}
		repetitorsLimitStr := r.URL.Query().Get("repetitors_limit")
		repetitorsLimit, err := strconv.ParseInt(repetitorsLimitStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
		repetitors, err := repetitorService.GetRepetitors(repetitorsOffset, repetitorsLimit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(repetitors)
	}
}
