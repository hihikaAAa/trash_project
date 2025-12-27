// Package profileservice
package profileservice

import (
	"context"
	"database/sql"
	"fmt"

	domainerrors "github.com/hihikaAAa/TrashProject/internal/domain/domain_errors"
	"github.com/hihikaAAa/TrashProject/internal/domain/person"
	"github.com/hihikaAAa/TrashProject/internal/domain/user"
	"github.com/hihikaAAa/TrashProject/internal/domain/worker"
	"github.com/hihikaAAa/TrashProject/internal/dto"
	addressrepo "github.com/hihikaAAa/TrashProject/internal/postgres/address_repo"
	userrepo "github.com/hihikaAAa/TrashProject/internal/postgres/user_repo"
	workerrepo "github.com/hihikaAAa/TrashProject/internal/postgres/worker_repo"
	errorswrapper "github.com/hihikaAAa/TrashProject/internal/service/errors_wrapper"
)

type ProfileService struct{
	WorkerRepo workerrepo.WorkerRepository
	AddressRepo addressrepo.AddressRepository
	UserRepo    userrepo.UserRepository
}

func NewProfileService(db *sql.DB) *ProfileService{
	wRepo := workerrepo.NewWorkerRepository(db)
	aRepo := addressrepo.NewAddressRepository(db)
	uRepo := userrepo.NewUserRepository(db)
	return &ProfileService{
		WorkerRepo: wRepo,
		AddressRepo: aRepo,
		UserRepo: uRepo,
	}
}

func (p *ProfileService) UpdatePersonality(ctx context.Context, input dto.ProfileInput)(dto.ProfileOutput, error){
	switch input.Role{
		case worker.WorkerRole:
			wrk, err := p.WorkerRepo.GetByID(ctx,input.ID)
			if err != nil{
				return dto.ProfileOutput{}, errorswrapper.WrapRepoErr("p.workerRepo.GetByID", err)
			}
			newPrs, err := person.NewPerson(input.FirstName, input.Surname, input.LastName, input.Role)
			if err != nil{
				return dto.ProfileOutput{}, fmt.Errorf("new person: %w", err)
			}
			err = wrk.UpdateWorker(*newPrs)
			if err != nil{
				return dto.ProfileOutput{}, fmt.Errorf("update worker: %w", err)
			}
			_, err = p.WorkerRepo.UpdateWorker(ctx,wrk)
			if err != nil{
				return dto.ProfileOutput{}, errorswrapper.WrapRepoErr("p.workerRepo.UpdateWorker", err)
			}

		case user.UserRole:
			usr, err := p.UserRepo.GetByID(ctx,input.ID)
			if err != nil{
				return dto.ProfileOutput{}, errorswrapper.WrapRepoErr("p.UserRepo.GetByID", err)
			}
			newPrs, err := person.NewPerson(input.FirstName, input.Surname, input.LastName, input.Role)
			if err != nil{
				return dto.ProfileOutput{}, fmt.Errorf("new person: %w", err)
			}
			err = usr.UpdateUser(*newPrs)
			if err != nil{
				return dto.ProfileOutput{}, fmt.Errorf("update user: %w", err)
			}
			_, err = p.UserRepo.UpdateUser(ctx,usr)
			if err != nil{
				return dto.ProfileOutput{}, errorswrapper.WrapRepoErr("p.userRepo.UpdateUser", err)
			}
	}
	return dto.ProfileOutput{}, domainerrors.ErrWrongRole
}
