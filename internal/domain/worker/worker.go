package worker

import(
	"github.com/google/uuid"

	"github.com/hihikaAAa/TrashProject/internal/domain/task"
)

type Worker struct{
	id uuid.UUID
	first_name string
	surname string
	last_name string
	taskList []task.Task
	// TODO : Добавить район для работы 
}