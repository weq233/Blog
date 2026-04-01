package controllers

import (
	"blog-system/models"
	"blog-system/utils"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	"golang.org/x/crypto/bcrypt"
)

// 认证控制器 - API
type AuthController struct {
	web.Controller
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	CaptchaId  string `json:"captcha_id"`
	CaptchaAns string `json:"captcha_ans"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// GetCurrentUser 获取当前登录用户信息
func (c *AuthController) GetCurrentUser() {
	// 从 JWT 中间件中获取用户 ID
	userIDStr := c.Ctx.Input.Param("userID")

	if userIDStr == "" {
		c.Data["json"] = map[string]interface{}{
			"code":    401,
			"message": "未登录或 token 无效",
		}
		c.ServeJSON()
		return
	}

	// 将字符串转换为 int
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    401,
			"message": "用户 ID 格式错误",
		}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	var user models.User
	err = o.QueryTable(new(models.User)).Filter("id", userID).One(&user)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    404,
			"message": "用户不存在",
		}
		c.ServeJSON()
		return
	}

	// 返回用户信息（不包含密码）
	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"id":        user.Id,
			"username":  user.Username,
			"email":     user.Email,
			"nickname":  user.Nickname,
			"avatar":    user.Avatar,
			"role":      user.Role,
			"status":    user.Status,
			"createdAt": user.CreatedAt,
		},
	}
	c.ServeJSON()
}

// Login 用户登录（带验证码和失败限制）
func (c *AuthController) Login() {
	var req LoginRequest

	body, err := io.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
		}
		c.ServeJSON()
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "JSON 解析失败",
		}
		c.ServeJSON()
		return
	}

	if !utils.VerifyCaptcha(req.CaptchaId, req.CaptchaAns, true) {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "验证码错误",
		}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	var user models.User

	err = o.QueryTable(new(models.User)).Filter("username", req.Username).One(&user)

	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    401,
			"message": "用户名或密码错误",
		}
		c.ServeJSON()
		return
	}

	// 检查账户是否被锁定
	if user.LockedUntil.After(time.Now()) {
		waitMinutes := int(time.Until(user.LockedUntil).Minutes()) + 1
		c.Data["json"] = map[string]interface{}{
			"code":    423,
			"message": fmt.Sprintf("账户已被锁定，请 %d 分钟后再试", waitMinutes),
		}
		c.ServeJSON()
		return
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		// 密码错误，增加失败次数
		user.FailedAttempts++
		if user.FailedAttempts >= 5 {
			// 锁定账户 10 分钟
			user.LockedUntil = time.Now().Add(10 * time.Minute)
			user.FailedAttempts = 0
		}
		o.Update(&user, "FailedAttempts", "LockedUntil")

		c.Data["json"] = map[string]interface{}{
			"code":    401,
			"message": "用户名或密码错误",
		}
		c.ServeJSON()
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		c.Data["json"] = map[string]interface{}{
			"code":    403,
			"message": "账户已被禁用",
		}
		c.ServeJSON()
		return
	}

	// 登录成功，重置失败次数
	user.FailedAttempts = 0
	user.LockedUntil = time.Time{}
	o.Update(&user, "FailedAttempts", "LockedUntil")

	// 生成 JWT 令牌
	token, err := utils.GenerateToken(user.Id, user.Username, user.Role)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "令牌生成失败",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"token": token,
			"user": map[string]interface{}{
				"id":       user.Id,
				"username": user.Username,
				"nickname": user.Nickname,
				"email":    user.Email,
				"role":     user.Role,
			},
		},
		"message": "登录成功",
	}
	c.ServeJSON()
}

// Register 用户注册
func (c *AuthController) Register() {
	var req RegisterRequest

	body, err := io.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
		}
		c.ServeJSON()
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "JSON 解析失败",
		}
		c.ServeJSON()
		return
	}

	// 参数验证
	if req.Username == "" || req.Password == "" || req.Email == "" {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "用户名、密码和邮箱不能为空",
		}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()

	// 检查用户名是否已存在
	var existingUser models.User
	err = o.QueryTable(new(models.User)).Filter("username", req.Username).One(&existingUser)
	if err == nil {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "用户名已存在",
		}
		c.ServeJSON()
		return
	}

	// 检查邮箱是否已注册
	err = o.QueryTable(new(models.User)).Filter("email", req.Email).One(&existingUser)
	if err == nil {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "邮箱已被注册",
		}
		c.ServeJSON()
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "密码加密失败",
		}
		c.ServeJSON()
		return
	}

	user := &models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Nickname: req.Nickname,
		Role:     1, // 默认普通用户
		Status:   1,
	}

	_, err = o.Insert(user)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "注册失败",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"code":    0,
		"message": "注册成功",
	}
	c.ServeJSON()
}

// ChangePassword 修改密码（需要登录）
func (c *AuthController) ChangePassword() {
	var req ChangePasswordRequest

	body, err := io.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
		}
		c.ServeJSON()
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "JSON 解析失败",
		}
		c.ServeJSON()
		return
	}

	// 获取当前用户
	username := c.Ctx.Input.Param("username")
	if username == "" {
		c.Data["json"] = map[string]interface{}{
			"code":    401,
			"message": "未登录",
		}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	var user models.User
	err = o.QueryTable(new(models.User)).Filter("username", username).One(&user)

	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    404,
			"message": "用户不存在",
		}
		c.ServeJSON()
		return
	}

	// 验证旧密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "原密码错误",
		}
		c.ServeJSON()
		return
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "密码加密失败",
		}
		c.ServeJSON()
		return
	}

	user.Password = string(hashedPassword)
	o.Update(&user, "Password")

	c.Data["json"] = map[string]interface{}{
		"code":    0,
		"message": "密码修改成功",
	}
	c.ServeJSON()
}

// GetCaptcha 获取验证码图片
func (c *AuthController) GetCaptcha() {
	captchaId := c.GetString("captcha_id", "")

	// 生成验证码
	id, _, err := utils.GenerateCaptcha(captchaId, &utils.CaptchaConfig{
		CaptchaLen:    4,
		MathChallenge: false,
	})
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "验证码生成失败",
		}
		c.ServeJSON()
		return
	}

	// TODO: 实际项目中应该生成真实的验证码图片
	// 这里为了演示，只返回 captcha_id，前端可以使用第三方服务生成图片
	// 或者使用 canvas 自己绘制

	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"captcha_id": id,
			// 实际应该返回 base64 编码的图片
			// "image": "data:image/png;base64,xxx",
			"tip": "请使用 captcha_id 调用前端验证码组件生成图片",
		},
	}
	c.ServeJSON()
}

// GetMathCaptcha 获取数学计算题验证码
func (c *AuthController) GetMathCaptcha() {
	id, question, _, err := utils.GenerateMathCaptcha()
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "验证码生成失败",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"captcha_id": id,
			"question":   question,
		},
	}
	c.ServeJSON()
}

// GetUserStats 获取用户统计数据
func (c *AuthController) GetUserStats() {
	// 从 JWT 中间件中获取用户 ID
	userIDStr := c.Ctx.Input.Param("userID")

	if userIDStr == "" {
		c.Data["json"] = map[string]interface{}{
			"code":    401,
			"message": "未登录或 token 无效",
		}
		c.ServeJSON()
		return
	}

	// 将字符串转换为 int
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    401,
			"message": "用户 ID 格式错误",
		}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()

	// 统计文章数量
	articleCount, err := o.QueryTable(new(models.Article)).Filter("author_id", userID).Count()
	if err != nil {
		articleCount = 0
	}

	// 统计评论数量
	commentCount, err := o.QueryTable(new(models.Comment)).Filter("user_id", userID).Count()
	if err != nil {
		commentCount = 0
	}

	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"articles":  articleCount,
			"comments":  commentCount,
			"favorites": 0,
			"following": 0,
			"followers": 0,
		},
	}
	c.ServeJSON()
}
