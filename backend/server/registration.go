package server

import (
	"data_base_project/service_logic"
	"data_base_project/types"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func SetupRegistrationRouter(
	authService service_logic.IAuthService,
	moderatorService service_logic.IModeratorService,
	clientService service_logic.IClientService,
	adminService service_logic.IAdminService,
	repetitorService service_logic.IRepetitorService,
) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc(REGISTRATION_CLIENT, RegistrationClientHandler(clientService, authService))
	router.HandleFunc(REGISTRATION_MODERATOR, RegistrationModeratorHandler(moderatorService, authService))
	router.HandleFunc(REGISTRATION_ADMIN, RegistrationAdminHandler(adminService, authService))
	router.HandleFunc(REGISTRATION_REPETITOR, RegistrationRepetitorHandler(repetitorService, authService))
	return router
}

func RegistrationModeratorHandler(
	moderatorService service_logic.IModeratorService,
	authService service_logic.IAuthService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "POST" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		log.Printf("Request body: %s", string(body))

		var initData types.InitModeratorData
		if err := json.Unmarshal(body, &initData); err != nil {
			log.Printf("Error unmarshaling request body: %v", err)
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}
		log.Printf("Parsed init data: %+v", initData)

		inSystem, err := authService.CheckLogin(initData.InitUserData.AuthData.Login)
		if err != nil {
			log.Printf("Error checking login: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if inSystem {
			log.Printf("User already exists: %s", initData.InitUserData.AuthData.Login)
			http.Error(w, "User already exists", http.StatusBadRequest)
			return
		}

		err = moderatorService.CreateModerator(initData)
		if err != nil {
			log.Printf("Error creating moderator: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Moderator created successfully: %s", initData.InitUserData.AuthData.Login)
		w.WriteHeader(http.StatusCreated)
	}
}

func RegistrationClientHandler(
	clientService service_logic.IClientService,
	authService service_logic.IAuthService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)

		if r.Method != "POST" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		log.Printf("Request body: %s", string(body))

		var initData types.InitClientData
		if err := json.Unmarshal(body, &initData); err != nil {
			log.Printf("Error unmarshaling request body: %v", err)
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}
		log.Printf("Parsed init data: %+v", initData)

		inSystem, err := authService.CheckLogin(initData.InitUserData.AuthData.Login)
		if err != nil {
			log.Printf("Error checking login: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if inSystem {
			log.Printf("User already exists: %s", initData.InitUserData.AuthData.Login)
			http.Error(w, "User already exists", http.StatusBadRequest)
			return
		}

		err = clientService.CreateClient(initData)
		if err != nil {
			log.Printf("Error creating client: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Client created successfully: %s", initData.InitUserData.AuthData.Login)
		w.WriteHeader(http.StatusCreated)
	}
}

func RegistrationAdminHandler(
	adminService service_logic.IAdminService,
	authService service_logic.IAuthService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)

		if r.Method != "POST" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		log.Printf("Request body: %s", string(body))

		var initData types.InitAdminData
		if err := json.Unmarshal(body, &initData); err != nil {
			log.Printf("Error unmarshaling request body: %v", err)
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}
		log.Printf("Parsed init data: %+v", initData)

		inSystem, err := authService.CheckLogin(initData.InitUserData.AuthData.Login)
		if err != nil {
			log.Printf("Error checking login: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if inSystem {
			log.Printf("User already exists: %s", initData.InitUserData.AuthData.Login)
			http.Error(w, "User already exists", http.StatusBadRequest)
			return
		}

		err = adminService.CreateAdmin(initData)
		if err != nil {
			log.Printf("Error creating admin: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Admin created successfully: %s", initData.InitUserData.AuthData.Login)
		w.WriteHeader(http.StatusCreated)
	}
}

func RegistrationRepetitorHandler(
	repetitorService service_logic.IRepetitorService,
	authService service_logic.IAuthService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)

		if r.Method != "POST" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		log.Printf("Request body: %s", string(body))

		var initData types.InitRepetitorData
		if err := json.Unmarshal(body, &initData); err != nil {
			log.Printf("Error unmarshaling request body: %v", err)
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}
		log.Printf("Parsed init data: %+v", initData)

		inSystem, err := authService.CheckLogin(initData.InitUserData.AuthData.Login)
		if err != nil {
			log.Printf("Error checking login: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if inSystem {
			log.Printf("User already exists: %s", initData.InitUserData.AuthData.Login)
			http.Error(w, "User already exists", http.StatusBadRequest)
			return
		}

		err = repetitorService.CreateRepetitor(initData)
		if err != nil {
			log.Printf("Error creating repetitor: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Repetitor created successfully: %s", initData.InitUserData.AuthData.Login)
		w.WriteHeader(http.StatusCreated)
	}
}
