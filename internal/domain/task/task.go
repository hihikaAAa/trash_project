package task

import (
	"time"

	"github.com/google/uuid"

	"github.com/hihikaAAa/TrashProject/internal/domain/address"
)

type Task struct{
	id uuid.UUID
	clientName string
	clientLastName string
	taskInfo address.Address
	workerID uuid.UUID
	createdAt time.Time
	closedAt time.Time
}