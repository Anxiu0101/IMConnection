package v1

import (
	"IMConnection/api"
	"IMConnection/pkg/logging"
	"IMConnection/pkg/util"
	"IMConnection/service"
	"github.com/gin-gonic/gin"
	"net/http"
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
	var accountService service.AccountService
	if err := c.ShouldBind(&accountService); err == nil {
		res := accountService.Register()
		c.JSON(http.StatusOK, res)
	} else {
		logging.Info(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
	}
}

// UserLogin godoc
// @Summary User Login
// @Schemes
// @Description help User login site
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {string} UserLogin
// @Router /user/login [post]
func UserLogin(c *gin.Context) {
	var accountService service.AccountService
	if err := c.ShouldBind(&accountService); err == nil {
		res := accountService.Login()
		c.JSON(http.StatusOK, res)
	} else {
		logging.Info(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
	}
}

func ResetUserPassword(c *gin.Context) {
	var accountService service.AccountService
	if err := c.ShouldBind(&accountService); err == nil {
		claim, _ := util.ParseToken(c.GetHeader("Authorization"))
		res := accountService.ResetPassword(claim.ID, accountService.Password)
		c.JSON(http.StatusOK, res)
	} else {
		logging.Info(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
	}
}

func RefreshAccessToken(c *gin.Context) {

}

func GetUserInfo(c *gin.Context) {
	var userService service.UserService
	if err := c.ShouldBind(&userService); err == nil {
		claim, _ := util.ParseToken(c.GetHeader("Authorization"))
		res := userService.GetInfo(claim.ID)
		c.JSON(http.StatusOK, res)
	} else {
		logging.Info(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
	}
}

func UpdateUserInfo(c *gin.Context) {
	var userService service.UserService
	if err := c.ShouldBind(&userService); err == nil {
		claim, _ := util.ParseToken(c.GetHeader("Authorization"))
		res := userService.GetInfo(claim.ID)
		c.JSON(http.StatusOK, res)
	} else {
		logging.Info(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
	}
}
