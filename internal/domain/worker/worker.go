package worker

import(
	"github.com/google/uuid"

	"github.com/hihikaAAa/TrashProject/internal/domain/task"
)

type Worker struct{
	ID uuid.UUID
	FirstName string
	Surname string
	LastName string
	TaskList []task.Task
	IsActive bool
	// TODO : Добавить район для работы 
}