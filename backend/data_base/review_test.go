package data_base

import (
	"data_base_project/types"
	"testing"
	"time"
)

func TestCreateReviewTable(t *testing.T) {
	reviewRepository := CreateReviewRepository(globalDb, "test_review_table")
	if reviewRepository == nil {
		t.Errorf("Failed to create review repository")
	}
}

func TestInsertReview(t *testing.T) {
	InsertTestUser(1)
	InsertTestUser(2)
	defer globalDb.Exec("TRUNCATE TABLE test_review_table, test_user_table CASCADE")
	reviewRepository := CreateReviewRepository(globalDb, "test_review_table")
	review := types.Review{
		ClientID:    1,
		RepetitorID: 2,
		Rating:      5,
		Comment:     "Good job!",
		CreatedAt:   time.Now(),
	}
	id, err := reviewRepository.InsertReview(review)
	if err != nil {
		t.Errorf("Failed to insert review: %v", err)
	}
	resultReview := types.Review{}
	err = globalDb.QueryRow("SELECT * FROM test_review_table WHERE id = $1", id).Scan(&resultReview.ID, &resultReview.ClientID, &resultReview.RepetitorID, &resultReview.Rating, &resultReview.Comment, &resultReview.CreatedAt)
	if err != nil {
		t.Errorf("Failed to get review: %v", err)
	}
	if resultReview.ID != id {
		t.Errorf("Failed to get review: %v", resultReview.ID)
	}
	if resultReview.ClientID != 1 {
		t.Errorf("Failed to get review: %v", resultReview.ClientID)
	}
	if resultReview.RepetitorID != 2 {
		t.Errorf("Failed to get review: %v", resultReview.RepetitorID)
	}
	if resultReview.Rating != 5 {
		t.Errorf("Failed to get review: %v", resultReview.Rating)
	}
	if resultReview.Comment != "Good job!" {
		t.Errorf("Failed to get review: %v", resultReview.Comment)
	}
}

