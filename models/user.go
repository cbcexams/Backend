package models

import (
	"cbc-backend/config"
	"cbc-backend/utils"
	"errors"
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID        string    `orm:"pk;size(36);column(id)" json:"id"`
	Username  string    `orm:"unique;size(128)" json:"username"`
	Password  string    `orm:"size(128)" json:"-"`
	Email     string    `orm:"size(128);unique" json:"email"`
	Role      string    `orm:"size(20)" json:"role"` // admin, teacher, student
	CreatedAt time.Time `orm:"auto_now_add;type(timestamp)" json:"created_at"`
}

// TableName specifies the database table name
func (u *User) TableName() string {
	return "users"
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) GenerateToken() (string, error) {
	claims := jwt.MapClaims{
		"user_id":  u.ID,
		"username": u.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JWTSecret)
}

// GetUserByUsername retrieves a user by username
func GetUserByUsername(username string) (*User, error) {
	o := orm.NewOrm()
	user := &User{Username: username}
	err := o.Read(user, "Username")
	return user, err
}

// CreateUser creates a new user
func CreateUser(user *User) error {
	// Generate UUID for new user
	user.ID = uuid.New().String()

	o := orm.NewOrm()
	_, err := o.Insert(user)
	return err
}

func UpdateUser(user *User) error {
	o := orm.NewOrm()
	_, err := o.Update(user)
	return err
}

func DeleteUser(username string) error {
	o := orm.NewOrm()
	_, err := o.Delete(&User{Username: username})
	return err
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return config.JWTSecret, nil
	})
	return token, err
}

func EnsureUsersTable() error {
	sql := `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(36) PRIMARY KEY,
		username VARCHAR(128) UNIQUE NOT NULL,
		password VARCHAR(128) NOT NULL,
		email VARCHAR(128) UNIQUE,
		role VARCHAR(20) DEFAULT 'user',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	)`

	o := orm.NewOrm()
	_, err := o.Raw(sql).Exec()
	return err
}

// PasswordReset represents a password reset request
type PasswordReset struct {
	ID        int    `orm:"pk;auto;column(id)"`
	UserID    string `orm:"column(user_id);size(36)"`
	Token     string `orm:"size(100)"`
	ExpiresAt time.Time
	Used      bool      `orm:"default(false)"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
}

func CreatePasswordReset(email string) (string, error) {
	o := orm.NewOrm()

	// Find user by email
	var user User
	if err := o.QueryTable("users").Filter("email", email).One(&user); err != nil {
		return "", fmt.Errorf("user not found")
	}

	// Generate reset token
	token := utils.GenerateRandomString(32)

	// Create reset record
	reset := &PasswordReset{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour), // Token valid for 24 hours
	}

	_, err := o.Insert(reset)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ResetPassword(token, newPassword string) error {
	o := orm.NewOrm()

	// Find valid reset token
	var reset PasswordReset
	err := o.QueryTable("password_resets").
		Filter("token", token).
		Filter("expires_at__gt", time.Now()).
		Filter("used", false).
		One(&reset)

	if err != nil {
		return fmt.Errorf("invalid or expired reset token")
	}

	// Get user
	var user User
	if err := o.QueryTable("users").Filter("id", reset.UserID).One(&user); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	user.Password = string(hashedPassword)
	if _, err := o.Update(&user, "Password"); err != nil {
		return err
	}

	// Mark token as used
	reset.Used = true
	_, err = o.Update(&reset, "Used")
	return err
}

func EnsurePasswordResetTable() error {
	sql := `
	CREATE TABLE IF NOT EXISTS password_resets (
		id SERIAL PRIMARY KEY,
		user_id VARCHAR(36) REFERENCES users(id),
		token VARCHAR(100) NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		used BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	)`

	o := orm.NewOrm()
	_, err := o.Raw(sql).Exec()
	return err
}

// DeleteUserByID deletes a user by their ID
func DeleteUserByID(userID string) error {
	o := orm.NewOrm()

	// No need to convert to int anymore since we're using UUID
	// Delete password resets first using raw SQL
	_, err := o.Raw("DELETE FROM password_resets WHERE user_id = ?", userID).Exec()
	if err != nil {
		fmt.Printf("Warning: Failed to delete password resets: %v\n", err)
	}

	// Delete user
	_, err = o.Delete(&User{ID: userID})
	if err != nil {
		return err
	}

	return nil
}

// PromoteUserToAdmin promotes a user to admin role
func PromoteUserToAdmin(userID string) error {
	o := orm.NewOrm()

	// Get user
	user := &User{ID: userID}
	if err := o.Read(user); err != nil {
		return err
	}

	// Update role
	user.Role = "admin"
	_, err := o.Update(user, "Role")
	return err
}
