package controllers

import (
	"errors"
	"fmt"
	"trabalho/initializers"
	"trabalho/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func generateSequentialID() (string, error) {
	var count int64
	if err := initializers.DB.Unscoped().Model(&models.Post{}).Count(&count).Error; err != nil {
		return "", err
	}

	nextID := count + 1
	return fmt.Sprintf("%d", nextID), nil
}

func PostsCreate(c *gin.Context) {
	//Pegar dados
	var body struct {
		Title  string
		Artist string
		Price  float64
		Gender string
	}

	c.Bind(&body)

	sequentialID, err := generateSequentialID()
	if err != nil {
		c.Status(500)
		return
	}

	//Criar um post
	post := models.Post{
		ID:     sequentialID,
		Title:  body.Title,
		Artist: body.Artist,
		Price:  body.Price,
		Gender: body.Gender,
	}

	if err := validateCreatePostFields(post); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result := initializers.DB.Create(&post) // pass pointer of data to Create
	if result.Error != nil {
		c.Status(400)
		return
	}
	//Retorno
	c.JSON(200, gin.H{
		"post": post,
	})
}

func validateCreatePostFields(post models.Post) error {
	if post.Title == "" || post.Artist == "" || post.Price == 0 || post.Gender == "" {
		return errors.New("Todos os campos são obrigatórios")
	}
	return nil
}

func PostsIndex(c *gin.Context) {
	//Pegar os Posts
	var posts []models.Post
	initializers.DB.Find(&posts)

	//Responder esses Posts
	c.JSON(200, gin.H{
		"post": posts,
	})
}

func PostsShow(c *gin.Context) {
	// Pegar o ID através da URL
	id := c.Param("id")

	// Achar o post que queremos mostrar
	var post models.Post
	result := initializers.DB.First(&post, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{
				"error": "Post não encontrado",
			})
			return
		}

		c.JSON(500, gin.H{
			"error": "Erro no servidor: " + result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"post": post,
	})
}

func PostsUpdate(c *gin.Context) {
	//Pegar o ID atraves da URl
	id := c.Param("id")

	//Pegar os dados do body
	var body struct {
		Title  string
		Artist string
		Price  float64
		Gender string
	}
	c.Bind(&body)

	//Achar o post que queremos atualizar
	var post models.Post
	initializers.DB.First(&post, id)

	//Atualizar
	initializers.DB.Model(&post).Updates(models.Post{
		Title:  body.Title,
		Artist: body.Artist,
		Price:  body.Price,
		Gender: body.Gender,
	})

	//Responder esses Posts
	c.JSON(200, gin.H{
		"posts": post,
	})
}

func PostsDelete(c *gin.Context) {
	id := c.Param("id")
	fmt.Println(id)

	var post models.Post
	initializers.DB.Delete(&post, id)

	//Responder esses Posts
	c.Status(200)
}
