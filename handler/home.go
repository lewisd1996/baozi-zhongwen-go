package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/lewisd1996/baozi-zhongwen/view/home"
)

type HomeHandler struct {
}

func NewHomeHandler() HomeHandler {
	return HomeHandler{}
}

func (h HomeHandler) HandleHomeShow(c echo.Context) error {
	userId := c.Get("user_id").(string)
	return Render(c, home.Show(userId, c.Path()))
}
