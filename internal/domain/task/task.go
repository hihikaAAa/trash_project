package task

import (
	"time"

	"github.com/google/uuid"

)
type Status string 

const(
	StatusOpen Status = "OPEN"
	StatusInProgress Status = "IN PROGRESS"
	StatusDone Status = "DONE"
	StatusCanceled Status = "CANCELED"
)
type Task struct{
	Id uuid.UUID
	ClientID uuid.UUID
	AddressID uuid.UUID
	WorkerID *uuid.UUID
	CreatedAt time.Time
	ClosedAt *time.Time
}
