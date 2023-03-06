package admin

import (
	"accounts/data"
	"accounts/middlewares"
	"accounts/models"
	"accounts/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func addAdminProvidersRoutes(rg *gin.RouterGroup) {
	rg.GET("/providers", func(c *gin.Context) {
		var providers []models.Provider
		if middlewares.TenantDB(c).Find(&providers).Error != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, utils.Filter(providers, models.Provider2Dto))
	})
	rg.GET("/providers/:providerId", func(c *gin.Context) {
		providerId := c.Param("providerId")
		var provider models.Provider
		if middlewares.TenantDB(c).First(&provider, "id = ?", providerId).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, provider.Dto())
	})
	rg.POST("/providers", func(c *gin.Context) {
		tenant := middlewares.GetTenant(c)
		var provider models.Provider
		err := c.BindJSON(&provider)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		provider.TenantId = tenant.Id
		if data.DB.Create(&provider).Error != nil {
			c.Status(http.StatusConflict)
			return
		}
		c.JSON(http.StatusOK, provider.Dto())
	})
	rg.PUT("/providers/:providerId", func(c *gin.Context) {
		providerId := c.Param("providerId")
		var provider models.Provider
		if middlewares.TenantDB(c).First(&provider, "id = ?", providerId).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		var p models.Provider
		err := c.BindJSON(&p)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		provider.Name = p.Name
		if data.DB.Save(&provider).Error != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, provider.Dto())
	})
	rg.DELETE("/providers/:providerId", func(c *gin.Context) {
		providerId := c.Param("providerId")
		var provider models.Provider
		if middlewares.TenantDB(c).First(&provider, "id = ?", providerId).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		if data.DB.Delete(&provider).Error != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusNoContent)
	})
}
