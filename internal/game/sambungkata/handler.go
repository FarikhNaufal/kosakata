package sambungkata

import (
	"encoding/json"
	"fmt"
	"kosakata/internal/utils/response"
	"math"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type WordHandler struct {
	wordService Service
}

func NewHandler(wordService Service) *WordHandler {
	return &WordHandler{wordService}
}

func (h *WordHandler) ShowAllWord(ctx *gin.Context) {
	word, err := h.wordService.FindAll()
	if err != nil {
		response.Failed(ctx, http.StatusInternalServerError, nil)
		return
	}
	response.Success(ctx, "Success", word)
}

func (h *WordHandler) ShowWord(ctx *gin.Context) {
	id := ctx.Param("id")
	word, err := h.wordService.FindById(id)

	if err != nil {
		response.Failed(ctx, http.StatusInternalServerError, nil)
		return
	}
	response.Success(ctx, "Success", word)
}

func (h *WordHandler) GetTodayWord(ctx *gin.Context) {
	word, err := h.wordService.FindTodayWord()
	if err != nil {
		response.Failed(ctx, http.StatusInternalServerError, nil)
		return
	}
	res := TodayWordDTO{
		ID:    word.ID,
		Start: word.Start,
	}
	response.Success(ctx, "Success", res)
}

func (h *WordHandler) CheckingNextWord(ctx *gin.Context) {
	var nextWordRequest NextWordRequest
	err := ctx.ShouldBindBodyWithJSON(&nextWordRequest)
	if err != nil {
		errMsgs := []string{}
		for _, e := range err.(validator.ValidationErrors) {
			errMsg := fmt.Sprintf("Error validasi pada field %s, detail: %s", e.Field(), e.Tag())
			errMsgs = append(errMsgs, errMsg)
		}

		response.Failed(ctx, http.StatusBadRequest, nil, &errMsgs)
		return
	}

	word, err := h.wordService.FindTodayWord()
	if err != nil {
		response.Failed(ctx, http.StatusInternalServerError, nil, err)
		return
	}

	var list []string
	if err := json.Unmarshal(word.List, &list); err != nil {
		return
	}

	success := false
	prev_word := nextWordRequest.PrevWord
	var clue string
	progress := 0.0

	if nextWordRequest.PrevWord == nil {
		if strings.EqualFold(nextWordRequest.NextWord, list[0]) {
			success = true
			progress = 1 / float64(len(list)+1) * 100
		}
	} else {
		for i := range list {
			if i < len(list)-1 &&
				strings.EqualFold(*nextWordRequest.PrevWord, list[i]) &&
				strings.EqualFold(nextWordRequest.NextWord, list[i+1]) {
				success = true
				progress = float64(i+2) / float64(len(list)+1) * 100

				clue = list[i+1]
				break
			}

			if len(list)-1 == i {
				if strings.EqualFold(nextWordRequest.NextWord, word.End) {
					success = true
					progress = float64(i+2) / float64(len(list)+1) * 100
					break
				}
				clue = word.End
			}

		}
	}

	if !success {
		res := WrongWordDTO{
			Clue:     []string{string(clue[0]), string(clue[len(clue)-1])},
			Length:   len(clue),
			PrevWord: prev_word,
		}

		msg := "Wrong next word."
		response.Failed(ctx, http.StatusUnprocessableEntity, &msg, res)
		return
	}
	response.Success(ctx, "Success",
		gin.H{
			"prev_word": nextWordRequest.NextWord,
			"progress":  math.Round(progress*100) / 100,
		},
	)

}

// function WordMatcher()

func (h *WordHandler) StoreWord(ctx *gin.Context) {
	var wordRequest WordRequest

	err := ctx.ShouldBindBodyWithJSON(&wordRequest)
	if err != nil {
		errMsgs := []string{}
		for _, e := range err.(validator.ValidationErrors) {
			errMsg := fmt.Sprintf("Error validasi pada field %s, detail: %s", e.Field(), e.Tag())
			errMsgs = append(errMsgs, errMsg)
		}

		response.Failed(ctx, http.StatusBadRequest, nil, errMsgs)

		return
	}

	word, err := h.wordService.StoreWord(wordRequest)
	if err != nil {
		response.Failed(ctx, http.StatusInternalServerError, nil, err)
		return
	}

	response.Success(ctx, "Success", word)
}
