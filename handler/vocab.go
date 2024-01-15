package handler

import (
	"math/rand"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lewisd1996/baozi-zhongwen/model"
	"github.com/lewisd1996/baozi-zhongwen/view/components"
)

type VocabHandler struct {
	entries []model.Vocab
}

func NewVocabHandler(entries []model.Vocab) VocabHandler {
	return VocabHandler{entries}
}

func (h VocabHandler) HandleVocabShow(c echo.Context) error {
	// Generate random number between 1 and len(entries)
	entryNumber := rand.Intn(len(h.entries)) + 1

	// Get the entry from the slice
	entry := h.entries[entryNumber]

	english := strings.Join(entry.Translations["eng"], ", ")
	return HTML(c, components.VocabCard(entry.Hanzi, entry.Pinyin, english, entry.Level))
}