func TestInsertReviewInSeq(t *testing.T) {
	InsertTestUser(1)
	InsertTestUser(2)
	defer globalDb.Exec("TRUNCATE TABLE test_review_table, test_user_table CASCADE")
	reviewRepository := CreateReviewRepository(globalDb, "test_review_table")
	tx, err := globalDb.Begin()
	if err != nil {
		t.Errorf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()
	review := types.Review{
		ClientID:    1,
		RepetitorID: 2,
		Rating:      5,
		Comment:     "Good job!",
		CreatedAt:   time.Now(),
	}
	id, err := reviewRepository.InsertReviewInSeq(tx, review)
	if err != nil {
		t.Errorf("Failed to insert review: %v", err)
	}
	resultReview := types.Review{}
	err = tx.QueryRow("SELECT * FROM test_review_table WHERE id = $1", id).Scan(&resultReview.ID, &resultReview.ClientID, &resultReview.RepetitorID, &resultReview.Rating, &resultReview.Comment, &resultReview.CreatedAt)
	if err != nil {
		t.Errorf("Failed to get review: %v", err)
	}
	if resultReview.ID != id {
		t.Errorf("Failed to get review: %v", resultReview.ID)
	}
	if resultReview.ClientID != 1 {
		t.Errorf("Failed to get review: %v", resultReview.ClientID)
	}
	if resultReview.RepetitorID != 2 {
		t.Errorf("Failed to get review: %v", resultReview.RepetitorID)
	}
	if resultReview.Rating != 5 {
		t.Errorf("Failed to get review: %v", resultReview.Rating)
	}
	if resultReview.Comment != "Good job!" {
		t.Errorf("Failed to get review: %v", resultReview.Comment)
	}
}

func TestGetReview(t *testing.T) {
	InsertTestUser(1)
	InsertTestUser(2)
	defer globalDb.Exec("TRUNCATE TABLE test_review_table, test_user_table CASCADE")
	reviewRepository := CreateReviewRepository(globalDb, "test_review_table")
	review := types.Review{
		ClientID:    1,
		RepetitorID: 2,
		Rating:      5,
		Comment:     "Good job!",
		CreatedAt:   time.Now(),
	}
	id, err := reviewRepository.InsertReview(review)
	if err != nil {
		t.Errorf("Failed to insert review: %v", err)
	}
	resultReview, err := reviewRepository.GetReview(id)
	if err != nil {
		t.Errorf("Failed to get review: %v", err)
	}
	if resultReview.ID != id {
		t.Errorf("Failed to get review: %v", resultReview.ID)
	}
	if resultReview.ClientID != 1 {
		t.Errorf("Failed to get review: %v", resultReview.ClientID)
	}
	if resultReview.RepetitorID != 2 {
		t.Errorf("Failed to get review: %v", resultReview.RepetitorID)
	}
	if resultReview.Rating != 5 {
		t.Errorf("Failed to get review: %v", resultReview.Rating)
	}
	if resultReview.Comment != "Good job!" {
		t.Errorf("Failed to get review: %v", resultReview.Comment)
	}
}

func TestGetReviewsByClientID(t *testing.T) {
	InsertTestUser(1)
	InsertTestUser(2)
	InsertTestUser(3)
	InsertTestUser(4)
	defer globalDb.Exec("TRUNCATE TABLE test_review_table, test_user_table CASCADE")
	reviewRepository := CreateReviewRepository(globalDb, "test_review_table")
	review := types.Review{
		ClientID:    1,
		RepetitorID: 2,
		Rating:      5,
		Comment:     "Good job!",
		CreatedAt:   time.Now(),
	}
	firstId, err := reviewRepository.InsertReview(review)
	if err != nil {
		t.Errorf("Failed to insert review: %v", err)
	}
	review2 := types.Review{
		ClientID:    2,
		RepetitorID: 3,
		Rating:      4,
		Comment:     "Fine job!",
		CreatedAt:   time.Now(),
	}
	_, err = reviewRepository.InsertReview(review2)
	if err != nil {
		t.Errorf("Failed to insert review: %v", err)
	}
	review3 := types.Review{
		ClientID:    1,
		RepetitorID: 4,
		Rating:      3,
		Comment:     "Bad job!",
		CreatedAt:   time.Now(),
	}
	thirdId, err := reviewRepository.InsertReview(review3)
	if err != nil {
		t.Errorf("Failed to insert review: %v", err)
	}
	resultReviews, err := reviewRepository.GetReviewsByClientID(1, 0, 10)
	if err != nil {
		t.Errorf("Failed to get review: %v", err)
	}
	if len(resultReviews) != 2 {
		t.Errorf("Failed to get review: %v", resultReviews)
	}
	if resultReviews[0].ID != thirdId {
		t.Errorf("Failed to get review: %v", resultReviews[0].ID)
	}
	if resultReviews[1].ID != firstId {
		t.Errorf("Failed to get review: %v", resultReviews[1].ID)
	}
	if resultReviews[0].Rating != 3 {
		t.Errorf("Failed to get review: %v", resultReviews[0].Rating)
	}
	if resultReviews[1].Rating != 5 {
		t.Errorf("Failed to get review: %v", resultReviews[1].Rating)
	}
	if resultReviews[0].Comment != "Bad job!" {
		t.Errorf("Failed to get review: %v", resultReviews[0].Comment)
	}
	if resultReviews[1].Comment != "Good job!" {
		t.Errorf("Failed to get review: %v", resultReviews[1].Comment)
	}
	resultReviews, err = reviewRepository.GetReviewsByClientID(1, 1, 10)
	if err != nil {
		t.Errorf("Failed to get review: %v", err)
	}
	if len(resultReviews) != 1 {
		t.Errorf("Failed to get review: %v", resultReviews)
	}
	if resultReviews[0].ID != firstId {
		t.Errorf("Failed to get review: %v", resultReviews[0].ID)
	}
	if resultReviews[0].Rating != 5 {
		t.Errorf("Failed to get review: %v", resultReviews[0].Rating)
	}
	if resultReviews[0].Comment != "Good job!" {
		t.Errorf("Failed to get review: %v", resultReviews[0].Comment)
	}
	resultReviews, err = reviewRepository.GetReviewsByClientID(3, 0, 10)
	if err != nil {
		t.Errorf("Failed to get review: %v", err)
	}
	if len(resultReviews) != 0 {
		t.Errorf("Failed to get review: %v", resultReviews)
	}
}

func TestGetReviewsByRepetitorID(t *testing.T) {
	InsertTestUser(1)
	InsertTestUser(2)
	InsertTestUser(3)
	InsertTestUser(4)
	defer globalDb.Exec("TRUNCATE TABLE test_review_table, test_user_table CASCADE")
	reviewRepository := CreateReviewRepository(globalDb, "test_review_table")
	review := types.Review{
		ClientID:    1,
		RepetitorID: 2,
		Rating:      5,
		Comment:     "Good job!",
		CreatedAt:   time.Now(),
	}
	firstId, err := reviewRepository.InsertReview(review)
	if err != nil {
		t.Errorf("Failed to insert review: %v", err)
	}
	review2 := types.Review{
		ClientID:    4,
		RepetitorID: 1,
		Rating:      4,
		Comment:     "Fine job!",
		CreatedAt:   time.Now(),
	}
	_, err = reviewRepository.InsertReview(review2)
	if err != nil {
		t.Errorf("Failed to insert review: %v", err)
	}
	review3 := types.Review{
		ClientID:    3,
		RepetitorID: 2,
		Rating:      3,
		Comment:     "Bad job!",
		CreatedAt:   time.Now(),
	}
	thirdId, err := reviewRepository.InsertReview(review3)
	if err != nil {
		t.Errorf("Failed to insert review: %v", err)
	}
	resultReviews, err := reviewRepository.GetReviewsByRepetitorID(2, 0, 10)
	if err != nil {
		t.Errorf("Failed to get review: %v", err)
	}
	if len(resultReviews) != 2 {
		t.Errorf("Failed to get review: %v", resultReviews)
	}
	if resultReviews[0].ID != thirdId {
		t.Errorf("Failed to get review: %v", resultReviews[0].ID)
	}
	if resultReviews[1].ID != firstId {
		t.Errorf("Failed to get review: %v", resultReviews[1].ID)
	}
	if resultReviews[0].Rating != 3 {
		t.Errorf("Failed to get review: %v", resultReviews[0].Rating)
	}
	if resultReviews[1].Rating != 5 {
		t.Errorf("Failed to get review: %v", resultReviews[1].Rating)
	}
	if resultReviews[0].Comment != "Bad job!" {
		t.Errorf("Failed to get review: %v", resultReviews[0].Comment)
	}
	if resultReviews[1].Comment != "Good job!" {
		t.Errorf("Failed to get review: %v", resultReviews[1].Comment)
	}
	resultReviews, err = reviewRepository.GetReviewsByRepetitorID(2, 1, 10)
	if err != nil {
		t.Errorf("Failed to get review: %v", err)
	}
	if len(resultReviews) != 1 {
		t.Errorf("Failed to get review: %v", resultReviews)
	}
	if resultReviews[0].ID != firstId {
		t.Errorf("Failed to get review: %v", resultReviews[0].ID)
	}
	if resultReviews[0].Rating != 5 {
		t.Errorf("Failed to get review: %v", resultReviews[0].Rating)
	}
	if resultReviews[0].Comment != "Good job!" {
		t.Errorf("Failed to get review: %v", resultReviews[0].Comment)
	}
	resultReviews, err = reviewRepository.GetReviewsByRepetitorID(3, 0, 10)
	if err != nil {
		t.Errorf("Failed to get review: %v", err)
	}
	if len(resultReviews) != 0 {
		t.Errorf("Failed to get review: %v", resultReviews)
	}
}

func TestUpdateReview(t *testing.T) {
	InsertTestUser(1)
	InsertTestUser(2)
	defer globalDb.Exec("TRUNCATE TABLE test_review_table, test_user_table CASCADE")
	reviewRepository := CreateReviewRepository(globalDb, "test_review_table")
	review := types.Review{
		ClientID:    1,
		RepetitorID: 2,
		Rating:      5,
		Comment:     "Good job!",
		CreatedAt:   time.Now(),
	}
	id, err := reviewRepository.InsertReview(review)
	if err != nil {
		t.Errorf("Failed to insert review: %v", err)
	}
	err = reviewRepository.UpdateReview(types.Review{
		ID:          id,
		ClientID:    1,
		RepetitorID: 2,
		Rating:      4,
		Comment:     "Bad job!",
	})
	if err != nil {
		t.Errorf("Failed to update review: %v", err)
	}
	resultReview, err := reviewRepository.GetReview(id)
	if err != nil {
		t.Errorf("Failed to get review: %v", err)
	}
	if resultReview.Rating != 4 {
		t.Errorf("Failed to get review: %v", resultReview.Rating)
	}
	if resultReview.Comment != "Bad job!" {
		t.Errorf("Failed to get review: %v", resultReview.Comment)
	}
	err = reviewRepository.UpdateReview(types.Review{
		ID:          id + 3,
		ClientID:    1,
		RepetitorID: 2,
		Rating:      5,
		Comment:     "Good job!",
	})
	if err == nil {
		t.Errorf("Try to update non-existent review")
	}
}
