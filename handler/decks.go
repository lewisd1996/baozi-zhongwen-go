package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/lewisd1996/baozi-zhongwen/view/decks"
)

type DecksHandler struct{}

func NewDecksHandler() DecksHandler {
	return DecksHandler{}
}

func (h DecksHandler) HandleDecksShow(c echo.Context) error {
	return Render(c, decks.Show(c.Path()))
}
