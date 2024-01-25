package internal

import (
	"errors"
	"github.com/gin-gonic/gin"
	"reflect"
	"strconv"
)

func BindUri(c *gin.Context, obj any) *Api {
	a := New(c)
	return a.BindUri(obj)
}

func BindJson(c *gin.Context, obj any) *Api {
	a := New(c)
	return a.BindJson(obj)
}

func BindUriAndJson(c *gin.Context, obj any) *Api {
	a := New(c)
	return a.BindUriAndJson(obj)
}

func setUriValue(c *gin.Context, obj any) error {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		fd := t.Field(i)

		if tag := fd.Tag.Get("uri"); tag != "" {
			uri := c.Param(tag)
			if uri == "" {
				if required := fd.Tag.Get("binding"); required != "required" {
					continue
				}
				return errors.New(tag + " uri parameter should not be empty")
			}
			field := reflect.ValueOf(obj).Elem().Field(i)
			switch field.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				num, err := strconv.ParseInt(uri, 10, 64)
				if err != nil {
					return errors.New("failed to parse " + tag + " parameter as integer")
				}
				field.SetInt(num)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				num, err := strconv.ParseUint(uri, 10, 64)
				if err != nil {
					return errors.New("failed to parse " + tag + " parameter as unsigned integer")
				}
				field.SetUint(num)
			case reflect.Float32, reflect.Float64:
				num, err := strconv.ParseFloat(uri, 64)
				if err != nil {
					return errors.New("failed to parse " + tag + " parameter as float")
				}
				field.SetFloat(num)
			case reflect.String:
				field.SetString(uri)
			default:
				return errors.New("unsupported field type for " + tag + " parameter")
			}
		}
	}
	return nil
}
