package reset

import (
	"alfred/internal/controller/internal"
	"alfred/internal/endpoint/resp"
	"alfred/internal/model"
	"alfred/internal/service"
	"alfred/internal/service/reset"
	"alfred/pkg/global"
	"alfred/pkg/middlewares"
	"alfred/pkg/utils"
	"github.com/gin-gonic/gin"
	"time"
)

// SmsAvailable 校验租户是否配置SMS
//
//	@Summary	check sms available
//	@Schemes
//	@Description	check sms available
//	@Tags			reset
//	@Param			tenant		path		string	true	"tenant"	default(default)
//	@Success		200
//	@Router			/accounts/{tenant}/reset/smsAvailable [get]
func SmsAvailable(c *gin.Context) {
	tenant := internal.GetTenant(c)
	var provider model.Provider
	if err := global.DB.First(&provider, "tenant_id = ? AND name = ?", tenant.Id, "sms").Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get provider err")
		return
	}

	//查询在sms_connectors表中是否存在对应的TenantId的记录
	var connectors []model.SmsConnector
	if err := global.DB.Where("tenant_id = ?", tenant.Id).Find(&connectors).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "get connectors err")
		return
	}

	//同时满足两个条件才返回true
	if len(connectors) > 0 && provider.Type == "sms" {
		resp.SuccessWithData(c, true)
		return
	} else {
		resp.ErrorRequestWithMsg(c, "sms is not available")
	}
}

// VerifyResetPasswordRequest ForgotPassword godoc 发起忘记密码请求
//
//	@Summary	forgot password
//	@Schemes
//	@Description	forgot password
//	@Tags			reset
//	@Param			tenant		path		string	true	"tenant"	default(default)
//	@Param			verifyMethod	formData	string	true	"verify method"
//	@Param			passCodePayload	formData	string	true	"pass code payload"
//	@Param	        areaCode	formData	string	false	"area code" default(+86)
//	@Success		200
//	@Router			/accounts/{tenant}/reset/verifyResetPasswordRequest [post]
func VerifyResetPasswordRequest(c *gin.Context) {
	verifyMethod := c.PostForm("verifyMethod")
	tenant := internal.GetTenant(c)
	if verifyMethod == "phonePassCode" {
		resetProvider, err := reset.GetResetAuthProvider(tenant.Id, "sms")
		phoneNumber := c.PostForm("passCodePayload")
		areaCode := c.PostForm("areaCode")
		fullPhoneNumber := areaCode + phoneNumber

		//先查询user表中是否有对应的手机号
		var user model.User
		if err = global.DB.First(&user, "phone = ?", phoneNumber).Error; err != nil {
			resp.ErrorSqlFirst(c, err, "user not exist")
			return
		}

		_, code, err := resetProvider.ResetAuth(fullPhoneNumber, tenant.Id)
		if err != nil {
			resp.ErrorUnknown(c, err, "provider auth err")
			return
		}
		// 获取重置密码令牌
		token, err := service.GetResetPasswordToken(c, fullPhoneNumber)
		if err != nil {
			resp.ErrorUnknown(c, err, "provider auth err")
			return
		}

		//将Code和Token存入token_codes表中
		if err = global.DB.Create(&model.TokenCode{
			Code:      code,
			Token:     token,
			CreatedAt: time.Now(),
			ClientId:  tenant.Name,
			TenantId:  tenant.Id,
			Sub:       phoneNumber,
			Type:      "resetPassword",
		}).Error; err != nil {
			resp.ErrorSqlCreate(c, err, "create token code err")
			return
		}

	} else if verifyMethod == "emailPassCode" {
		// TODO: 处理电子邮件验证码的验证逻辑
		resp.ErrorRequestWithMsg(c, "verify method not supported")
	} else {
		resp.ErrorRequestWithMsg(c, "verify method not supported")
	}
}

// GetResetPasswordToken 生成一次性的重置密码的token
//
//	@Summary	get reset password token
//	@Schemes
//	@Description	get reset password token
//	@Tags			reset
//	@Param			tenant		path		string	true	"tenant"	default(default)
//	@Param			code		formData	string	true	"code"
//	@Success		200
//	@Router			/accounts/{tenant}/reset/getResetPasswordToken [post]
func GetResetPasswordToken(c *gin.Context) {
	//清除token_codes表中的resetPassword类型的过期记录
	service.ClearResetPasswordTokenCode("resetPassword")

	code := c.PostForm("code")
	var tokenCode model.TokenCode
	if err := global.DB.First(&tokenCode, "code = ?", code).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get token code err , code not exist")
		return
	}

	//清除当前code对应的token_codes表中的记录
	service.ClearTokenCode(code)

	//将对应的token返回
	resp.SuccessWithData(c, tokenCode.Token)
}

// PasswordReset godoc 重置密码
//
//	@Summary	reset password
//	@Schemes
//	@Description	reset password
//	@Tags			reset
//	@Param			Authorization	header	string	true	"token"
//	@Param			tenant		path		string	true	"tenant"	default(default)
//	@Param          newPassword formData	string	true	"new password"
//	@Success		200
//	@Router			/accounts/{tenant}/reset/resetPassword [post]
func PasswordReset(c *gin.Context) {
	tenant := internal.GetTenant(c)
	phone, exist := c.Get("sub")
	if !exist {
		resp.ErrorUnauthorized(c, nil, "sub not exist")
		return
	}
	//去掉phone中的+86
	phone = phone.(string)[3:]

	//获取用户信息
	user, err := service.GetUserByPhone(phone.(string), tenant.Id)
	if err != nil {
		resp.ErrorSqlSelect(c, err, "get user err")
		return
	}
	//获取新密码
	newPassword := c.PostForm("newPassword")
	if newPassword == "" {
		resp.ErrorRequestWithMsg(c, "password should not be null")
		return
	}

	hash, err := utils.HashPassword(newPassword)
	if err != nil {
		resp.ErrorUnknown(c, err, "password hash err")
		return
	}

	user.PasswordHash = hash
	if err = internal.TenantDB(c).Save(&user).Error; err != nil {
		resp.ErrorSqlUpdate(c, err, "update user err")
		return
	}

	resp.Success(c)
}

func AddResetRouter(rg *gin.RouterGroup) {
	rg.GET("/reset/smsAvailable", SmsAvailable)
	rg.POST("/reset/verifyResetPasswordRequest", VerifyResetPasswordRequest)
	rg.POST("/reset/getResetPasswordToken", GetResetPasswordToken)
	rg.POST("/reset/resetPassword", middlewares.ResetPasswordToken, PasswordReset)
}
