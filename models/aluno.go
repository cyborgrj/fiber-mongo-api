package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Aluno struct {
	Id             primitive.ObjectID `json:"id,omitempty"`
	Name           string             `json:"name,omitempty" validate:"required,min=3,max=32"`
	Serie          string             `json:"serie,omitempty" validate:"required"`
	Cpf            string             `json:"cpf,omitempty" validate:"required"`
	Email          string             `json:"email,omitempty" validate:"required,min=6,max=32"`
	DataNascimento string             `json:"data,omitempty" validate:"required"`
	Idade          int
}
