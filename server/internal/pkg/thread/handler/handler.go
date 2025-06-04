package handler

import (
	"server/internal/pkg/domain"

	"github.com/labstack/echo"
)

type Handler struct {
	ThreadSvc domain.ThreadService
}

func (h Handler) GetThread(ctx echo.Context) error {
	tid := ctx.Param("tid")

	t, err := h.ThreadSvc.Get(tid)
	if err != nil {
		return err
	}

	return ctx.JSON(200, t)
}

func (h Handler) CreateThread(ctx echo.Context) error {
	var thread domain.Thread

	err := ctx.Bind(&thread)
	if err != nil {
		return err
	}

	err = h.ThreadSvc.Create(thread)
	if err != nil {
		return err
	}

	return ctx.NoContent(200)
}
