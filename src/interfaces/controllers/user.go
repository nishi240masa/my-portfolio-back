package controllers

import (
	"net/http"
	"w3st/dto"
	"w3st/usecase"

	"github.com/gin-gonic/gin"
)


type UserController struct{
	userUsecase usecase.UserUsecase
}

func NewUserController(userUsecase usecase.UserUsecase) *UserController {
	return &UserController{
		userUsecase: userUsecase,
	}
}

func (controller *UserController) Signup(c *gin.Context) {
	var input dto.SignupData

	// リクエストのバインド
	if  err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ユーザー登録
	user, err := controller.userUsecase.CreateUser(input.Name, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (controller *UserController) Login(c *gin.Context) {
	var input dto.LoginData

	// リクエストのバインド
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ログイン
	token, err := controller.userUsecase.FindUser(input.Email)
		if err != nil {
			if err.Error() == "User not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})

	}
