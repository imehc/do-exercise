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
	"github.com/imehc/do-exercise/server/util"
	"github.com/redis/go-redis/v9"
)

type UserApi struct{}

const (
	BindEmailPrefix      = "bind_email_"
	RebindEmailPrefix    = "rebind_email_"
	ModifyPasswordPrefix = "modify_password_"
	ForgotPasswordPrefix = "forgot_password_"
	LoginWithEmailPrefix = "login_with_email_"
)

// userNotFoundErr 与 service 层返回的错误标识保持一致，同时也是 i18n 的翻译 key。
const userNotFoundErr = "userNotFound"

// checkEmail 检查邮箱是否是本人或者已发送验证码。
// userIds 是本次验证码的候选账号集合（匿名邮箱流程下同一邮箱可能命中多个租户的账号）。
func (s *UserApi) checkEmail(ctx *gin.Context, context context.Context, iRedis *redis.Client, userIds []string, emailType string) (*string, *string, error) {
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

	if emailData.SharesUser(userIds) {
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
func (s *UserApi) setEmailCache(ctx context.Context, iRedis *redis.Client, emailType, email, code string, userIds []string, minutes time.Duration) error {
	emailCache := request.EmailCache{
		UserIds: userIds,
		Code:    code,
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

// checkEmailSendLimit 按邮箱地址限流验证码发送，防止邮件轰炸。
// 每个邮箱 1 小时内最多发送 maxSendsPerHour 封，超出后与 cache 仍在生效时的行为一致。
func (s *UserApi) checkEmailSendLimit(context context.Context, iRedis *redis.Client, email string) error {
	const (
		key              = "email_send_limit_%s"
		maxSendsPerHour  = 5
		limitWindowHours = 1
	)
	cnt, err := iRedis.Incr(context, fmt.Sprintf(key, email)).Result()
	if err != nil {
		return err
	}
	if cnt == 1 {
		iRedis.Expire(context, fmt.Sprintf(key, email), time.Duration(limitWindowHours)*time.Hour)
	}
	if cnt > maxSendsPerHour {
		return errors.New("emailSendLimit")
	}
	return nil
}

// sendEmailCode 发送邮箱验证码。
// userIds 为空表示「当前登录用户自己」的流程（绑定/换绑/改密），从会话取。
func (s *UserApi) sendEmailCode(ctx *gin.Context, emailType string, data *util.EmailData, userIds []string) {
	iRedis := global.Redis
	context := context.Background()

	uIds := userIds
	if len(uIds) == 0 {
		uIds = []string{ctx.MustGet("userId").(string)}
	}

	to, _, err := s.checkEmail(ctx, context, iRedis, uIds, emailType)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	minute := time.Duration(10) * time.Minute
	code := util.GenerateRandomNumber(6)
	if err := s.checkEmailSendLimit(context, iRedis, *to); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	if err := s.setEmailCache(context, iRedis, emailType, *to, code, uIds, minute); err != nil {
		response.BadRequest(ctx, "emailSendFailed")
		return
	}

	data.Minute = int(minute.Minutes())
	data.VerificationCode = code

	if err := s.sendEmail(global.Config.Email.User, *to, *data); err != nil {
		err := s.clearEmailCache(context, iRedis, emailType, *to)
		if err != nil {
			response.BadRequest(ctx, "emailSendFailed")
			return
		}
		response.BadRequest(ctx, "emailSendFailed")
		return
	}
	response.NoContent(ctx)
}

// bindEmail 绑定邮箱
func (s *UserApi) bindEmail(ctx *gin.Context, emailType string) {
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
		response.BadRequest(ctx, err.Error())
		return
	}

	defer func() {
		_ = s.clearEmailCache(context, iRedis, emailType, req.Email)
	}()

	if cache.Code != req.Code || !cache.Allows(userId) {
		response.BadRequest(ctx, "captchaError")
		return
	}

	if err = userService.BindEmail(util.DB(ctx), userId, req.Email); err != nil {
		response.BadRequest(ctx, err.Error())
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
	}, nil)
}

// SendRebindEmailCode 发送换绑邮箱验证码
func (s *UserApi) SendRebindEmailCode(ctx *gin.Context) {
	s.sendEmailCode(ctx, RebindEmailPrefix, &util.EmailData{
		EmailTitle:       "验证码",
		VerificationType: "绑定新邮箱",
		GreetingText:     "感谢您使用我们的服务，请使用以下验证码完成新邮箱绑定：",
	}, nil)
}

// SendModifyPasswordCode 发送修改密码验证码
func (s *UserApi) SendModifyPasswordCode(ctx *gin.Context) {
	s.sendEmailCode(ctx, ModifyPasswordPrefix, &util.EmailData{
		EmailTitle:       "验证码",
		VerificationType: "重置密码",
		GreetingText:     "感谢您使用我们的服务，请使用以下验证码完成密码重置：",
	}, nil)
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
	userId := ctx.MustGet("userId").(string)

	req := &request.UserModifyPasswordReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	oldPassword, err := shared.RSACrypto.DecryptWithKey(req.PublicKey, req.OldPassword, false)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	password, err := shared.RSACrypto.DecryptWithKey(req.PublicKey, req.Password, true)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	if err := userService.UpdatePassword(util.DB(ctx), userId, oldPassword, password, ctx.GetString("accessToken")); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.NoContent(ctx)
}

// UpdateProfile 修改用户基本信息
func (s *UserApi) UpdateProfile(ctx *gin.Context) {
	userId := ctx.MustGet("userId").(string)

	var req request.UserModifyProfileReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	err := userService.UpdateProfile(util.DB(ctx), userId, req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.NoContent(ctx)
}

// GetProfile 获取用户基本信息
func (s *UserApi) GetProfile(ctx *gin.Context) {
	userId := ctx.MustGet("userId").(string)

	user, err := userService.GetProfile(util.DB(ctx), userId)
	if err != nil {
		// userId 取自已验签的 token，查不到说明账号已被删除，属资源不存在而非参数错误。
		// 不用 401：前端对 401 会刷新 token 后重试，对不存在的账号会陷入循环。
		if err.Error() == userNotFoundErr {
			response.NotFoundMessage(ctx, err.Error())
			return
		}
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, user)
}

// GetMenu 获取用户菜单
func (s *UserApi) GetMenu(ctx *gin.Context) {
	userId := ctx.MustGet("userId").(string)

	menu, err := userService.GetMenu(util.DB(ctx), userId)
	if err != nil {
		if err.Error() == userNotFoundErr {
			response.NotFoundMessage(ctx, err.Error())
			return
		}
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, menu)
}
