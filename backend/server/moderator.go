package server

import (
	"data_base_project/service_logic"
	"data_base_project/types"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func SetupModeratorRouter(
	transactionService service_logic.ITransactionService,
	contractService service_logic.IContractService,
	moderatorService service_logic.IModeratorService,
) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc(MODERATOR_GET_PROFILE, ModeratorGetProfileHandler(moderatorService))
	router.HandleFunc(MODERATOR_GET_TRANSACTION_TO_APPROVE, ModeratorGetTransactionsToApproveHandler(transactionService))
	router.HandleFunc(MODERATOR_APPROVE_TRANSACTION, ModeratorApproveTransactionHandler(transactionService))
	router.HandleFunc(MODERATOR_GET_CONTRACTS, ModeratorGetContractsHandler(contractService))
	router.HandleFunc(MODERATOR_BAN_CONTRACT, ModeratorBanContractHandler(contractService))
	return router
}

func ModeratorGetProfileHandler(moderatorService service_logic.IModeratorService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		moderatorIDStr := r.URL.Query().Get("id")
		moderatorID, err := strconv.ParseInt(moderatorIDStr, 10, 64)
		if err != nil {
			log.Printf("Invalid moderator ID: %v", err)
			http.Error(w, "Invalid moderator ID", http.StatusBadRequest)
			return
		}
		moderator, err := moderatorService.GetModeratorProfile(moderatorID)
		if err != nil {
			log.Printf("Error getting moderator: %v", err)
			http.Error(w, "Error getting moderator", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(moderator)
	}
}

func ModeratorGetTransactionsToApproveHandler(transactionService service_logic.ITransactionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		transaction, err := transactionService.GetPendingContractPaymentTransaction()
		if err != nil {
			log.Printf("Error getting transaction: %v", err)
			http.Error(w, "Error getting transaction", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(transaction)
	}
}

func ModeratorApproveTransactionHandler(transactionService service_logic.ITransactionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		transactionIDStr := r.URL.Query().Get("id")
		if transactionIDStr == "" {
			log.Printf("Transaction ID is required")
			http.Error(w, "Transaction ID is required", http.StatusBadRequest)
			return
		}
		transactionID, err := strconv.ParseInt(transactionIDStr, 10, 64)
		if err != nil {
			log.Printf("Invalid transaction ID: %v", err)
			http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
			return
		}
		err = transactionService.ApproveTransaction(transactionID)
		if err != nil {
			log.Printf("Error approving transaction: %v", err)
			http.Error(w, "Error approving transaction", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func ModeratorGetContractsHandler(contractService service_logic.IContractService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		fromStr := r.URL.Query().Get("from")
		if fromStr == "" {
			log.Printf("From is required")
			http.Error(w, "From is required", http.StatusBadRequest)
			return
		}
		from, err := strconv.ParseInt(fromStr, 10, 64)
		if err != nil {
			log.Printf("Invalid from: %v", err)
			http.Error(w, "Invalid from", http.StatusBadRequest)
			return
		}
		sizeStr := r.URL.Query().Get("size")
		if sizeStr == "" {
			log.Printf("Size is required")
			http.Error(w, "Size is required", http.StatusBadRequest)
			return
		}
		size, err := strconv.ParseInt(sizeStr, 10, 64)
		if err != nil {
			log.Printf("Invalid size: %v", err)
			http.Error(w, "Invalid size", http.StatusBadRequest)
			return
		}
		contracts, err := contractService.GetAllContracts(from, size)
		if err != nil {
			log.Printf("Error getting contracts: %v", err)
			http.Error(w, "Error getting contracts", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(contracts)
	}
}

func ModeratorBanContractHandler(contractService service_logic.IContractService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "POST" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		contractIDStr := r.URL.Query().Get("id")
		if contractIDStr == "" {
			log.Printf("Contract ID is required")
			http.Error(w, "Contract ID is required", http.StatusBadRequest)
			return
		}
		contractID, err := strconv.ParseInt(contractIDStr, 10, 64)
		if err != nil {
			log.Printf("Invalid contract ID: %v", err)
			http.Error(w, "Invalid contract ID", http.StatusBadRequest)
			return
		}
		contract, err := contractService.GetContract(contractID)
		if err != nil {
			log.Printf("Error getting contract: %v", err)
			http.Error(w, "Error getting contract", http.StatusInternalServerError)
			return
		}
		if contract.Status != types.ContractStatusActive && contract.Status != types.ContractStatusPending {
			log.Printf("Contract is not in valid status")
			http.Error(w, "Contract is not in valid status", http.StatusBadRequest)
			return
		}
		err = contractService.UpdateContractStatus(contractID, types.ContractStatusBanned)
		if err != nil {
			log.Printf("Error banning contract: %v", err)
			http.Error(w, "Error banning contract", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
