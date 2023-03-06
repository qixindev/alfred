package admin

import (
	"accounts/data"
	"accounts/middlewares"
	"accounts/models"
	"accounts/models/dto"
	"accounts/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func addAdminUsersRoutes(rg *gin.RouterGroup) {
	rg.GET("/users", func(c *gin.Context) {
		var users []models.User
		if middlewares.TenantDB(c).Find(&users).Error != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, utils.Filter(users, models.User2AdminDto))
	})

	rg.GET("/users/:userId", func(c *gin.Context) {
		userId := c.Param("userId")
		var user models.User
		if middlewares.TenantDB(c).First(&user, "id = ?", userId).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, user.AdminDto())
	})

	rg.POST("/users", func(c *gin.Context) {
		tenant := middlewares.GetTenant(c)
		var user models.User
		err := c.BindJSON(&user)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		user.TenantId = tenant.Id
		if data.DB.Create(&user).Error != nil {
			c.Status(http.StatusConflict)
			return
		}
		c.JSON(http.StatusOK, user.AdminDto())
	})

	rg.PUT("/users/:userId", func(c *gin.Context) {
		userId := c.Param("userId")
		var user models.User
		if middlewares.TenantDB(c).First(&user, "id = ?", userId).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		var u models.User
		err := c.BindJSON(&u)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		user.Username = u.Username
		user.FirstName = u.FirstName
		user.LastName = u.LastName
		user.DisplayName = u.DisplayName
		user.Email = u.Email
		user.EmailVerified = u.EmailVerified
		user.Phone = u.Phone
		user.PhoneVerified = u.PhoneVerified
		user.TwoFactorEnabled = u.TwoFactorEnabled
		user.Disabled = u.Disabled
		if data.DB.Save(&user).Error != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, user.AdminDto())
	})

	rg.DELETE("/users/:userId", func(c *gin.Context) {
		userId := c.Param("userId")
		var user models.User
		if middlewares.TenantDB(c).First(&user, "id = ?", userId).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		if data.DB.Delete(&user).Error != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusNoContent)
	})

	rg.GET("/users/:userId/groups", func(c *gin.Context) {
		userId := c.Param("userId")
		var user models.User
		if middlewares.TenantDB(c).First(&user, "id = ?", userId).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		var groupUsers []models.GroupUser
		if data.DB.Joins("Group", "group_users.group_id = groups.id AND group_users.tenant_id = groups.tenant_id").
			Find(&groupUsers, "group_users.tenant_id = ? AND user_id = ?", user.TenantId, user.Id).Error != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		groups := utils.Filter(groupUsers, func(gu models.GroupUser) dto.GroupMemberDto {
			return dto.GroupMemberDto{
				Id:   gu.GroupId,
				Name: gu.Group.Name,
				Role: gu.Role,
			}
		})
		c.JSON(http.StatusOK, groups)
	})

	rg.PUT("/users/:userId/groups/:groupId", func(c *gin.Context) {
		userId := c.Param("userId")
		var user models.User
		if middlewares.TenantDB(c).First(&user, "id = ?", userId).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		groupId := c.Param("groupId")
		var group models.Group
		if middlewares.TenantDB(c).First(&group, "id = ?", groupId).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		var gu dto.GroupMemberDto
		if c.BindJSON(&gu) != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		var groupUser models.GroupUser
		if middlewares.TenantDB(c).First(groupUser, "group_id = ? AND user_id = ?", group.Id, user.Id).Error != nil {
			// Not found, create one.
			groupUser.UserId = user.Id
			groupUser.GroupId = group.Id
			groupUser.TenantId = user.TenantId
			groupUser.Role = gu.Role
		} else {
			// Found, update it.
			groupUser.Role = gu.Role
			if middlewares.TenantDB(c).Save(&groupUser).Error != nil {
				c.Status(http.StatusInternalServerError)
				return
			}
		}
		c.JSON(http.StatusOK, groupUser.GroupMemberDto())
	})

	rg.DELETE("/users/:userId/groups/:groupId", func(c *gin.Context) {
		userId := c.Param("userId")
		var user models.User
		if middlewares.TenantDB(c).First(&user, "id = ?", userId).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		groupId := c.Param("groupId")
		var groupUser models.GroupUser
		if middlewares.TenantDB(c).First(&groupUser, "user_id = ? AND group_id = ?", user.Id, groupId).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		if middlewares.TenantDB(c).Delete(&groupUser).Error != nil {
			c.Status(http.StatusInternalServerError)
		}
		c.Status(http.StatusNoContent)
	})
}
