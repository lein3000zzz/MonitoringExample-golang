package handler

import (
	"server/internal/pkg/domain"

	"github.com/labstack/echo"
)

type Handler struct {
	CommentSvc domain.CommentService
}

func (h Handler) Create(ctx echo.Context) error {
	var comment domain.Comment

	err := ctx.Bind(&comment)
	if err != nil {
		return err
	}

	tid := ctx.Param("tid")

	return h.CommentSvc.Create(tid, comment)
}

func (h Handler) Like(ctx echo.Context) error {
	tid := ctx.Param("tid")
	cid := ctx.Param("cid")

	return h.CommentSvc.Like(tid, cid)
}
