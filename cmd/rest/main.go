package main

import (
	"context"

	"github.com/VanLavr/Diploma-fin/internal/controllers/rest"
	"github.com/VanLavr/Diploma-fin/internal/infrastructure/postgres"
	application "github.com/VanLavr/Diploma-fin/internal/services"
	"github.com/VanLavr/Diploma-fin/utils"
	"github.com/VanLavr/Diploma-fin/utils/auth"
	"github.com/VanLavr/Diploma-fin/utils/config"
	"github.com/VanLavr/Diploma-fin/utils/errors"
)

func main() {
	// --- for student
	// create GET /allDebts endpoint +
	// create POST /notification endpoint +

	// --- for teacher
	// create GET /allExams endpoint +
	// create POST /setDate endpoint +

	// --- infra
	// create db schema
	// run all in docker containers

	cfg, err := config.ReadConfig()
	errors.FatalOnError(err)

	repository := postgres.NewRepository(cfg)
	errors.FatalOnError(utils.CreateAdmin(cfg, repository))

	studentApp := application.NewStudentUsecase(repository)
	teacherApp := application.NewTeacherUsecase(repository)
	examApp := application.NewExamUsecase(repository)
	groupApp := application.NewGroupUsecase(repository)
	fileApp := application.NewFileUsecase(repository)

	server := rest.NewServer(
		cfg,
		rest.NewStudentHandler(studentApp),
		rest.NewTeacherHandler(teacherApp, studentApp),
		rest.NewAuthHandler(teacherApp, studentApp, auth.NewAuthMiddleware(cfg, repository)),
		rest.NewExamHandler(examApp),
		rest.NewGroupHandler(groupApp),
		rest.NewFileHandler(fileApp),
	)

	errors.FatalOnError(server.Start(context.Background()))
}
