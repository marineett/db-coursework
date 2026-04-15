package data_base

import (
	"data_base_project/types"
	"encoding/json"
	"testing"
	"time"
)

func TestCreateResumeRepository(t *testing.T) {
	err := CreateResumeTable(globalDb, "test_resume_table", "test_user_table")
	if err != nil {
		t.Errorf("Failed to create resume table: %v", err)
	}
}

func TestInsertResume(t *testing.T) {
	InsertTestUser(1)
	defer globalDb.Exec("TRUNCATE TABLE test_resume_table, test_user_table CASCADE")
	resumeRepository := CreateResumeRepository(globalDb, "test_resume_table")
	resume := types.Resume{
		RepetitorID: 1,
		Title:       "Test Resume",
		Description: "Test Description",
		Prices:      map[string]int{"go": 1000, "typescript": 2000},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	id, err := resumeRepository.InsertResume(resume)
	if err != nil {
		t.Errorf("Failed to insert resume: %v", err)
	}
	resultResume := types.Resume{}
	var pricesJSON []byte
	err = globalDb.QueryRow("SELECT * FROM test_resume_table WHERE id = $1", id).Scan(&resultResume.ID, &resultResume.RepetitorID, &resultResume.Title, &resultResume.Description, &pricesJSON, &resultResume.CreatedAt, &resultResume.UpdatedAt)
	if err != nil {
		t.Errorf("Failed to get resume: %v", err)
	}
	err = json.Unmarshal(pricesJSON, &resultResume.Prices)
	if err != nil {
		t.Errorf("Failed to get resume: %v", err)
	}
	if resultResume.RepetitorID != 1 {
		t.Errorf("Failed to get resume: %v", resultResume.RepetitorID)
	}
	if resultResume.Title != "Test Resume" {
		t.Errorf("Failed to get resume: %v", resultResume.Title)
	}
	if resultResume.Description != "Test Description" {
		t.Errorf("Failed to get resume: %v", resultResume.Description)
	}
	if resultResume.Prices["go"] != 1000 {
		t.Errorf("Failed to get resume: %v", resultResume.Prices["go"])
	}
	if resultResume.Prices["typescript"] != 2000 {
		t.Errorf("Failed to get resume: %v", resultResume.Prices["typescript"])
	}

}

func TestGetResume(t *testing.T) {
	InsertTestUser(1)
	defer globalDb.Exec("TRUNCATE TABLE test_resume_table, test_user_table CASCADE")
	resumeRepository := CreateResumeRepository(globalDb, "test_resume_table")
	resume := types.Resume{
		RepetitorID: 1,
		Title:       "Test Resume",
		Description: "Test Description",
		Prices:      map[string]int{"go": 1000, "typescript": 2000},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	id, err := resumeRepository.InsertResume(resume)
	if err != nil {
		t.Errorf("Failed to insert resume: %v", err)
	}
	resultResume, err := resumeRepository.GetResume(id)
	if err != nil {
		t.Errorf("Failed to get resume: %v", err)
	}
	if resultResume.RepetitorID != 1 {
		t.Errorf("Failed to get resume: %v", resultResume.RepetitorID)
	}
	if resultResume.Title != "Test Resume" {
		t.Errorf("Failed to get resume: %v", resultResume.Title)
	}
	if resultResume.Description != "Test Description" {
		t.Errorf("Failed to get resume: %v", resultResume.Description)
	}
	if resultResume.Prices["go"] != 1000 {
		t.Errorf("Failed to get resume: %v", resultResume.Prices["go"])
	}
	if resultResume.Prices["typescript"] != 2000 {
		t.Errorf("Failed to get resume: %v", resultResume.Prices["typescript"])
	}
}

func TestUpdateResumeTitle(t *testing.T) {
	InsertTestUser(1)
	defer globalDb.Exec("TRUNCATE TABLE test_resume_table, test_user_table CASCADE")
	resumeRepository := CreateResumeRepository(globalDb, "test_resume_table")
	resume := types.Resume{
		RepetitorID: 1,
		Title:       "Test Resume",
		Description: "Test Description",
		Prices:      map[string]int{"go": 1000, "typescript": 2000},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	id, err := resumeRepository.InsertResume(resume)
	if err != nil {
		t.Errorf("Failed to insert resume: %v", err)
	}
	err = resumeRepository.UpdateResumeTitle(id, "New Title", time.Now())
	if err != nil {
		t.Errorf("Failed to update resume title: %v", err)
	}
	resultResume, err := resumeRepository.GetResume(id)
	if err != nil {
		t.Errorf("Failed to get resume: %v", err)
	}
	if resultResume.Title != "New Title" {
		t.Errorf("Failed to get resume: %v", resultResume.Title)
	}
	err = resumeRepository.UpdateResumeTitle(id+1, "New Title", time.Now())
	if err == nil {
		t.Errorf("Try to update non-existent resume")
	}
}

func TestUpdateResumeDescription(t *testing.T) {
	InsertTestUser(1)
	defer globalDb.Exec("TRUNCATE TABLE test_resume_table, test_user_table CASCADE")
	resumeRepository := CreateResumeRepository(globalDb, "test_resume_table")
	resume := types.Resume{
		RepetitorID: 1,
		Title:       "Test Resume",
		Description: "Test Description",
		Prices:      map[string]int{"go": 1000, "typescript": 2000},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	id, err := resumeRepository.InsertResume(resume)
	if err != nil {
		t.Errorf("Failed to insert resume: %v", err)
	}
	err = resumeRepository.UpdateResumeDescription(id, "New Description", time.Now())
	if err != nil {
		t.Errorf("Failed to update resume description: %v", err)
	}
	resultResume, err := resumeRepository.GetResume(id)
	if err != nil {
		t.Errorf("Failed to get resume: %v", err)
	}
	if resultResume.Description != "New Description" {
		t.Errorf("Failed to get resume: %v", resultResume.Description)
	}
	err = resumeRepository.UpdateResumeDescription(id+1, "New Description", time.Now())
	if err == nil {
		t.Errorf("Try to update non-existent resume")
	}
}

func TestUpdateResumePrices(t *testing.T) {
	InsertTestUser(1)
	defer globalDb.Exec("TRUNCATE TABLE test_resume_table, test_user_table CASCADE")
	resumeRepository := CreateResumeRepository(globalDb, "test_resume_table")
	resume := types.Resume{
		RepetitorID: 1,
		Title:       "Test Resume",
		Description: "Test Description",
		Prices:      map[string]int{"go": 1000, "typescript": 2000},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	id, err := resumeRepository.InsertResume(resume)
	if err != nil {
		t.Errorf("Failed to insert resume: %v", err)
	}
	err = resumeRepository.UpdateResumePrices(id, map[string]int{"C++": 2000, "typescript": 1000}, time.Now())
	if err != nil {
		t.Errorf("Failed to update resume prices: %v", err)
	}
	resultResume, err := resumeRepository.GetResume(id)
	if err != nil {
		t.Errorf("Failed to get resume: %v", err)
	}
	if resultResume.Prices["C++"] != 2000 {
		t.Errorf("Failed to get resume: %v", resultResume.Prices["C++"])
	}
	if resultResume.Prices["typescript"] != 1000 {
		t.Errorf("Failed to get resume: %v", resultResume.Prices["typescript"])
	}
	if _, ok := resultResume.Prices["go"]; ok {
		t.Errorf("Got price for non-existent language")
	}
	err = resumeRepository.UpdateResumePrices(id+1, map[string]int{"C++": 2000, "typescript": 1000}, time.Now())
	if err == nil {
		t.Errorf("Try to update non-existent resume")
	}
}

func TestDeleteResume(t *testing.T) {
	InsertTestUser(1)
	defer globalDb.Exec("TRUNCATE TABLE test_resume_table, test_user_table CASCADE")
	resumeRepository := CreateResumeRepository(globalDb, "test_resume_table")
	resume := types.Resume{
		RepetitorID: 1,
		Title:       "Test Resume",
		Description: "Test Description",
		Prices:      map[string]int{"go": 1000, "typescript": 2000},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	id, err := resumeRepository.InsertResume(resume)
	if err != nil {
		t.Errorf("Failed to insert resume: %v", err)
	}
	err = resumeRepository.DeleteResume(id)
	if err != nil {
		t.Errorf("Failed to delete resume: %v", err)
	}
	_, err = resumeRepository.GetResume(id)
	if err == nil {
		t.Errorf("Resume not deleted")
	}
	err = resumeRepository.DeleteResume(id + 1)
	if err == nil {
		t.Errorf("Try to delete non-existent resume")
	}
}
