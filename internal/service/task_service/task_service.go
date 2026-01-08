// Package taskservice
package taskservice

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hihikaAAa/TrashProject/internal/domain/task"
	"github.com/hihikaAAa/TrashProject/internal/dto"
	taskrepo "github.com/hihikaAAa/TrashProject/internal/postgres/task_repo"
	errorswrapper "github.com/hihikaAAa/TrashProject/internal/service/errors_wrapper"
)

type TaskAction string

const(
	TaskActionComplete TaskAction = "complete"
	TaskActionCancel TaskAction = "cancel"
	TaskActionDrop TaskAction = "drop"
)

type TaskService struct{
	TaskRepo    taskrepo.TaskRepository
}

func NewTaskService(db *sql.DB) *TaskService{
	tRepo := taskrepo.NewTaskRepository(db)
	return &TaskService{
		TaskRepo: tRepo,
	}
}

func (t *TaskService) CreateTask(ctx context.Context, input dto.TaskInput) (dto.TaskOutput, error) {
	tsk, err := task.NewTask(input.ClientID, input.AddressID, input.Time, task.Role(input.Role)) // Создали задачу
	if err != nil {
		return dto.TaskOutput{}, fmt.Errorf("new task: %w", err)
	}

	err = t.TaskRepo.AddTask(ctx, tsk) // Добавили задачу в БД
	if err != nil {
		return dto.TaskOutput{}, errorswrapper.WrapRepoErr("t.taskRepo.AddTask", err)
	}

	return dto.TaskOutput{
		TaskID: tsk.ID,
	}, nil
}

func (t *TaskService) CompleteTask(ctx context.Context, input dto.TaskInput) (dto.TaskOutput, error) {
	return t.ApplyAction(ctx, input, TaskActionComplete)
}

func (t *TaskService) CancelTask(ctx context.Context, input dto.TaskInput) (dto.TaskOutput, error) {
	return t.ApplyAction(ctx, input, TaskActionCancel)
}

func (t *TaskService) DropTask(ctx context.Context, input dto.TaskInput) (dto.TaskOutput, error) {
	return t.ApplyAction(ctx, input, TaskActionDrop)
}

func (t *TaskService) ApplyAction(ctx context.Context, input dto.TaskInput, action TaskAction) (dto.TaskOutput, error){
	tsk, err := t.TaskRepo.GetByID(ctx,input.ID)     // Получил задачу
	if err != nil{
		return dto.TaskOutput{}, errorswrapper.WrapRepoErr("t.taskRepo.GetByID", err)
	}
	
	switch action{
	case TaskActionComplete:
		err = tsk.CompleteTask(input.Time, task.Role(input.Role))
	case TaskActionCancel:
		err = tsk.CancelTask(input.Time, task.Role(input.Role))
	case TaskActionDrop:
		err = tsk.DropTask(task.Role(input.Role))
	default:
		return dto.TaskOutput{}, fmt.Errorf("unknown action: %s", action)
	}
	if err != nil {
		return dto.TaskOutput{}, fmt.Errorf("%s task: %w", action, err)
	}

	_, err = t.TaskRepo.UpdateStatus(ctx,tsk)  		// Обновил в БД
	if err != nil{
		return dto.TaskOutput{}, errorswrapper.WrapRepoErr("t.taskRepo.UpdateStatus", err)
	}

	return dto.TaskOutput{
		TaskID: tsk.ID,
	}, nil
	
}
