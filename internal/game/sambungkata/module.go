package sambungkata

import (

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(group *gin.RouterGroup, db *gorm.DB) {
	r := NewRepository(db)
	s := NewService(r)
	h := NewHandler(s)

	RegisterRoute(group, h)
}
