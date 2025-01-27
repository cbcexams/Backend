package controllers

import (
	"cbc-backend/models"
	"cbc-backend/utils"
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
)

// UserController handles user-related operations
type UserController struct {
	beego.Controller
}

// Post handles user registration
func (u *UserController) Post() {
	fmt.Println("\n==================================================")
	fmt.Println("              User Registration                    ")
	fmt.Println("==================================================")

	// Get registration data
	username := u.GetString("username")
	password := u.GetString("password")
	email := u.GetString("email")

	fmt.Printf("Registering user: %s, Email: %s\n", username, email)

	// Create new user
	user := &models.User{
		Username: username,
		Password: password,
		Email:    email,
		Role:     "user", // Default role
	}

	// Add user to database
	if err := models.AddUser(user); err != nil {
		fmt.Printf("Error adding user: %v\n", err)
		utils.SendResponse(&u.Controller, false, "", nil, err)
		return
	}

	// Generate JWT token
	token, err := user.GenerateToken()
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		utils.SendResponse(&u.Controller, false, "", nil, err)
		return
	}

	utils.SendResponse(&u.Controller, true, "User registered successfully", map[string]string{
		"token": token,
	}, nil)
}

// Login handles user authentication
func (u *UserController) Login() {
	fmt.Println("\n==================================================")
	fmt.Println("              User Login                          ")
	fmt.Println("==================================================")

	// Get login credentials
	username := u.GetString("username")
	password := u.GetString("password")

	fmt.Printf("Login attempt for user: %s\n", username)

	// Get user from database
	user, err := models.GetUser(username)
	if err != nil {
		fmt.Printf("Error finding user: %v\n", err)
		utils.SendResponse(&u.Controller, false, "", nil, fmt.Errorf("invalid credentials"))
		return
	}

	// Validate password
	if err := user.ValidatePassword(password); err != nil {
		fmt.Printf("Invalid password for user %s: %v\n", username, err)
		utils.SendResponse(&u.Controller, false, "", nil, fmt.Errorf("invalid credentials"))
		return
	}

	// Generate JWT token
	token, err := user.GenerateToken()
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		utils.SendResponse(&u.Controller, false, "", nil, err)
		return
	}

	utils.SendResponse(&u.Controller, true, "Login successful", map[string]string{
		"token": token,
	}, nil)
}
