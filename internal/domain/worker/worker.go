package worker

import(
	"github.com/google/uuid"

	"github.com/hihikaAAa/TrashProject/internal/domain/task"
)

type Worker struct{
	ID uuid.UUID
	First_name string
	Surname string
	Last_name string
	TaskList []task.Task
	Is_active bool
	// TODO : Добавить район для работы 
}