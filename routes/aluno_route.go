package routes

import (
	"fiber-mongo-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func AlunosRoute(app *fiber.App) {
	//Todas as rotas relacionadas aos alunos virão aqui
	app.Post("/aluno", controllers.CreateAluno)

	app.Get("/aluno/:alunoId", controllers.GetAluno)

	app.Get("/alunos", controllers.GetAlunos)

	app.Delete("/aluno/:alunoId", controllers.DeleteAluno)

	app.Put("/aluno/:alunoId", controllers.UpdateAluno)
}
