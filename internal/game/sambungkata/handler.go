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

// GetTodayWord godoc
// @Summary      Get today's word
// @Description  Mengambil kata yang digunakan untuk hari ini
// @Tags         Words
// @Produce      json
// @Success      200  {object}  TodayWordDTO
// @Router       /word/today/{id} [get]
func (h *WordHandler) GetTodayWord(ctx *gin.Context) {
	id := ctx.Param("id")
	word, err := h.wordService.FindTodayWord(id)
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

type CheckWord struct {
	Success  bool
	Clue     string
	Progress float64
	Position []WordPosition
}

type WordPosition struct {
	Position int    `json:"position"`
	Color    string `json:"color"`
	Char     string `json:"char"`
}

func CheckWordPosition(word, answer string) (wp []WordPosition) {
	maxLen := len(word)
	if len(answer) > maxLen {
		maxLen = len(answer)
	}

	result := make([]int, maxLen)
	freq := make(map[rune]int)

	for _, ch := range answer {
		freq[ch]++
	}

	// step 1: posisi sama
	for i := 0; i < maxLen; i++ {
		if i < len(word) && i < len(answer) && word[i] == answer[i] {
			result[i] = 1
			freq[rune(word[i])]--
		}
	}

	// step 2: posisi beda
	for i := 0; i < len(word); i++ {
		ch := rune(word[i])
		if result[i] == 0 && freq[ch] > 0 {
			result[i] = 2
			freq[ch]--
		}
	}

	colorMap := map[int]string{
		0: "#FF8B94", // salah (error)
		1: "#A8E6CF", // benar (success)
		2: "#FFD3B6", // ada tapi di posisi beda (warning)
	}

	for i, val := range result {
		char := "?"
		if i < len(word) {
			char = string(word[i])
		}
		wp = append(wp, WordPosition{
			Position: val,
			Color:    colorMap[val],
			Char:     char,
		})
	}

	return wp
}

func MatchingWord(list []string, nextWordRequest NextWordRequest, word Word) (res CheckWord, error error) {
	res.Success = false
	res.Progress = 0.0

	if nextWordRequest.PrevWord == nil {
		res.Position = CheckWordPosition(nextWordRequest.NextWord, list[0])
		res.Clue = list[0]
		if strings.EqualFold(nextWordRequest.NextWord, list[0]) {
			res.Success = true
			res.Progress = 1 / float64(len(list)+1) * 100
		}
	} else {
		for i := range list {
			if len(list)-1 == i {
				res.Clue = word.End
				res.Progress = float64(i+1) / float64(len(list)+1) * 100
				res.Position = CheckWordPosition(nextWordRequest.NextWord, word.End)
				if strings.EqualFold(nextWordRequest.NextWord, word.End) {
					// kondisi cocok
					res.Success = true
					res.Progress = float64(i+2) / float64(len(list)+1) * 100
					break
				}
			}

			if i < len(list)-1 {
				res.Position = CheckWordPosition(nextWordRequest.NextWord, list[i+1])
				res.Progress = float64(i+1) / float64(len(list)+1) * 100
				if strings.EqualFold(*nextWordRequest.PrevWord, list[i]) &&
					strings.EqualFold(nextWordRequest.NextWord, list[i+1]) {
					// kondisi cocok
					res.Success = true
					res.Progress = float64(i+2) / float64(len(list)+1) * 100
					break
				} else if strings.EqualFold(*nextWordRequest.PrevWord, list[i]) {
					res.Clue = list[i+1]
					break
				}
			}
		}
	}

	res.Progress = math.Round(res.Progress*100) / 100
	return
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

	word, err := h.wordService.FindTodayWord(nextWordRequest.ID)
	if err != nil {
		response.Failed(ctx, http.StatusInternalServerError, nil, err)
		return
	}

	var list []string
	if err := json.Unmarshal(word.List, &list); err != nil {
		response.Failed(ctx, http.StatusInternalServerError, nil, err)
		return
	}

	matcher, err := MatchingWord(list, nextWordRequest, word)
	if err != nil {
		response.Failed(ctx, http.StatusInternalServerError, nil, err)
		return
	}

	if !matcher.Success {
		res := WrongWordDTO{
			// Clue:     []string{string(matcher.Clue[0])},
			Clue:     []string{string(matcher.Clue[0])},
			Length:   len(matcher.Clue),
			PrevWord: nextWordRequest.PrevWord,
			Details:  matcher.Position,
			Progress: matcher.Progress,
		}

		msg := "Wrong next word."
		response.Failed(ctx, http.StatusUnprocessableEntity, &msg, res)
		return
	}
	response.Success(ctx, "Success",
		gin.H{
			"prev_word": nextWordRequest.NextWord,
			"progress":  matcher.Progress,
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
