package server

import (
	"data_base_project/service_logic"
	"data_base_project/types"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

func SetupRepetitorRouter(
	repetitorService service_logic.IRepetitorService,
	contractService service_logic.IContractService,
	transactionService service_logic.ITransactionService,
	resumeService service_logic.IResumeService,
) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc(REPETITOR_GET_PROFILE, RepetitorGetProfileHandler(repetitorService))
	router.HandleFunc(REPETITOR_GET_CONTRACTS, RepetitorGetContractsHandler(contractService))
	router.HandleFunc(REPETITOR_GET_AVAILABLE_CONTRACTS, RepetitorGetAvailableContractsHandler(contractService))
	router.HandleFunc(REPETITOR_ACCEPT_CONTRACT, RepetitorAcceptContractHandler(repetitorService, contractService))
	router.HandleFunc(REPETITOR_MAKE_REVIEW, RepetitorMakeReviewHandler(contractService))
	router.HandleFunc(REPETITOR_PAY_FOR_CONTRACT, RepetitorPayForContractHandler(repetitorService, contractService, transactionService))
	router.HandleFunc(REPETITOR_CANCEL_CONTRACT, RepetitorCancelContractHandler(repetitorService, contractService))
	router.HandleFunc(REPETITOR_COMPLETE_CONTRACT, RepetitorCompleteContractHandler(repetitorService, contractService))
	router.HandleFunc(REPETITOR_CHANGE_RESUME, RepetitorChangeResumeHandler(repetitorService, resumeService))
	return router
}

