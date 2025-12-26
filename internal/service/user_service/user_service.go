// Package userservice
package userservice

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hihikaAAa/TrashProject/internal/domain/address"
	domainerrors "github.com/hihikaAAa/TrashProject/internal/domain/domain_errors"
	"github.com/hihikaAAa/TrashProject/internal/domain/person"
	"github.com/hihikaAAa/TrashProject/internal/domain/task"
	"github.com/hihikaAAa/TrashProject/internal/domain/user"
	dto "github.com/hihikaAAa/TrashProject/internal/dto"
	addressrepo "github.com/hihikaAAa/TrashProject/internal/postgres/address_repo"
	taskrepo "github.com/hihikaAAa/TrashProject/internal/postgres/task_repo"
	userrepo "github.com/hihikaAAa/TrashProject/internal/postgres/user_repo"
)

type UserService struct {
	UserRepo    userrepo.UserRepository
	AddressRepo addressrepo.AddressRepository
	TaskRepo    taskrepo.TaskRepository
}

func NewUserService(db *sql.DB) *UserService {
	uRepo := userrepo.NewUserRepository(db)
	aRepo := addressrepo.NewAddressRepository(db)
	tRepo := taskrepo.NewTaskRepository(db)
	return &UserService{
		UserRepo:    uRepo,
		AddressRepo: aRepo,
		TaskRepo:    tRepo,
	}
}

func (u *UserService) CreateProfile(ctx context.Context, input dto.UserInput) (dto.UserOutput, error) { // Сейчас при первом создании юзера, передается пустой ID. Спросить, как это хендлить
	if err := u.UserRepo.CheckNotExists(ctx, input.Email); err != nil { // Проверил на существование
		return dto.UserOutput{}, domainerrors.ErrUserExists
	}
	adrs, err := address.NewAddress(input.Street, input.HouseNumber, input.Entrance, input.FloorNumber, input.ApartmentNumber)
	if err != nil {
		return dto.UserOutput{}, fmt.Errorf("new address: %w", err)
	}

	prs, err := person.NewPerson(input.FirstName, input.Surname, input.LastName)
	if err != nil {
		return dto.UserOutput{}, fmt.Errorf("new person: %w", err)
	}

	user, err := user.NewUser(prs, adrs.ID)
	if err != nil {
		return dto.UserOutput{}, fmt.Errorf("new user: %w", err)
	}
	// TODO: СЕЙЧАС КОД БЕЗ ТРАНЗАКЦИЙ, ДОБАВИТЬ ИХ
	err = u.AddressRepo.AddAddress(ctx, adrs)
	if err != nil {
		return dto.UserOutput{}, fmt.Errorf("u.addressRepo.AddAddress: %w", err)
	}
	err = u.UserRepo.AddUser(ctx, user)
	if err != nil {
		return dto.UserOutput{}, fmt.Errorf("u.userRepo.AddUser: %w", err)
	}
	return dto.UserOutput{
		UserID: user.ID,
	}, nil
}

func (u *UserService) UpdateUserInfo(ctx context.Context, input dto.UserInput) (dto.UserOutput, error) {
	user, err := u.UserRepo.GetByID(ctx, input.ID) // Получили юзера по ID
	if err != nil {
		return dto.UserOutput{}, fmt.Errorf("u.userRepo.GetByID: %w", err)
	}

	prs, err := person.NewPerson(input.FirstName, input.Surname, input.LastName) // Создали новую персоналию
	if err != nil {
		return dto.UserOutput{}, fmt.Errorf("new person: %w", err)
	}

	err = user.UpdateUser(*prs)
	if err != nil {
		return dto.UserOutput{}, fmt.Errorf("update user: %w", err)
	} // Обновили информацию о персоналии

	_, err = u.UserRepo.UpdateUser(ctx, user) // Обновили инфу в БД
	if err != nil {
		return dto.UserOutput{}, fmt.Errorf("u.userRepo.UpdateUser: %w", err)
	}

	addrs, err := u.AddressRepo.GetByID(ctx, user.AddressID) // Получили адрес по айдишнику
	if err != nil {
		return dto.UserOutput{}, fmt.Errorf("u.addressRepo.GetByID: %w", err)
	}

	err = addrs.UpdateAddress(input.Street, input.HouseNumber, input.Entrance, input.FloorNumber, input.ApartmentNumber) // Обновили адрес
	if err != nil {
		return dto.UserOutput{}, fmt.Errorf("update address: %w", err)
	}

	_, err = u.AddressRepo.UpdateAddress(ctx, addrs) // Мб нет смысла возвращать адрес, тк он просто меняется в бд, id остается // Обновили инфу в БД
	if err != nil {
		return dto.UserOutput{}, fmt.Errorf("u.addressRepo.UpdateAddress: %w", err)
	}

	return dto.UserOutput{
		UserID: user.ID,
	}, nil
}

func (u *UserService) CreateTask(ctx context.Context, input dto.TaskInput) (dto.TaskOutput, error) {
	tsk, err := task.NewTask(input.ClientID, input.AddressID, input.Time) // Создали задачу
	if err != nil {
		return dto.TaskOutput{}, fmt.Errorf("new task: %w", err)
	}

	err = u.TaskRepo.AddTask(ctx, tsk) // Добавили задачу в БД
	if err != nil {
		return dto.TaskOutput{}, fmt.Errorf("u.taskRepo.AddTask : %w", err)
	}

	return dto.TaskOutput{
		TaskID: tsk.ID,
	}, nil
}
