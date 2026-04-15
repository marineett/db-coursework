package server

import (
	"data_base_project/service_logic"
	"data_base_project/types"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

func SetupContractRouter(
	contractService service_logic.IContractService,
	reviewService service_logic.IReviewService,
	lessonService service_logic.ILessonService,
) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc(CONTRACT_GET, ContractGetContractHandler(contractService))
	router.HandleFunc(CONTRACT_GET_REVIEW, ContractGetReviewHandler(reviewService))
	router.HandleFunc(ADD_LESSON, ContractAddLessonHandler(lessonService))
	router.HandleFunc(GET_LESSONS, ContractGetLessonsHandler(lessonService))
	return router
}

func ContractGetContractHandler(contractService service_logic.IContractService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
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
		contract, err := contractService.GetContract(int64(contractID))
		if err != nil {
			log.Printf("Error getting contract: %v", err)
			http.Error(w, "Error getting contract", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(contract)
		log.Printf("Contract retrieved: %v", contract)
		w.WriteHeader(http.StatusOK)
	}
}

func ContractGetReviewHandler(reviewService service_logic.IReviewService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		reviewIDStr := r.URL.Query().Get("review_id")
		reviewID, err := strconv.Atoi(reviewIDStr)
		if err != nil {
			log.Printf("Error converting reviewID to int: %v", err)
			http.Error(w, "Invalid reviewID", http.StatusBadRequest)
			return
		}
		review, err := reviewService.GetReview(int64(reviewID))
		if err != nil {
			log.Printf("Error getting review: %v", err)
			http.Error(w, "Error getting review", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(review)
	}
}

func ContractAddLessonHandler(lessonService service_logic.ILessonService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "POST" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		lesson := types.Lesson{}
		err := json.NewDecoder(r.Body).Decode(&lesson)
		if err != nil {
			log.Printf("Error decoding lesson: %v", err)
			http.Error(w, "Error decoding lesson", http.StatusBadRequest)
			return
		}
		lesson.CreatedAt = time.Now()
		lessonID, err := lessonService.CreateLesson(lesson)
		if err != nil {
			log.Printf("Error creating lesson: %v", err)
			http.Error(w, "Error creating lesson", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(lessonID)
	}
}

func ContractGetLessonsHandler(lessonService service_logic.ILessonService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		contractIDStr := r.URL.Query().Get("contract_id")
		contractID, err := strconv.ParseInt(contractIDStr, 10, 64)
		if err != nil {
			log.Printf("Error converting contractID to int: %v", err)
			http.Error(w, "Invalid contractID", http.StatusBadRequest)
			return
		}
		offsetStr := r.URL.Query().Get("lessons_offset")
		offset, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			log.Printf("Error converting from to int: %v", err)
			http.Error(w, "Invalid from", http.StatusBadRequest)
			return
		}
		sizeStr := r.URL.Query().Get("lessons_size")
		size, err := strconv.ParseInt(sizeStr, 10, 64)
		if err != nil {
			log.Printf("Error converting size to int: %v", err)
			http.Error(w, "Invalid size", http.StatusBadRequest)
			return
		}
		lessons, err := lessonService.GetLessons(contractID, offset, size)
		if err != nil {
			log.Printf("Error getting lessons: %v", err)
			http.Error(w, "Error getting lessons", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(lessons)
	}
}
