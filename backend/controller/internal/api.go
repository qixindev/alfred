package internal

import (
	"errors"
	"github.com/gin-gonic/gin"
)

type Api struct {
	c     *gin.Context
	Error error
}

func New(c *gin.Context) *Api {
	return &Api{
		c: c,
	}
}

func (a *Api) setError(err error) *Api {
	a.Error = err
	return a
}

func (a *Api) SetCtx(c *gin.Context) *Api {
	if a.c == nil {
		a.c = c
	}
	return a
}

func (a *Api) BindUri(obj any) *Api {
	if a.c == nil {
		return a.setError(errors.New("gin context should not be nil"))
	}
	if err := setUriValue(a.c, obj, "uri"); err != nil {
		return a.setError(err)
	}
	return a
}

func (a *Api) BindQuery(obj any) *Api {
	if a.c == nil {
		return a.setError(errors.New("gin context should not be nil"))
	}
	if err := a.c.ShouldBindQuery(obj); err != nil {
		return a.setError(err)
	}
	return a
}

func (a *Api) BindJson(obj any) *Api {
	if a.c == nil {
		return a.setError(errors.New("gin context should not be nil"))
	}
	if err := a.c.ShouldBindJSON(obj); err != nil {
		return a.setError(err)
	}
	return a
}

func (a *Api) BindUriAndJson(obj any) *Api {
	if a.c == nil {
		return a.setError(errors.New("gin context should not be nil"))
	}
	if err := a.c.ShouldBindJSON(obj); err != nil {
		return a.setError(err)
	}
	if err := setUriValue(a.c, obj, "uri"); err != nil {
		return a.setError(err)
	}
	return a
}
