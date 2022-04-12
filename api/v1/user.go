package v1

import (
	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// UserRegister godoc
// @Summary User Register
// @Schemes
// @Description help User Create Account
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {string} UserRegister
// @Router /user/register [post]
func UserRegister(c *gin.Context) {
	c.JSON(200, "check")
}

func UserLogin(c *gin.Context) {

}

func GetUserInfo(c *gin.Context) {

}

func UpdateUserInfo(c *gin.Context) {

}

func ResetUserPassword(c *gin.Context) {

}
