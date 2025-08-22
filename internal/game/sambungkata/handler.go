package sambungkata

import (
	"encoding/json"
	"fmt"
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
		ctx.JSON(ctx.Request.Response.StatusCode, gin.H{
			"error": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": word,
	})
}

func (h *WordHandler) ShowWord(ctx *gin.Context) {
	id := ctx.Param("id")
	word, err := h.wordService.FindById(id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": word,
	})
}

func (h *WordHandler) GetTodayWord(ctx *gin.Context) {
	word, err := h.wordService.FindTodayWord()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": word,
	})
}

func (h *WordHandler) CheckingWord(ctx *gin.Context) {
	var nextWordRequest NextWordRequest
	err := ctx.ShouldBindBodyWithJSON(&nextWordRequest)
	if err != nil {
		errMsgs := []string{}
		for _, e := range err.(validator.ValidationErrors) {
			errMsg := fmt.Sprintf("Error validasi pada field %s, detail: %s", e.Field(), e.Tag())
			errMsgs = append(errMsgs, errMsg)
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"error": errMsgs})
		return
	}

	word, err := h.wordService.FindTodayWord()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var list []string
	if err := json.Unmarshal(word.List, &list); err != nil {
		return
	}
	success := false
	for i := 0; i < len(list)-1; i++ {
		if nextWordRequest.PrevWord == nil {
			if strings.EqualFold(word.Start, nextWordRequest.NextWord) {
				success = true
			}
		} else {
			for i, w := range list {
				if strings.EqualFold(*nextWordRequest.PrevWord, w) && i+1 < len(list) && strings.EqualFold(nextWordRequest.NextWord, list[i+1]) {
					success = true
					break
				}
			}
		}
	}
	if !success {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"data": gin.H{
				"prev_word": nextWordRequest.PrevWord,
				"message":   "wrong guess",
				"success":   success,
			},
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"prev_word": nextWordRequest.NextWord,
			"success":   success,
		},
	})

	// fmt.Println(word)
}

func (h *WordHandler) StoreWord(ctx *gin.Context) {
	var wordRequest WordRequest

	err := ctx.ShouldBindBodyWithJSON(&wordRequest)
	if err != nil {
		errMsgs := []string{}
		for _, e := range err.(validator.ValidationErrors) {
			errMsg := fmt.Sprintf("Error validasi pada field %s, detail: %s", e.Field(), e.Tag())
			errMsgs = append(errMsgs, errMsg)
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"error": errMsgs})
		return
	}

	word, err := h.wordService.StoreWord(wordRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": word,
	})
}
