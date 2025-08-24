package sambungkata

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Service interface {
	FindAll() ([]Word, error)
	FindById(ID string) (Word, error)
	FindRandomWord() (Word, error)
	FindTodayWord() (Word, error)
	StoreWord(wordRequest WordRequest)(Word, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) FindAll() ([]Word, error) {
	return s.repository.FindAll()
}

func (s *service) FindById(id string) (Word, error) {
	return s.repository.FindById(id)
}

func (s *service) FindRandomWord() (Word, error) {
	return s.repository.FindRandomWord()
}

func (s *service) FindTodayWord() (Word, error) {
	word, err := s.repository.FindTodayWord()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return s.repository.FindRandomWord()
	}

	return word, err
}


func (s *service) StoreWord(wordRequest WordRequest) (Word, error) {
	word := Word{
		Start:     wordRequest.Start,
		End:       wordRequest.End,
		List:      wordRequest.List,
		CreatedAt: time.Now(),
		ReleaseAt: wordRequest.ReleaseAt,
	}

	newWord, err := s.repository.StoreWord(word)
	return newWord, err
}
