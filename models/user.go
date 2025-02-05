package models

import (
	"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID        int       `orm:"pk;auto;column(id)" json:"id"`
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
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = u.Username
	claims["role"] = u.Role
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte("your-secret-key")) // Use environment variable in production
	if err != nil {
		return "", err
	}

	return tokenString, nil
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
		return []byte("your-secret-key"), nil
	})
	return token, err
}

func EnsureUsersTable() error {
	// Create table if it doesn't exist
	sql := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
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
