package models

import (
	"fiber-mongo-api/custom_errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/paemuri/brdoc"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Aluno struct {
	Id             primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Name           string             `json:"name,omitempty" validate:"required,min=3,max=32"`
	Serie          string             `json:"serie,omitempty" validate:"required"`
	Cpf            string             `json:"cpf,omitempty" validate:"required"`
	Email          string             `json:"email,omitempty" validate:"required,email"`
	DataNascimento string             `json:"data,omitempty" validate:"required"`
	Idade          int                `json:"idade,omitempty" bson:"idade,omitempty"`
	Endereco       Address
}

func (a Aluno) BodyParser(ctx *fiber.Ctx) (*Aluno, error) {
	payload := &Aluno{}
	err := ctx.BodyParser(payload)
	if err != nil {
		return nil, err
	}

	endereco, err := ToAdress(payload.Endereco.Cep)
	if err != nil {
		return nil, err
	}

	payload.Endereco = *endereco

	return payload, nil

}

func (a Aluno) QueryParamParser(ctx *fiber.Ctx) (string, error) {
	id := ctx.Params("alunoId")
	if id == "" {
		return "", custom_errors.ErrIDInvalido
	}
	return id, nil
}

func (a Aluno) IsValid() error {
	validate := validator.New()
	err := validate.Struct(a)
	if err != nil {
		return err
	}

	if !brdoc.IsCPF(a.Cpf) {
		return custom_errors.ErrCPFInvalido
	}

	return nil
}
