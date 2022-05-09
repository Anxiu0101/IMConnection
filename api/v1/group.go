package v1

import (
	"IMConnection/model"
	"IMConnection/pkg/logging"
	"IMConnection/pkg/util"
	"IMConnection/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateGroup(c *gin.Context) {
	var groupService service.GroupService
	if err := c.ShouldBind(groupService); err == nil {
		res := groupService.Create()
		c.JSON(http.StatusOK, res)
	} else {
		logging.Info(err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse(err))
	}
}

func InviteMember(c *gin.Context) {
	var groupService service.GroupService
	if err := c.ShouldBind(&groupService); err == nil {
		claim, _ := util.ParseToken(c.GetHeader("Authorization"))
		res := groupService.Invite(claim.ID)
		c.JSON(http.StatusOK, res)
	} else {
		logging.Info(err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse(err))
	}
}
