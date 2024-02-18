package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/lewisd1996/baozi-zhongwen/app"
)

type HealthHandler struct {
	app *app.App
}

/* -------------------------------------------------------------------------- */
/*                                    Init                                    */
/* -------------------------------------------------------------------------- */

func NewHealthHandler(a *app.App) HealthHandler {
	return HealthHandler{app: a}
}

/* --------------------------------- Health --------------------------------- */

func (h HealthHandler) HandleHealth(c echo.Context) error {
	return c.NoContent(200)
}
