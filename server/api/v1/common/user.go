package common

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/global/shared"
	"github.com/imehc/do-exercise/server/model/common/request"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/common/status"
	"github.com/imehc/do-exercise/server/util"
	"github.com/redis/go-redis/v9"
)

type UserApi struct{}

const (
	BindEmailPrefix      = "bind_email_"
	RebindEmailPrefix    = "rebind_email_"
	ModifyPasswordPrefix = "modify_password_"
	ForgotPasswordPrefix = "forgot_password_"
)

// checkEmail 检查邮箱是否是本人或者已发送验证码
func (s *UserApi) checkEmail(ctx *gin.Context, context context.Context, iRedis *redis.Client, userId string, emailType string) (*string, *string, error) {
	req := &request.Email{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		return nil, nil, errors.New("invalidParameter")
	}
	redisKey := fmt.Sprintf("%s%s", emailType, req.Email)
	emailInfo, err := iRedis.Get(context, redisKey).Result()
	if err != nil {
		if err == redis.Nil {
			return &req.Email, nil, nil
		}
		return nil, nil, errors.New("emailSendFailed")
	}
	// 解析email信息
	var emailData request.EmailCache
	if err := json.Unmarshal([]byte(emailInfo), &emailData); err != nil {
		return nil, nil, errors.New("emailSendFailed")
	}

	if emailData.UserId == userId {
		return nil, &emailData.Code, errors.New("emailSendLimit")
	} else {
		return nil, nil, errors.New("emailExists")
	}
}

// sendEmail 发送邮件
func (s *UserApi) sendEmail(from, to string, data util.EmailData) error {
	return shared.Email.SendEmail(from, to, data.EmailTitle, data)
}

// getEmailCache 获取缓存的邮箱验证码
func (s *UserApi) getEmailCache(ctx context.Context, iRedis *redis.Client, emailType, email string) (*request.EmailCache, error) {
	emailInfo, err := iRedis.Get(ctx, fmt.Sprintf("%s%s", emailType, email)).Result()
	if err != nil {
		return nil, errors.New("captchaNotExist")
	}
	var emailData request.EmailCache
	if err := json.Unmarshal([]byte(emailInfo), &emailData); err != nil {
		return nil, errors.New("captchaNotExist")
	}
	return &emailData, nil
}

// setEmailCache 将验证码存入redis
func (s *UserApi) setEmailCache(ctx context.Context, iRedis *redis.Client, emailType, email, code string, userId string, minutes time.Duration) error {
	emailCache := request.EmailCache{
		UserId: userId,
		Code:   code,
	}
	emailCacheJson, err := json.Marshal(emailCache)
	if err != nil {
		return err
	}
	if err := iRedis.Set(ctx, fmt.Sprintf("%s%s", emailType, email), emailCacheJson, minutes).Err(); err != nil {
		return err
	}

	return nil
}

// clearEmailCache 清除邮箱验证码
func (s *UserApi) clearEmailCache(ctx context.Context, iRedis *redis.Client, emailType, email string) error {
	if err := iRedis.Del(ctx, fmt.Sprintf("%s%s", emailType, email)).Err(); err != nil {
		return err
	}
	return nil
}

// sendEmailCode 发送邮箱验证码
func (s *UserApi) sendEmailCode(ctx *gin.Context, emailType string, data *util.EmailData) {
	lang := ctx.GetString("lang")
	iRedis := global.Redis
	context := context.Background()
	userId := ctx.MustGet("userId").(string)

	to, _, err := s.checkEmail(ctx, context, iRedis, userId, emailType)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}

	minute := time.Duration(10) * time.Minute
	code := util.GenerateRandomNumber(6)
	if err := s.setEmailCache(context, iRedis, emailType, *to, code, userId, minute); err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate("emailSendFailed", lang),
		})
		return
	}

	data.Minute = int(minute.Minutes())
	data.VerificationCode = code

	if err := s.sendEmail(global.Config.Email.User, *to, *data); err != nil {
		err := s.clearEmailCache(context, iRedis, emailType, *to)
		if err != nil {
			response.BadRequest(ctx, response.ValidationError{
				Type:    status.BAD_REQUEST_MSG,
				Message: global.I18.Translate("emailSendFailed", lang),
			})
			return
		}
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate("emailSendFailed", lang),
		})
		return
	}
	response.NoContent(ctx)
}

