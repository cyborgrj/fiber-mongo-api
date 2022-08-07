package routes

import (
	"fiber-mongo-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func AlunosRoute(app *fiber.App) {
	//Todas as rotas relacionadas aos alunos virão aqui
	app.Post("/aluno", controllers.CreateAluno)

	app.Get("/aluno/:userId", controllers.GetAluno)

	app.Get("/alunos", controllers.GetAlunos)
}
