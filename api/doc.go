package api

import (
	"CollabDoc-go/store"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Doc struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func CreateDoc(c *gin.Context) {
	var doc Doc
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	id := store.SaveDoc(doc.Name)
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func GetDoc(c *gin.Context) {
	id := c.Param("id")
	doc := store.LoadDoc(id)
	c.JSON(http.StatusOK, doc)
}
