package models

import (
	"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int       `orm:"pk;auto" json:"id"`
	Username  string    `orm:"size(128);unique" json:"username"`
	Password  string    `orm:"size(128)" json:"-"`
	Email     string    `orm:"size(128);unique" json:"email"`
	Role      string    `orm:"size(20)" json:"role"` // admin, teacher, student
	CreatedAt time.Time `orm:"auto_now_add;type(timestamp)" json:"created_at"`
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

func AddUser(user *User) error {
	if err := user.HashPassword(); err != nil {
		return err
	}

	o := orm.NewOrm()
	_, err := o.Insert(user)
	return err
}

func GetUser(username string) (*User, error) {
	o := orm.NewOrm()
	user := User{Username: username}
	err := o.Read(&user, "Username")
	if err != nil {
		return nil, err
	}
	return &user, nil
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
