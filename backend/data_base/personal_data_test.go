package data_base

import (
	"data_base_project/types"
	"log"
	"testing"
	"time"
)

func TestCreatePersonalDataRepository(t *testing.T) {
	personalDataRepository := CreatePersonalDataRepository(globalDb, "test_personal_data_table")
	if personalDataRepository == nil {
		t.Errorf("Failed to create personal data repository")
	}
}

func TestInsertPersonalData(t *testing.T) {
	personalData := types.PersonalData{
		TelephoneNumber: "+88005553535",
		Email:           "test@example.com",
		PassportData: types.PassportData{
			PassportNumber:   "1234567890",
			PassportDate:     time.Now(),
			PassportSeries:   "1024",
			PassportIssuedBy: "test",
		},
		FirstName:  "Jhon",
		LastName:   "Doe",
		MiddleName: "Jhonovich",
	}
	personalDataRepository := CreatePersonalDataRepository(globalDb, "test_personal_data_table")
	lastInsertedID, err := personalDataRepository.InsertPersonalData(personalData)
	if err != nil {
		log.Fatalf("Error inserting personal data: %v", err)
	}
	defer globalDb.Exec("TRUNCATE TABLE test_user_table, test_auth_table, test_personal_data_table CASCADE")
	getResult, err := personalDataRepository.GetPersonalData(lastInsertedID)
	if err != nil {
		log.Fatalf("Error getting personal data: %v", err)
	}
	if getResult.ID != lastInsertedID {
		t.Errorf("Expected ID %d, got %d", lastInsertedID, getResult.ID)
	}
	if getResult.TelephoneNumber != "+88005553535" {
		t.Errorf("Expected TelephoneNumber %s, got %s", "+88005553535", getResult.TelephoneNumber)
	}
	if getResult.Email != "test@example.com" {
		t.Errorf("Expected Email %s, got %s", "test@example.com", getResult.Email)
	}
	if getResult.PassportNumber != "1234567890" {
		t.Errorf("Expected PassportNumber %s, got %s", "1234567890", getResult.PassportNumber)
	}
	if getResult.PassportSeries != "1024" {
		t.Errorf("Expected PassportSeries %s, got %s", "1024", getResult.PassportSeries)
	}
	if getResult.FirstName != "Jhon" {
		t.Errorf("Expected FirstName %s, got %s", "Jhon", getResult.FirstName)
	}
	if getResult.LastName != "Doe" {
		t.Errorf("Expected LastName %s, got %s", "Doe", getResult.LastName)
	}
	if getResult.MiddleName != "Jhonovich" {
		t.Errorf("Expected MiddleName %s, got %s", "Jhonovich", getResult.MiddleName)
	}
}

func TestGetPersonalData(t *testing.T) {
	personalDataRepository := CreatePersonalDataRepository(globalDb, "test_personal_data_table")
	personalData := types.PersonalData{
		TelephoneNumber: "+88005553535",
		Email:           "test@example.com",
		PassportData: types.PassportData{
			PassportNumber:   "1234567890",
			PassportDate:     time.Now(),
			PassportSeries:   "1024",
			PassportIssuedBy: "test",
		},
		FirstName:  "Jhon",
		LastName:   "Doe",
		MiddleName: "Jhonovich",
	}
	lastInsertedID, err := personalDataRepository.InsertPersonalData(personalData)
	if err != nil {
		log.Fatalf("Error inserting personal data: %v", err)
	}
	defer globalDb.Exec("TRUNCATE TABLE test_user_table, test_auth_table, test_personal_data_table CASCADE")
	getResult, err := personalDataRepository.GetPersonalData(lastInsertedID)
	if err != nil {
		log.Fatalf("Error getting personal data: %v", err)
	}
	if getResult.ID != lastInsertedID {
		t.Errorf("Expected ID %d, got %d", lastInsertedID, getResult.ID)
	}
	if getResult.TelephoneNumber != "+88005553535" {
		t.Errorf("Expected TelephoneNumber %s, got %s", "+88005553535", getResult.TelephoneNumber)
	}
	if getResult.Email != "test@example.com" {
		t.Errorf("Expected Email %s, got %s", "test@example.com", getResult.Email)
	}
	if getResult.PassportNumber != "1234567890" {
		t.Errorf("Expected PassportNumber %s, got %s", "1234567890", getResult.PassportNumber)
	}
	if getResult.PassportSeries != "1024" {
		t.Errorf("Expected PassportSeries %s, got %s", "1024", getResult.PassportSeries)
	}
	if getResult.FirstName != "Jhon" {
		t.Errorf("Expected FirstName %s, got %s", "Jhon", getResult.FirstName)
	}
	if getResult.LastName != "Doe" {
		t.Errorf("Expected LastName %s, got %s", "Doe", getResult.LastName)
	}
	if getResult.MiddleName != "Jhonovich" {
		t.Errorf("Expected MiddleName %s, got %s", "Jhonovich", getResult.MiddleName)
	}

}

