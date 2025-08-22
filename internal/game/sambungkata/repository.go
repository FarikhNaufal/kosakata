package sambungkata

import (
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]Word, error)
	FindById(ID string) (Word, error)
	FindRandomWord() (Word, error)
	FindTodayWord() (Word, error)
	StoreWord(word Word) (Word, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]Word, error) {
	var word []Word
	err := r.db.Raw("SELECT * FROM words").Scan(&word).Error

	return word, err
}

func (r *repository) FindById(id string) (Word, error) {
	var word Word
	err := r.db.Raw("SELECT * FROM words WHERE id=?", id).First(&word).Error

	return word, err
}

func (r *repository) FindRandomWord() (Word, error) {
	var word Word
	err := r.db.Order("RAND()").First(&word).Error

	return word, err
}

func (r *repository) FindTodayWord() (Word, error) {
	var word Word
	today := time.Now().Format("2006-01-02")

	err := r.db.Raw("SELECT * FROM words WHERE release_at=?", today).First(&word).Error

	return word, err
}

func (r *repository) StoreWord(word Word) (Word, error) {

	err := r.db.Create(&word).Error

	return word, err
}
