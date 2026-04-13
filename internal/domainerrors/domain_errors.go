// Package domainerrors
package domainerrors

import (
	"errors"
	"net/http"

	httpres "github.com/hihikaAAa/trash_project/pkg/http_res"
)

var (
	ErrForbidden = httpres.NewHTTPError(
		errors.New("Доступ запрещен"),
		http.StatusForbidden,
		httpres.CodeDenied,
		"denied",
	)

	ErrTaskNotFound = httpres.NewHTTPError(
		errors.New("Задача не найдена"),
		http.StatusNotFound,
		httpres.CodeNotFound,
		"not_found",
	)

	ErrWrongRole = httpres.NewHTTPError(
		errors.New("Неправильная роль"),
		http.StatusForbidden,
		httpres.CodeDenied,
		"denied",
	)

	ErrBadTaskInfo = httpres.NewHTTPError(
		errors.New("Невалидное тело запроса задачи"),
		http.StatusBadRequest,
		httpres.CodeBadRequest,
		"bad_request",
	)
	ErrBadTaskTime = httpres.NewHTTPError(
		errors.New("Невалидное время задачи"),
		http.StatusBadRequest,
		httpres.CodeBadRequest,
		"bad_request",
	)
	ErrBadTaskWorker = httpres.NewHTTPError(
		errors.New("Невалидный работник"),
		http.StatusBadRequest,
		httpres.CodeBadRequest,
		"bad_request",
	)
	ErrBadTaskAddress = httpres.NewHTTPError(
		errors.New("Невалидный адрес задачи"),
		http.StatusBadRequest,
		httpres.CodeBadRequest,
		"bad_request",
	)
	ErrBadTaskDescription = httpres.NewHTTPError(
		errors.New("невалидное описание задачи"),
		http.StatusBadRequest,
		httpres.CodeBadRequest,
		"bad_request",
	)

	ErrInvalidStatusTransition = httpres.NewHTTPError(
		errors.New("Невалидный статус"),
		http.StatusConflict,
		httpres.CodeConflict,
		"conflict",
	)
	ErrTaskCanceled = httpres.NewHTTPError(
		errors.New("задача уже отменена"),
		http.StatusConflict,
		httpres.CodeConflict,
		"conflict",
	)
	ErrTaskDone = httpres.NewHTTPError(
		errors.New("Задача уже выполнена"),
		http.StatusConflict,
		httpres.CodeConflict,
		"conflict",
	)
	ErrTaskNotInProgress = httpres.NewHTTPError(
		errors.New("Завершить можно только ту задачу, которая в процессе исполнения"),
		http.StatusConflict,
		httpres.CodeConflict,
		"conflict",
	)
)