// bindEmail 绑定邮箱
func (s *UserApi) bindEmail(ctx *gin.Context, emailType string) {
	lang := ctx.GetString("lang")
	iRedis := global.Redis
	context := context.Background()
	userId := ctx.MustGet("userId").(string)

	req := &request.BindEmailReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	cache, err := s.getEmailCache(context, iRedis, emailType, req.Email)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}

	defer func() {
		_ = s.clearEmailCache(context, iRedis, emailType, req.Email)
	}()

	if cache.Code != req.Code || cache.UserId != userId {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate("captchaError", lang),
		})
		return
	}

	if err = userService.BindEmail(request.BindEmailReq{
		Id:    userId,
		Email: req.Email,
	}); err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}
	response.NoContent(ctx)
}

// SendBindEmailCode 发送绑定邮箱验证码
func (s *UserApi) SendBindEmailCode(ctx *gin.Context) {
	s.sendEmailCode(ctx, BindEmailPrefix, &util.EmailData{
		EmailTitle:       "验证码",
		VerificationType: "绑定邮箱",
		GreetingText:     "感谢您使用我们的服务，请使用以下验证码完成邮箱绑定：",
	})
}

// SendRebindEmailCode 发送换绑邮箱验证码
func (s *UserApi) SendRebindEmailCode(ctx *gin.Context) {
	s.sendEmailCode(ctx, RebindEmailPrefix, &util.EmailData{
		EmailTitle:       "验证码",
		VerificationType: "绑定新邮箱",
		GreetingText:     "感谢您使用我们的服务，请使用以下验证码完成新邮箱绑定：",
	})
}

// SendModifyPasswordCode 发送修改密码验证码
func (s *UserApi) SendModifyPasswordCode(ctx *gin.Context) {
	s.sendEmailCode(ctx, ModifyPasswordPrefix, &util.EmailData{
		EmailTitle:       "验证码",
		VerificationType: "重置密码",
		GreetingText:     "感谢您使用我们的服务，请使用以下验证码完成密码重置：",
	})
}

// BindEmail 绑定邮箱
func (s *UserApi) BindEmail(ctx *gin.Context) {
	s.bindEmail(ctx, BindEmailPrefix)
}

// RebindEmail 换绑邮箱
func (s *UserApi) RebindEmail(ctx *gin.Context) {
	s.bindEmail(ctx, RebindEmailPrefix)
}

// UpdatePassword 修改密码
func (s *UserApi) UpdatePassword(ctx *gin.Context) {
	lang := ctx.GetString("lang")
	userId := ctx.MustGet("userId").(string)

	req := &request.UserModifyPasswordReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	oldPassword, err := shared.RSACrypto.DecryptWithKey(req.PublicKey, req.OldPassword)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}
	password, err := shared.RSACrypto.DecryptWithKey(req.PublicKey, req.Password)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}
	req.Id = userId
	req.OldPassword = oldPassword
	req.Password = password

	if err := userService.UpdatePassword(*req); err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}
	response.NoContent(ctx)
}

// UpdateProfile 修改用户基本信息
func (s *UserApi) UpdateProfile(ctx *gin.Context) {
	lang := ctx.GetString("lang")
	userId := ctx.MustGet("userId").(string)

	var req request.UserModifyProfileReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	req.Id = userId

	err := userService.UpdateProfile(req)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}
	response.NoContent(ctx)
}

// GetProfile 获取用户基本信息
func (s *UserApi) GetProfile(ctx *gin.Context) {
	lang := ctx.GetString("lang")
	userId := ctx.MustGet("userId").(string)

	user, err := userService.GetProfile(userId)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}
	response.Success(ctx, user)
}

// GetMenu 获取用户菜单
func (s *UserApi) GetMenu(ctx *gin.Context) {
	lang := ctx.GetString("lang")
	userId := ctx.MustGet("userId").(string)

	menu, err := userService.GetMenu(userId)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}
	response.Success(ctx, menu)
}