func RepetitorGetProfileHandler(repetitorService service_logic.IRepetitorService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		repetitorIDStr := r.URL.Query().Get("id")
		repetitorID, err := strconv.ParseInt(repetitorIDStr, 10, 64)
		if err != nil {
			log.Printf("Error converting repetitorID to int: %v", err)
			http.Error(w, "Invalid repetitorID", http.StatusBadRequest)
			return
		}
		reviewsOffsetStr := r.URL.Query().Get("reviews_offset")
		reviewsOffset, err := strconv.ParseInt(reviewsOffsetStr, 10, 64)
		if err != nil {
			log.Printf("Error converting reviewsOffset to int: %v", err)
			http.Error(w, "Invalid reviewsOffset", http.StatusBadRequest)
			return
		}
		reviewsLimitStr := r.URL.Query().Get("reviews_limit")
		reviewsLimit, err := strconv.ParseInt(reviewsLimitStr, 10, 64)
		if err != nil {
			log.Printf("Error converting reviewsLimit to int: %v", err)
			http.Error(w, "Invalid reviewsLimit", http.StatusBadRequest)
			return
		}
		repetitor, err := repetitorService.GetRepetitorProfile(int64(repetitorID), int64(reviewsOffset), int64(reviewsLimit))
		if err != nil {
			log.Printf("Error getting repetitor: %v", err)
			http.Error(w, "Error getting repetitor", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(repetitor)
		log.Printf("Repetitor retrieved: %v", repetitor)
	}
}

func RepetitorGetContractsHandler(contractService service_logic.IContractService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		repetitorIDStr := r.URL.Query().Get("repetitor_id")
		repetitorID, err := strconv.Atoi(repetitorIDStr)
		if err != nil {
			log.Printf("Error converting repetitorID to int: %v", err)
			http.Error(w, "Invalid repetitorID", http.StatusBadRequest)
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
		log.Printf("id: %d, offset: %d, limit: %d, status: %d", repetitorID, offset, limit, status)
		contracts, err := contractService.GetRepetitorContractList(int64(repetitorID), int64(offset), int64(limit), types.ContractStatus(status))
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

func RepetitorGetAvailableContractsHandler(
	contractService service_logic.IContractService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
		contracts, err := contractService.GetContractList(int64(offset), int64(limit), types.ContractStatus(status))
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

func RepetitorAcceptContractHandler(
	repetitorService service_logic.IRepetitorService,
	contractService service_logic.IContractService,
) http.HandlerFunc {
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

		repetitorIDStr := r.URL.Query().Get("repetitor_id")
		repetitorID, err := strconv.Atoi(repetitorIDStr)
		if err != nil {
			log.Printf("Error converting repetitorID to int: %v", err)
			http.Error(w, "Invalid repetitorID", http.StatusBadRequest)
			return
		}
		_, err = repetitorService.GetRepetitorData(int64(repetitorID))
		if err != nil {
			log.Printf("Error getting repetitor: %v", err)
			http.Error(w, "Error getting repetitor", http.StatusInternalServerError)
			return
		}
		err = contractService.UpdateContractRepetitorID(int64(contractID), int64(repetitorID))
		if err != nil {
			log.Printf("Error updating contract repetitor: %v", err)
			http.Error(w, "Error updating contract repetitor", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Printf("Contract %d updated with repetitor %d", contractID, repetitorID)
	}
}

func RepetitorMakeReviewHandler(
	contractService service_logic.IContractService,
) http.HandlerFunc {
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
		err = contractService.CreateContractReviewRepetitor(int64(contractID), review)
		if err != nil {
			log.Printf("Error creating review: %v", err)
			http.Error(w, "Error creating review", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		log.Printf("Review created: %v", review)
	}
}

func RepetitorPayForContractHandler(
	repetitorService service_logic.IRepetitorService,
	contractService service_logic.IContractService,
	transactionService service_logic.ITransactionService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "POST" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		userIDStr := r.URL.Query().Get("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			log.Printf("Error converting userID to int: %v", err)
			http.Error(w, "Invalid userID", http.StatusBadRequest)
			return
		}
		contractIDStr := r.URL.Query().Get("contract_id")
		contractID, err := strconv.Atoi(contractIDStr)
		if err != nil {
			log.Printf("Error converting contractID to int: %v", err)
			http.Error(w, "Invalid contractID", http.StatusBadRequest)
			return
		}
		amountStr := r.URL.Query().Get("amount")
		amount, err := strconv.Atoi(amountStr)
		if err != nil {
			log.Printf("Error converting amount to int: %v", err)
			http.Error(w, "Invalid amount", http.StatusBadRequest)
			return
		}
		_, err = contractService.GetContract(int64(contractID))
		if err != nil {
			log.Printf("Error getting contract: %v", err)
			http.Error(w, "Error getting contract", http.StatusInternalServerError)
			return
		}
		transactionID, err := transactionService.CreateContractPaymentTransaction(int64(amount), int64(userID))
		if err != nil {
			log.Printf("Error creating transaction: %v", err)
			http.Error(w, "Error creating transaction", http.StatusInternalServerError)
			return
		}
		err = contractService.UpdateContractPaymentStatus(int64(contractID), types.PaymentStatusPaid)
		if err != nil {
			log.Printf("Error updating contract status: %v", err)
			http.Error(w, "Error updating contract status", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		log.Printf("Transaction created: %v", transactionID)
	}
}

func RepetitorCancelContractHandler(
	repetitorService service_logic.IRepetitorService,
	contractService service_logic.IContractService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		repetitorIDStr := r.URL.Query().Get("id")
		repetitorID, err := strconv.Atoi(repetitorIDStr)
		if err != nil {
			log.Printf("Error converting repetitorID to int: %v", err)
			http.Error(w, "Invalid repetitorID", http.StatusBadRequest)
			return
		}
		contractIDStr := r.URL.Query().Get("c_id")
		contractID, err := strconv.Atoi(contractIDStr)
		if err != nil {
			log.Printf("Error converting contractID to int: %v", err)
			http.Error(w, "Invalid contractID", http.StatusBadRequest)
			return
		}
		contract, err := contractService.GetContract(int64(contractID))
		if err != nil {
			log.Printf("Error getting contract: %v", err)
			http.Error(w, "Error getting contract", http.StatusInternalServerError)
			return
		}
		if contract.RepetitorID != int64(repetitorID) {
			log.Printf("Contract doesn't belong to repetitor")
			http.Error(w, "Contract doesn't belong to repetitor", http.StatusNotFound)
			return
		}
		if contract.Status != types.ContractStatusActive {
			log.Printf("Contract is not active")
			http.Error(w, "Contract is not active", http.StatusBadRequest)
			return
		}
		err = contractService.UpdateContractStatus(int64(contractID), types.ContractStatusCancelled)
		if err != nil {
			log.Printf("Error updating contract status: %v", err)
			http.Error(w, "Error updating contract status", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func RepetitorCompleteContractHandler(
	repetitorService service_logic.IRepetitorService,
	contractService service_logic.IContractService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		repetitorIDStr := r.URL.Query().Get("id")
		repetitorID, err := strconv.Atoi(repetitorIDStr)
		if err != nil {
			log.Printf("Error converting repetitorID to int: %v", err)
			http.Error(w, "Invalid repetitorID", http.StatusBadRequest)
			return
		}
		contractIDStr := r.URL.Query().Get("c_id")
		contractID, err := strconv.Atoi(contractIDStr)
		if err != nil {
			log.Printf("Error converting contractID to int: %v", err)
			http.Error(w, "Invalid contractID", http.StatusBadRequest)
			return
		}
		contract, err := contractService.GetContract(int64(contractID))
		if err != nil {
			log.Printf("Error getting contract: %v", err)
			http.Error(w, "Error getting contract", http.StatusInternalServerError)
			return
		}
		if contract.RepetitorID != int64(repetitorID) {
			log.Printf("Contract doesn't belong to repetitor")
			http.Error(w, "Contract doesn't belong to repetitor", http.StatusNotFound)
			return
		}
		err = contractService.UpdateContractStatus(int64(contractID), types.ContractStatusCompleted)
		if err != nil {
			log.Printf("Error updating contract status: %v", err)
			http.Error(w, "Error updating contract status", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func RepetitorChangeResumeHandler(
	repetitorService service_logic.IRepetitorService,
	resumeService service_logic.IResumeService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "POST" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		repetitorIDStr := r.URL.Query().Get("id")
		repetitorID, err := strconv.ParseInt(repetitorIDStr, 10, 64)
		if err != nil {
			log.Printf("Error converting repetitorID to int: %v", err)
			http.Error(w, "Invalid repetitorID", http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		repetitor, err := repetitorService.GetRepetitorData(repetitorID)
		if err != nil {
			log.Printf("Error getting repetitor: %v", err)
			http.Error(w, "Error getting repetitor", http.StatusInternalServerError)
			return
		}
		resume := types.Resume{}
		if err := json.Unmarshal(body, &resume); err != nil {
			log.Printf("Error unmarshalling request body: %v", err)
			http.Error(w, "Error unmarshalling request body", http.StatusBadRequest)
			return
		}
		err = resumeService.UpdateResumeTitle(repetitor.ResumeID, resume.Title)
		if err != nil {
			log.Printf("Error updating repetitor resume: %v", err)
			http.Error(w, "Error updating repetitor resume", http.StatusInternalServerError)
			return
		}
		err = resumeService.UpdateResumeDescription(repetitor.ResumeID, resume.Description)
		if err != nil {
			log.Printf("Error updating repetitor resume: %v", err)
			http.Error(w, "Error updating repetitor resume", http.StatusInternalServerError)
			return
		}
		err = resumeService.UpdateResumePrices(repetitor.ResumeID, resume.Prices)
		if err != nil {
			log.Printf("Error updating repetitor resume: %v", err)
			http.Error(w, "Error updating repetitor resume", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
