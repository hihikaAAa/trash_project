// Package userservice
package services

/*
import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hihikaAAa/trash_project/internal/domain/address"
	"github.com/hihikaAAa/trash_project/internal/domain/person"
	"github.com/hihikaAAa/trash_project/internal/domain/task"
	"github.com/hihikaAAa/trash_project/internal/domain/user"
	dto "github.com/hihikaAAa/trash_project/internal/dto"
	addressrepo "github.com/hihikaAAa/trash_project/internal/postgres/address_repo"
	taskrepo "github.com/hihikaAAa/trash_project/internal/postgres/task_repo"
	userrepo "github.com/hihikaAAa/trash_project/internal/postgres/user_repo"
	errorswrapper "github.com/hihikaAAa/trash_project/internal/service/errors_wrapper"
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
		return dto.UserOutput{}, errorswrapper.WrapRepoErr("u.userRepo.CheckNotExists", err)
	}
	adrs, err := address.NewAddress(input.Street, input.HouseNumber, input.Entrance, input.FloorNumber, input.ApartmentNumber)
	if err != nil {
		return dto.UserOutput{}, fmt.Errorf("new address: %w", err)
	}

	prs, err := person.NewPerson(input.FirstName, input.Surname, input.LastName)
	if err != nil {
		return dto.UserOutput{}, fmt.Errorf("new person: %w", err)
	}

	usr, err := user.NewUser(prs, adrs.ID)
	if err != nil {
		return dto.UserOutput{}, fmt.Errorf("new user: %w", err)
	}
	// TODO: СЕЙЧАС КОД БЕЗ ТРАНЗАКЦИЙ, ДОБАВИТЬ ИХ
	err = u.AddressRepo.AddAddress(ctx, adrs)
	if err != nil {
		return dto.UserOutput{}, errorswrapper.WrapRepoErr("u.addressRepo.AddAddress",err)
	}
	err = u.UserRepo.AddUser(ctx, usr)
	if err != nil {
		return dto.UserOutput{}, errorswrapper.WrapRepoErr("u.userRepo.AddUser", err)
	}
	return dto.UserOutput{
		UserID: usr.ID,
	}, nil
}

func (u *UserService) UpdateUserInfo(ctx context.Context, input dto.UserInput) (dto.UserOutput, error) {
	usr, err := u.UserRepo.GetByID(ctx, input.ID) // Получили юзера по ID
	if err != nil {
		return dto.UserOutput{}, errorswrapper.WrapRepoErr("u.userRepo.GetByID", err)
	}

	prs, err := person.NewPerson(input.FirstName, input.Surname, input.LastName) // Создали новую персоналию
	if err != nil {
		return dto.UserOutput{}, fmt.Errorf("new person: %w", err)
	}

	err = usr.UpdateUser(*prs)
	if err != nil {
		return dto.UserOutput{}, fmt.Errorf("update user: %w", err)
	} // Обновили информацию о персоналии

	_, err = u.UserRepo.UpdateUser(ctx, usr) // Обновили инфу в БД
	if err != nil {
		return dto.UserOutput{}, errorswrapper.WrapRepoErr("u.userRepo.UpdateUser",err)
	}

	addrs, err := u.AddressRepo.GetByID(ctx, usr.AddressID) // Получили адрес по айдишнику
	if err != nil {
		return dto.UserOutput{}, errorswrapper.WrapRepoErr("u.addressRepo.GetByID",err)
	}

	err = addrs.UpdateAddress(input.Street, input.HouseNumber, input.Entrance, input.FloorNumber, input.ApartmentNumber) // Обновили адрес
	if err != nil {
		return dto.UserOutput{}, fmt.Errorf("update address: %w", err)
	}

	_, err = u.AddressRepo.UpdateAddress(ctx, addrs) // Мб нет смысла возвращать адрес, тк он просто меняется в бд, id остается // Обновили инфу в БД
	if err != nil {
		return dto.UserOutput{}, errorswrapper.WrapRepoErr("u.addressRepo.UpdateAddress", err)
	}

	return dto.UserOutput{
		UserID: usr.ID,
	}, nil
}
*/
