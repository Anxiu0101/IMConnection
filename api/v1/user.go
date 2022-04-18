package v1

import (
	"IMConnection/api"
	"IMConnection/pkg/logging"
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
	if err := c.ShouldBind(&accountService); err != nil {
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
	if err := c.ShouldBind(&accountService); err != nil {
		res := accountService.Login()
		c.JSON(http.StatusOK, res)
	} else {
		logging.Info(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
	}
}

func ResetUserPassword(c *gin.Context) {
	var accountService service.AccountService
	if err := c.ShouldBind(&accountService); err != nil {
		//TODO 重设密码接口从 access token 中解析用户 id
		res := accountService.ResetPassword(1, accountService.Password)
		c.JSON(http.StatusOK, res)
	} else {
		logging.Info(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
	}
}

func GetUserInfo(c *gin.Context) {

}

func UpdateUserInfo(c *gin.Context) {

}