func TestUpdatePersonalData(t *testing.T) {
	personalDataRepository := CreatePersonalDataRepository(globalDb, "test_personal_data_table")
	personalData := types.PersonalData{
		TelephoneNumber: "+88005553535",
		Email:           "test@example.com",
		PassportData: types.PassportData{
			PassportNumber:   "1234567890",
			PassportDate:     time.Now(),
			PassportSeries:   "1024",
			PassportIssuedBy: "test",
		},
		FirstName:  "Jhon",
		LastName:   "Doe",
		MiddleName: "Jhonovich",
	}
	lastInsertedID, err := personalDataRepository.InsertPersonalData(personalData)
	if err != nil {
		log.Fatalf("Error inserting personal data: %v", err)
	}
	defer globalDb.Exec("TRUNCATE TABLE test_user_table, test_auth_table, test_personal_data_table CASCADE")
	updatedPersonalData := types.PersonalData{
		TelephoneNumber: "+88005553536",
		Email:           "test2@example.com",
		PassportData: types.PassportData{
			PassportNumber:   "1234567891",
			PassportDate:     time.Now(),
			PassportSeries:   "1025",
			PassportIssuedBy: "test2",
		},
		FirstName:  "Jhon",
		LastName:   "Doe",
		MiddleName: "Jhonovich",
	}
	err = personalDataRepository.UpdatePersonalData(lastInsertedID, updatedPersonalData)
	if err != nil {
		log.Fatalf("Error updating personal data: %v", err)
	}
	getResult, err := personalDataRepository.GetPersonalData(lastInsertedID)
	if err != nil {
		log.Fatalf("Error getting personal data: %v", err)
	}
	if getResult.TelephoneNumber != "+88005553536" {
		t.Errorf("Expected TelephoneNumber %s, got %s", "+88005553536", getResult.TelephoneNumber)
	}
	if getResult.Email != "test2@example.com" {
		t.Errorf("Expected Email %s, got %s", "test2@example.com", getResult.Email)
	}
	if getResult.PassportNumber != "1234567891" {
		t.Errorf("Expected PassportNumber %s, got %s", "1234567891", getResult.PassportNumber)
	}
	if getResult.PassportSeries != "1025" {
		t.Errorf("Expected PassportSeries %s, got %s", "1025", getResult.PassportSeries)
	}
	if getResult.FirstName != "Jhon" {
		t.Errorf("Expected FirstName %s, got %s", "Jhon", getResult.FirstName)
	}
	if getResult.LastName != "Doe" {
		t.Errorf("Expected LastName %s, got %s", "Doe", getResult.LastName)
	}
	if getResult.MiddleName != "Jhonovich" {
		t.Errorf("Expected MiddleName %s, got %s", "Jhonovich", getResult.MiddleName)
	}
}
