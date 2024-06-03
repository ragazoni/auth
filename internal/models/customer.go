package models

import (
	"auth/internal/errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type Customer struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name,omitempty"`
	CPF       string             `json:"cpf" bson:"cpf,omitempty"`
	BirthDate string             `json:"birthdate" bson:"birthdate,omitempty"`
	Email     string             `json:"email" bson:"email,omitempty"`
	Password  string             `json:"password" bson:"password,omitempty"`
	Type      string             `json:"type" bson:"type,omitempty"`
	IsAdmin   bool               `json:"isadmin" bson:"isadmin,omitempty"`
}

func (c Customer) Validate() error {
	errList := []string{}
	if c.Name == "" {
		errList = append(errList, "nome completo")
	}
	if c.CPF == "" {
		errList = append(errList, "documento")
	}
	if c.BirthDate == "" {
		errList = append(errList, "data nascimento")
	}
	if c.Email == "" {
		errList = append(errList, "email")
	}
	if c.Password == "" {
		errList = append(errList, "senha")
	}
	if len(errList) > 0 {
		return errors.NewRequiredFieldsErrorList(errList)

	}

	return nil
}

func (c *Customer) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	c.Password = string(hashedPassword)
	return nil
}

func (c *Customer) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(password))
	return err == nil
}
