package server

import (
	"data_base_project/service_logic"
	"data_base_project/types"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func SetupAuthorizeRouter(authService service_logic.IAuthService) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc(AUTH_AUTHORIZE, AuthorizeHandler(authService))

	return router
}

func AuthorizeHandler(authService service_logic.IAuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Receied request: %s %s", r.Method, r.URL.Path)
		if r.Method != "POST" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		var authData types.AuthData
		if err := json.Unmarshal(body, &authData); err != nil {
			log.Printf("Error unmarshaling request body: %v", err)
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}
		verdict, err := authService.Authorize(authData)
		if err != nil {
			log.Printf("Error authorizing: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(verdict)
		log.Printf("Authorized: %v", verdict)
	}
}
