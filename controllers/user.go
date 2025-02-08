package controllers

import (
	"cbc-backend/models"
	"cbc-backend/utils"
	"encoding/json"
	"fmt"
	"os"
	"strings"

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
	Email    string `json:"email"`
}

// PromoteRequest represents the admin promotion request body
type PromoteRequest struct {
	SecretKey string `json:"secret_key"` // Additional security measure
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

	// Create user with default role
	user := &models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Role:     "user", // Default role - cannot be overridden from request
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

	// Generate JWT token with role
	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		utils.SendResponse(&c.Controller, false, "Failed to generate token", nil, err)
		return
	}

	// Return both token and user ID in response
	utils.SendResponse(&c.Controller, true, "Login successful", map[string]interface{}{
		"token":    token,
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
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

// ForgotPassword initiates password reset
func (c *UserController) ForgotPassword() {
	var request struct {
		Email string `json:"email"`
	}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &request); err != nil {
		utils.SendResponse(&c.Controller, false, "Invalid request", nil, err)
		return
	}

	// Generate reset token
	resetToken, err := models.CreatePasswordReset(request.Email)
	if err != nil {
		utils.SendResponse(&c.Controller, false, "Failed to create reset token", nil, err)
		return
	}

	// TODO: Send email with reset link
	// For now, return token in response (in production, send via email)
	utils.SendResponse(&c.Controller, true, "Reset token generated", map[string]string{
		"reset_token": resetToken,
	}, nil)
}

// ResetPassword handles the password reset
func (c *UserController) ResetPassword() {
	var request struct {
		ResetToken  string `json:"reset_token"`
		NewPassword string `json:"new_password"`
	}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &request); err != nil {
		utils.SendResponse(&c.Controller, false, "Invalid request", nil, err)
		return
	}

	if err := models.ResetPassword(request.ResetToken, request.NewPassword); err != nil {
		utils.SendResponse(&c.Controller, false, "Failed to reset password", nil, err)
		return
	}

	utils.SendResponse(&c.Controller, true, "Password reset successfully", nil, nil)
}

// Delete handles user deletion
func (c *UserController) Delete() {
	// Get user ID from URL - remove curly braces if present
	uid := strings.Trim(c.Ctx.Input.Param(":uid"), "{}")
	if uid == "" {
		utils.SendResponse(&c.Controller, false, "User ID is required", nil, nil)
		return
	}

	// Get claims from JWT token
	claims, err := utils.GetJWTClaims(c.Ctx)
	if err != nil {
		utils.SendResponse(&c.Controller, false, "Invalid token", nil, err)
		return
	}

	// Get user ID directly as string
	tokenUserID := claims["user_id"].(string)

	// Safely get role from claims
	var userRole string
	if role, ok := claims["role"]; ok && role != nil {
		userRole = role.(string)
	}

	// Only allow users to delete their own account or admin users
	if tokenUserID != uid && userRole != "admin" {
		utils.SendResponse(&c.Controller, false, "Unauthorized to delete this user", nil, nil)
		return
	}

	// Delete user
	if err := models.DeleteUserByID(uid); err != nil {
		utils.SendResponse(&c.Controller, false, "Failed to delete user", nil, err)
		return
	}

	utils.SendResponse(&c.Controller, true, "User deleted successfully", nil, nil)
}

// PromoteToAdmin promotes a user to admin role
func (c *UserController) PromoteToAdmin() {
	// Get user ID from URL
	uid := strings.Trim(c.Ctx.Input.Param(":uid"), "{}")
	if uid == "" {
		utils.SendResponse(&c.Controller, false, "User ID is required", nil, nil)
		return
	}

	// Get claims from JWT token
	claims, err := utils.GetJWTClaims(c.Ctx)
	if err != nil {
		utils.SendResponse(&c.Controller, false, "Invalid token", nil, err)
		return
	}

	// Only existing admins can promote others
	userRole, ok := claims["role"].(string)
	if !ok || userRole != "admin" {
		utils.SendResponse(&c.Controller, false, "Only administrators can promote users", nil, nil)
		return
	}

	// Verify secret key
	var req PromoteRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		utils.SendResponse(&c.Controller, false, "Invalid request body", nil, err)
		return
	}

	// Verify against environment variable or config
	if req.SecretKey != os.Getenv("ADMIN_PROMOTION_KEY") {
		utils.SendResponse(&c.Controller, false, "Invalid promotion key", nil, nil)
		return
	}

	// Update user role to admin
	if err := models.PromoteUserToAdmin(uid); err != nil {
		utils.SendResponse(&c.Controller, false, "Failed to promote user", nil, err)
		return
	}

	utils.SendResponse(&c.Controller, true, "User promoted to admin successfully", nil, nil)
}
