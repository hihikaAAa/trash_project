package task

import (
	"time"

	"github.com/google/uuid"

)
type Status string 

const(
	StatusOpen Status = "OPEN"
	StatusInProgress Status = "IN_PROGRESS"
	StatusDone Status = "DONE"
	StatusCanceled Status = "CANCELED"
)

type Task struct{
	ID uuid.UUID
	ClientID uuid.UUID
	AddressID uuid.UUID
	Status Status 
	WorkerID *uuid.UUID
	CreatedAt time.Time
	ClosedAt *time.Time
}
