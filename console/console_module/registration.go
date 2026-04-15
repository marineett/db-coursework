package console_module

import (
	"bufio"
	"data_base_project/service_logic"
	"data_base_project/types"
	"fmt"
	"os"
	"strconv"
	"time"
)

func GetPersonalData() types.PersonalData {
	personalData := types.PersonalData{}
	fmt.Println("Enter first name:")
	fmt.Scanln(&personalData.FirstName)
	fmt.Println("Enter last name:")
	fmt.Scanln(&personalData.LastName)
	fmt.Println("Enter middle name:")
	fmt.Scanln(&personalData.MiddleName)
	fmt.Println("Enter email:")
	fmt.Scanln(&personalData.Email)
	fmt.Println("Enter telephone number:")
	fmt.Scanln(&personalData.TelephoneNumber)
	passportData := types.PassportData{}
	fmt.Println("Enter passport number:")
	fmt.Scanln(&passportData.PassportNumber)
	fmt.Println("Enter passport series:")
	fmt.Scanln(&passportData.PassportSeries)
	passportData.PassportDate = time.Now()
	fmt.Println("Enter passport issued by:")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := scanner.Text()
		passportData.PassportIssuedBy = input
	}
	personalData.PassportData = passportData
	return personalData
}

func GetAuthData() types.AuthData {
	authData := types.AuthData{}
	fmt.Println("Enter login:")
	fmt.Scanln(&authData.Login)
	fmt.Println("Enter password:")
	fmt.Scanln(&authData.Password)
	return authData
}

func RegisterClient(serviceModule *service_logic.ServiceModule) {
	personalData := GetPersonalData()
	authData := GetAuthData()
	initData := types.InitClientData{
		InitUserData: types.InitUserData{
			PersonalData: personalData,
			AuthData:     authData,
		},
	}
	err := serviceModule.ClientService.CreateClient(initData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Client registered successfully")
}

func RegisterRepetitor(serviceModule *service_logic.ServiceModule) {
	personalData := GetPersonalData()
	authData := GetAuthData()
	initData := types.InitRepetitorData{
		InitUserData: types.InitUserData{
			PersonalData: personalData,
			AuthData:     authData,
		},
	}
	err := serviceModule.RepetitorService.CreateRepetitor(initData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Repetitor registered successfully")
}

func RegisterModerator(serviceModule *service_logic.ServiceModule) {
	personalData := GetPersonalData()
	authData := GetAuthData()
	initData := types.InitModeratorData{
		InitUserData: types.InitUserData{
			PersonalData: personalData,
			AuthData:     authData,
		},
	}
	err := serviceModule.ModeratorService.CreateModerator(initData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Moderator registered successfully")
}

func RegisterAdmin(serviceModule *service_logic.ServiceModule) {
	personalData := GetPersonalData()
	authData := GetAuthData()
	fmt.Println("Enter admin salary:")
	var salaryStr string
	fmt.Scanln(&salaryStr)
	salary, err := strconv.ParseInt(salaryStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	initData := types.InitAdminData{
		InitUserData: types.InitUserData{
			PersonalData: personalData,
			AuthData:     authData,
		},
		Salary: salary,
	}
	err = serviceModule.AdminService.CreateAdmin(initData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Admin registered successfully")
}

func PrintRegistrationMenu() {
	fmt.Println("Registration:")
	fmt.Println("1. Register as a client")
	fmt.Println("2. Register as a repetitor")
	fmt.Println("3. Register as a moderator")
	fmt.Println("4. Register as an admin")
}

func UserRegistrationWork(serviceModule *service_logic.ServiceModule) {
	PrintRegistrationMenu()
	var choiceStr string
	fmt.Scanln(&choiceStr)
	choice, err := strconv.Atoi(choiceStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if choice < 1 || choice > 4 {
		fmt.Println("Invalid choice")
		return
	}
	switch choice {
	case 1:
		RegisterClient(serviceModule)
	case 2:
		RegisterRepetitor(serviceModule)
	case 3:
		RegisterModerator(serviceModule)
	case 4:
		RegisterAdmin(serviceModule)
	}
}
