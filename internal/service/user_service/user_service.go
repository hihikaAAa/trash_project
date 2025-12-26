// Package userservice
package userservice

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hihikaAAa/TrashProject/internal/domain/address"
	domainerrors "github.com/hihikaAAa/TrashProject/internal/domain/domain_errors"
	"github.com/hihikaAAa/TrashProject/internal/domain/person"
	"github.com/hihikaAAa/TrashProject/internal/domain/user"
	dto "github.com/hihikaAAa/TrashProject/internal/dto"
	addressrepo "github.com/hihikaAAa/TrashProject/internal/postgres/address_repo"
	userrepo "github.com/hihikaAAa/TrashProject/internal/postgres/user_repo"
)

type UserService struct {
	db          *sql.DB
	UserRepo    userrepo.UserRepository
	AddressRepo addressrepo.AddressRepository
}

func NewUserService(db *sql.DB) *UserService {
	uRepo := userrepo.NewUserRepository(db)
	aRepo := addressrepo.NewAddressRepository(db)
	return &UserService{
		db:          db,
		UserRepo:    uRepo,
		AddressRepo: aRepo,
	}
}

func (u *UserService) CreateProfile(ctx context.Context, input dto.Input) (dto.Output, error) {
	if err := u.UserRepo.CheckNotExists(ctx, input.Email); err != nil { // Проверил на существование
		return dto.Output{}, domainerrors.ErrUserExists
	}
	adrs, err := address.NewAddress(input.Street, input.HouseNumber, input.Entrance, input.FloorNumber, input.ApartmentNumber)
	if err != nil {
		return dto.Output{}, fmt.Errorf("new address: %w", err)
	}

	prs, err := person.NewPerson(input.FirstName, input.Surname, input.LastName)
	if err != nil {
		return dto.Output{}, fmt.Errorf("new person: %w", err)
	}

	user, err := user.NewUser(prs, adrs.ID)
	if err != nil {
		return dto.Output{}, fmt.Errorf("new user: %w", err)
	}
	// TODO: СЕЙЧАС КОД БЕЗ ТРАНЗАКЦИЙ, ДОБАВИТЬ ИХ
	err = u.AddressRepo.AddAddress(ctx, adrs)
	if err != nil{
		return dto.Output{}, fmt.Errorf("u.addressRepo.AddAddress: %w", err)
	}
	err = u.UserRepo.AddUser(ctx, user)
	if err != nil{
		return dto.Output{}, fmt.Errorf("u.userRepo.AddUser: %w", err)
	}
	return dto.Output{
		ID:user.ID}, nil
}
