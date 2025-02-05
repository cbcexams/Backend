package controllers

import (
	"cbc-backend/models"
	"cbc-backend/utils"
	"encoding/json"
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"golang.org/x/crypto/bcrypt"
)

// UserController handles user-related operations
type UserController struct {
	beego.Controller
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SignupRequest represents the signup request body
type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Post handles user registration
func (c *UserController) Post() {
	fmt.Println("\n==================================================")
	fmt.Println("              User Registration                    ")
	fmt.Println("==================================================")

	// Ensure users table exists
	if err := models.EnsureUsersTable(); err != nil {
		utils.SendResponse(&c.Controller, false, "Failed to ensure users table exists", nil, err)
		return
	}

	var req SignupRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		utils.SendResponse(&c.Controller, false, "Invalid request body", nil, err)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.SendResponse(&c.Controller, false, "Failed to hash password", nil, err)
		return
	}

	// Create user
	user := &models.User{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	if err := models.CreateUser(user); err != nil {
		utils.SendResponse(&c.Controller, false, "Failed to create user", nil, err)
		return
	}

	utils.SendResponse(&c.Controller, true, "User created successfully", nil, nil)
}

// Login handles user authentication
func (c *UserController) Login() {
	fmt.Println("\n==================================================")
	fmt.Println("              User Login                          ")
	fmt.Println("==================================================")

	var req LoginRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		utils.SendResponse(&c.Controller, false, "Invalid request body", nil, err)
		return
	}

	// Get user from database
	user, err := models.GetUserByUsername(req.Username)
	if err != nil {
		utils.SendResponse(&c.Controller, false, "Invalid credentials", nil, err)
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.SendResponse(&c.Controller, false, "Invalid credentials", nil, err)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Username)
	if err != nil {
		utils.SendResponse(&c.Controller, false, "Failed to generate token", nil, err)
		return
	}

	utils.SendResponse(&c.Controller, true, "Login successful", map[string]string{
		"token": token,
	}, nil)
}

// Logout handles user logout
func (c *UserController) Logout() {
	// Get the token from Authorization header
	authHeader := c.Ctx.Input.Header("Authorization")
	if authHeader == "" {
		utils.SendResponse(&c.Controller, false, "No token provided", nil, nil)
		return
	}

	// For now, we'll just return success since JWT tokens are stateless
	utils.SendResponse(&c.Controller, true, "Logged out successfully", nil, nil)
}
