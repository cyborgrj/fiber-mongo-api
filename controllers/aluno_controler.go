package controllers

import (
	"context"
	"fiber-mongo-api/configs"
	"fiber-mongo-api/custom_errors"
	"fiber-mongo-api/models"
	"fiber-mongo-api/responses"
	"net/http"
	"time"

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

func CreateAluno(ctx *fiber.Ctx) error {

	aluno, err := models.Aluno{}.BodyParser(ctx)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(err.Error())
	}

	aluno.Id = primitive.NewObjectID()

	err = aluno.IsValid()
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(err.Error())
	}

	err = alunoCollection.FindOne(ctx.UserContext(), bson.M{"cpf": aluno.Cpf}).Decode(&aluno)
	if err != mongo.ErrNoDocuments {
		return ctx.Status(http.StatusBadRequest).JSON(custom_errors.ErrCPFJaCadastrado.Error())
	}

	_, err = alunoCollection.InsertOne(ctx.UserContext(), aluno)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(aluno)
}

func GetAluno(ctx *fiber.Ctx) error {
	alunoId, err := models.Aluno{}.QueryParamParser(ctx)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(err.Error())
	}
	var aluno models.Aluno

	objId, _ := primitive.ObjectIDFromHex(alunoId)

	err = alunoCollection.FindOne(ctx.UserContext(), bson.M{"_id": objId}).Decode(&aluno)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(aluno)

}

func GetAlunos(ctx *fiber.Ctx) error {
	cntx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var alunos []models.Aluno
	defer cancel()

	results, err := alunoCollection.Find(cntx, bson.M{})

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	//iterando os dados do banco
	defer results.Close(cntx)
	for results.Next(cntx) {
		var singleAluno models.Aluno
		if err = results.Decode(&singleAluno); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(err.Error())
		}

		// Calcular data de nascimento
		hoje := time.Now()
		datanasc, _ := time.Parse("02/01/2006", singleAluno.DataNascimento)
		singleAluno.Idade = age(datanasc, hoje)
		alunos = append(alunos, singleAluno)
	}
	return ctx.Status(http.StatusOK).JSON(responses.AlunoResponse{Status: http.StatusOK, Message: "Sucesso", Data: &fiber.Map{"data": alunos}})
}

func DeleteAluno(ctx *fiber.Ctx) error {
	alunoId, err := models.Aluno{}.QueryParamParser(ctx)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(err.Error())
	}
	var aluno models.Aluno

	objId, _ := primitive.ObjectIDFromHex(alunoId)

	err = alunoCollection.FindOneAndDelete(ctx.UserContext(), bson.M{"_id": objId}).Decode(&aluno)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(aluno)

}

func UpdateAluno(ctx *fiber.Ctx) error {

	aluno, err := models.Aluno{}.BodyParser(ctx)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(err.Error())
	}

	alunoId, errAluno := models.Aluno{}.QueryParamParser(ctx)
	if errAluno != nil {
		return ctx.Status(http.StatusBadRequest).JSON(err.Error())
	}

	objId, _ := primitive.ObjectIDFromHex(alunoId)

	err = aluno.IsValid()
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(err.Error())
	}

	update := bson.M{
		"name":                 aluno.Name,
		"serie":                aluno.Serie,
		"data":                 aluno.DataNascimento,
		"cpf":                  aluno.Cpf,
		"email":                aluno.Email,
		"endereco.logradouro":  aluno.Endereco.Logradouro,
		"endereco.cep":         aluno.Endereco.Cep,
		"endereco.bairro":      aluno.Endereco.Bairro,
		"endereco.localidade":  aluno.Endereco.Cidade,
		"endereco.uf":          aluno.Endereco.Uf,
		"endereco.complemento": aluno.Endereco.Complemento,
	}
	err = alunoCollection.FindOneAndUpdate(ctx.UserContext(), bson.M{"_id": objId}, bson.M{"$set": update}).Err()
	if err != nil {
		return err
	}

	err = alunoCollection.FindOne(ctx.UserContext(), bson.M{"_id": objId}).Decode(&aluno)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(&aluno)

}
