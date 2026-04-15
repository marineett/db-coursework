package server

import (
	"data_base_project/service_logic"
	"data_base_project/types"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func SetupClientRouter(
	clientService service_logic.IClientService,
	contractService service_logic.IContractService,
) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc(CLIENT_GET_PROFILE, ClientGetProfileHandler(clientService))
	router.HandleFunc(CLIENT_CREATE_CONTRACT, ClientCreateContractHandler(contractService))
	router.HandleFunc(CLIENT_GET_CONTRACTS, ClientGetContractsHandler(contractService))
	router.HandleFunc(CLIENT_MAKE_REVIEW, ClientMakeReviewHandler(contractService))
	return router
}

func ClientGetProfileHandler(clientService service_logic.IClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		clientIDStr := r.URL.Query().Get("id")
		clientID, err := strconv.Atoi(clientIDStr)
		if err != nil {
			log.Printf("Error converting clientID to int: %v", err)
			http.Error(w, "Invalid clientID", http.StatusBadRequest)
			return
		}
		reviewsOffsetStr := r.URL.Query().Get("reviews_offset")
		reviewsOffset, err := strconv.Atoi(reviewsOffsetStr)
		if err != nil {
			log.Printf("Error converting reviewsOffset to int: %v", err)
			http.Error(w, "Invalid reviewsOffset", http.StatusBadRequest)
			return
		}
		reviewsLimitStr := r.URL.Query().Get("reviews_limit")
		reviewsLimit, err := strconv.Atoi(reviewsLimitStr)
		if err != nil {
			log.Printf("Error converting reviewsLimit to int: %v", err)
			http.Error(w, "Invalid reviewsLimit", http.StatusBadRequest)
			return
		}
		client, err := clientService.GetClientProfile(int64(clientID), int64(reviewsOffset), int64(reviewsLimit))
		if err != nil {
			http.Error(w, "Error getting client", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(client)
		log.Printf("Client retrieved: %v", client)
	}
}

func ClientCreateContractHandler(contractService service_logic.IContractService) http.HandlerFunc {
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
		contractInitInfo := types.ContractInitInfo{}
		if err := json.Unmarshal(body, &contractInitInfo); err != nil {
			log.Printf("Error unmarshalling request body: %v", err)
			http.Error(w, "Error unmarshalling request body", http.StatusBadRequest)
			return
		}
		contractID, err := contractService.CreateContract(contractInitInfo)
		if err != nil {
			log.Printf("Error creating contract: %v", err)
			http.Error(w, "Error creating contract", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(contractID)
		log.Printf("Contract created with ID: %d", contractID)
		w.WriteHeader(http.StatusCreated)
	}
}

func ClientGetContractsHandler(contractService service_logic.IContractService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		clientIDStr := r.URL.Query().Get("client_id")
		clientID, err := strconv.Atoi(clientIDStr)
		if err != nil {
			log.Printf("Error converting clientID to int: %v", err)
			http.Error(w, "Invalid clientID", http.StatusBadRequest)
			return
		}
		offsetStr := r.URL.Query().Get("offset")
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			log.Printf("Error converting offset to int: %v", err)
			http.Error(w, "Invalid offset", http.StatusBadRequest)
			return
		}
		limitStr := r.URL.Query().Get("limit")
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			log.Printf("Error converting limit to int: %v", err)
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
		statusStr := r.URL.Query().Get("status")
		status, err := strconv.Atoi(statusStr)
		if err != nil {
			log.Printf("Error converting status to int: %v", err)
			http.Error(w, "Invalid status", http.StatusBadRequest)
			return
		}
		contracts, err := contractService.GetClientContractList(int64(clientID), int64(offset), int64(limit), types.ContractStatus(status))
		if err != nil {
			log.Printf("Error getting contracts: %v", err)
			http.Error(w, "Error getting contracts", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(contracts)
		log.Printf("Contracts retrieved: %v", contracts)
	}
}

func ClientMakeReviewHandler(contractService service_logic.IContractService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "POST" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		contractIDStr := r.URL.Query().Get("contract_id")
		contractID, err := strconv.Atoi(contractIDStr)
		if err != nil {
			log.Printf("Error converting contractID to int: %v", err)
			http.Error(w, "Invalid contractID", http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		review := types.Review{}
		if err := json.Unmarshal(body, &review); err != nil {
			log.Printf("Error unmarshalling request body: %v", err)
			http.Error(w, "Error unmarshalling request body", http.StatusBadRequest)
			return
		}
		review.CreatedAt = time.Now()
		err = contractService.CreateContractReviewClient(int64(contractID), review)
		if err != nil {
			log.Printf("Error making review: %v", err)
			http.Error(w, "Error making review", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		log.Printf("Review made: %v", review)
	}
}
