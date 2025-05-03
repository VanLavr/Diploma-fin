package utils

import (
	"context"

	"github.com/VanLavr/Diploma-fin/internal/domain/commands"
	"github.com/VanLavr/Diploma-fin/internal/domain/repositories"
	"github.com/VanLavr/Diploma-fin/utils/config"
	"github.com/VanLavr/Diploma-fin/utils/hasher"
)

func CreateAdmin(cfg *config.Config, repo repositories.Repository) error {
	pass := hasher.Hshr.Hash(cfg.AdminPass)

	if _, err := repo.CreateTeacher(context.TODO(), commands.CreateTeacher{
		FirstName:  "admin",
		LastName:   "admin",
		MiddleName: "admin",
		Email:      "admin",
		Password:   pass,
	}); err != nil {
		return err
	}
	return nil
}
