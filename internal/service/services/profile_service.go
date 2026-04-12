// Package profileservice
package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/hihikaAAa/TrashProject/internal/domain/domain_errors"
	"github.com/hihikaAAa/TrashProject/internal/domain/person"
	"github.com/hihikaAAa/TrashProject/internal/domain/user"
	"github.com/hihikaAAa/TrashProject/internal/domain/worker"
	"github.com/hihikaAAa/TrashProject/internal/dto"
	errorswrapper "github.com/hihikaAAa/TrashProject/internal/service/errors_wrapper"
	"github.com/hihikaAAa/trash_project/internal/repositories"
	"github.com/hihikaAAa/trash_project/pkg/logger"
	"github.com/jackc/pgx/v5"
)

func (p *ProfileService) UpdatePersonality(ctx context.Context, input dto.ProfileInput) (dto.ProfileOutput, error) {
	switch input.Role {
	case worker.WorkerRole:
		wrk, err := p.WorkerRepo.GetByID(ctx, input.ID)
		if err != nil {
			return dto.ProfileOutput{}, errorswrapper.WrapRepoErr("p.workerRepo.GetByID", err)
		}
		newPrs, err := person.NewPerson(input.FirstName, input.Surname, input.LastName, input.Role)
		if err != nil {
			return dto.ProfileOutput{}, fmt.Errorf("new person: %w", err)
		}
		err = wrk.UpdateWorker(*newPrs)
		if err != nil {
			return dto.ProfileOutput{}, fmt.Errorf("update worker: %w", err)
		}
		_, err = p.WorkerRepo.UpdateWorker(ctx, wrk)
		if err != nil {
			return dto.ProfileOutput{}, errorswrapper.WrapRepoErr("p.workerRepo.UpdateWorker", err)
		}

	case user.UserRole:
		usr, err := p.UserRepo.GetByID(ctx, input.ID)
		if err != nil {
			return dto.ProfileOutput{}, errorswrapper.WrapRepoErr("p.UserRepo.GetByID", err)
		}
		newPrs, err := person.NewPerson(input.FirstName, input.Surname, input.LastName, input.Role)
		if err != nil {
			return dto.ProfileOutput{}, fmt.Errorf("new person: %w", err)
		}
		err = usr.UpdateUser(*newPrs)
		if err != nil {
			return dto.ProfileOutput{}, fmt.Errorf("update user: %w", err)
		}
		_, err = p.UserRepo.UpdateUser(ctx, usr)
		if err != nil {
			return dto.ProfileOutput{}, errorswrapper.WrapRepoErr("p.userRepo.UpdateUser", err)
		}
	}
	return dto.ProfileOutput{}, domainerrors.ErrWrongRole
}
