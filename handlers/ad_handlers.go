package handlers

import (
	"fmt"
	"gomatri/models"
	"gomatri/storage"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdHandler struct {
	store *storage.AdSqlliteStore
}

func CreateNewAdHandler() *AdHandler {
	return &AdHandler{
		store: storage.NewSqliteAdsStore(),
	}
}

func (ah *AdHandler) GetAds(c *gin.Context) {
	log.Println("IN GET handler")

	ads, err := ah.store.Get()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("ADS List", ads)
}

func (ah *AdHandler) UpdateAd(c *gin.Context) {
	log.Println("IN PATCH AD handler")
	var input models.Ad
	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("IN PATCH AD handler ", &input)
	ID, err := ah.store.Update(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("USER Updated ID: ", ID)

}

func (ah *AdHandler) DeleteAd(c *gin.Context) {
	log.Println("IN Delete handler")

	id := c.Param("id")

	ID, err := ah.store.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("USER deleted ID: ", ID)
}
func (ah *AdHandler) GetAd(c *gin.Context) {
	log.Println("IN GET one AD handler")
	id := c.Param("id")

	user, err := ah.store.GetOne(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("...... Get Ad: ", user)
}

func (ah *AdHandler) CreateAd(c *gin.Context) {
	log.Println("IN Create AD handler")
	c.Request.ParseForm()
	for key, value := range c.Request.PostForm {
		fmt.Println(key,value)
	}
	var input models.Ad
	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("In POST AD hanlder INPUT IS :", input)

	ID, err := ah.store.Create(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("AD CREATED ID: ", ID)
}
