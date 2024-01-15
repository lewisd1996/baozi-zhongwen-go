package handler

import (
	"math/rand"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lewisd1996/baozi-zhongwen/model"
	"github.com/lewisd1996/baozi-zhongwen/view/home"
)

type HomeHandler struct {
	entries []model.Vocab
}

func NewHomeHandler(entries []model.Vocab) HomeHandler {
	return HomeHandler{entries}
}

func (h HomeHandler) HandleHomeShow(c echo.Context) error {

	// Generate random number between 1 and len(entries)
	entryNumber := rand.Intn(len(h.entries)) + 1

	// Get the entry from the slice
	entry := h.entries[entryNumber]

	english := strings.Join(entry.Translations["eng"], ", ")

	return Render(c, home.Show(entry.Hanzi, entry.Pinyin, english, entry.Level))
}
