package internal

import (
	"accounts/internal/model"
	"errors"
	"github.com/gin-gonic/gin"
)

type Api struct {
	c      *gin.Context
	Error  error
	Tenant *model.Tenant
}

func New(c *gin.Context) *Api {
	return &Api{
		c: c,
	}
}

func (a *Api) setError(err error) *Api {
	if a.Error != nil {
		a.Error = err
	}
	return a
}

func (a *Api) SetCtx(c *gin.Context) *Api {
	if a.c == nil {
		a.c = c
	}
	return a
}

func (a *Api) SetTenant() *Api {
	if a.c == nil {
		return a.setError(errors.New("gin context should not be nil"))
	}
	t := a.c.MustGet("tenant")
	tenant, ok := t.(*model.Tenant)
	if !ok {
		return a
	}
	a.Tenant = tenant
	return a
}
