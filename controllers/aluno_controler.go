package controllers

import (
	"context"
	"fiber-mongo-api/configs"
	"fiber-mongo-api/models"
	"fiber-mongo-api/responses"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func age(birthdate, today time.Time) int {
	today = today.In(birthdate.Location())
	ty, tm, td := today.Date()
	today = time.Date(ty, tm, td, 0, 0, 0, 0, time.UTC)
	by, bm, bd := birthdate.Date()
	birthdate = time.Date(by, bm, bd, 0, 0, 0, 0, time.UTC)
	if today.Before(birthdate) {
		return 0
	}
	age := ty - by
	anniversary := birthdate.AddDate(age, 0, 0)
	if anniversary.After(today) {
		age--
	}
	return age
}

var alunoCollection *mongo.Collection = configs.GetCollection(configs.DB, "alunos")
var validate = validator.New()

func CreateAluno(ctx *fiber.Ctx) error {
	cntx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var aluno models.Aluno
	defer cancel()

	//validar o body do request
	if err := ctx.BodyParser(&aluno); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(responses.AlunoResponse{Status: http.StatusBadRequest, Message: "Erro", Data: &fiber.Map{"data": err.Error()}})
	}

	//usar biblioteca validator para validar os campos requeridos
	if validationErr := validate.Struct(&aluno); validationErr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(responses.AlunoResponse{Status: http.StatusBadRequest, Message: "Erro", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newAluno := models.Aluno{
		Id:             primitive.NewObjectID(),
		Name:           aluno.Name,
		Serie:          aluno.Serie,
		Cpf:            aluno.Cpf,
		DataNascimento: aluno.DataNascimento,
	}

	result, err := alunoCollection.InsertOne(cntx, newAluno)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(responses.AlunoResponse{Status: http.StatusInternalServerError, Message: "Erro", Data: &fiber.Map{"data": err.Error()}})
	}
	return ctx.Status(http.StatusCreated).JSON(responses.AlunoResponse{Status: http.StatusCreated, Message: "Sucesso", Data: &fiber.Map{"data": result}})
}

func GetAluno(ctx *fiber.Ctx) error {
	cntx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	alunoId := ctx.Params("alunoId")
	var aluno models.Aluno
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(alunoId)

	err := alunoCollection.FindOne(cntx, bson.M{"id": objId}).Decode(&aluno)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(responses.AlunoResponse{Status: http.StatusInternalServerError, Message: "Erro", Data: &fiber.Map{"data": err.Error()}})
	}

	return ctx.Status(http.StatusOK).JSON(responses.AlunoResponse{Status: http.StatusOK, Message: "Sucesso", Data: &fiber.Map{"data": aluno}})

}

func GetAlunos(ctx *fiber.Ctx) error {
	cntx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var alunos []models.Aluno
	defer cancel()

	results, err := alunoCollection.Find(cntx, bson.M{})

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(responses.AlunoResponse{Status: http.StatusInternalServerError, Message: "Erro", Data: &fiber.Map{"data": err.Error()}})
	}

	//iterando os dados do banco
	defer results.Close(cntx)
	for results.Next(cntx) {
		var singleAluno models.Aluno
		if err = results.Decode(&singleAluno); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(responses.AlunoResponse{Status: http.StatusInternalServerError, Message: "Erro", Data: &fiber.Map{"data": err.Error()}})
		}

		hoje := time.Now()
		datanasc, _ := time.Parse("02/01/2006", singleAluno.DataNascimento)
		singleAluno.Idade = age(datanasc, hoje)
		alunos = append(alunos, singleAluno)

	}
	return ctx.Status(http.StatusOK).JSON(responses.AlunoResponse{Status: http.StatusOK, Message: "Sucesso", Data: &fiber.Map{"data": alunos}})
}
