package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/smtp"
	"github.com/jordan-wright/email"
)

// 邮箱配置
type EmailConfig struct {
	SMTPServer string
	SMTPPort   string
	Username   string
	Password   string
	FromName   string
}

var emailConfig *EmailConfig

// InitEmailConfig 初始化邮箱配置
func InitEmailConfig(smtpServer, smtpPort, username, password, fromName string) {
	emailConfig = &EmailConfig{
		SMTPServer: smtpServer,
		SMTPPort:   smtpPort,
		Username:   username,
		Password:   password,
		FromName:   fromName,
	}
}

// generateRandomCode 生成随机验证码
func generateRandomCode(length int) (string, error) {
	const charset = "0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

// GenerateEmailCode 生成邮箱验证码
func GenerateEmailCode() (string, error) {
	return generateRandomCode(6)
}

// SendVerificationEmail 发送验证邮件
func SendVerificationEmail(toEmail, code string) error {
	if emailConfig == nil {
		return fmt.Errorf("邮箱配置未初始化")
	}

	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", emailConfig.FromName, emailConfig.Username)
	e.To = []string{toEmail}
	e.Subject = "博客系统 - 邮箱验证码"
	e.Text = []byte(fmt.Sprintf(`尊敬的用户：

您的邮箱验证码是：%s

验证码有效期为 10 分钟，请尽快使用。

如非本人操作，请忽略此邮件。

此致
敬礼
博客系统团队`, code))

	// 发送邮件
	auth := smtp.PlainAuth("", emailConfig.Username, emailConfig.Password, emailConfig.SMTPServer)
	addr := fmt.Sprintf("%s:%s", emailConfig.SMTPServer, emailConfig.SMTPPort)
	return e.Send(addr, auth)
}

// SendResetPasswordEmail 发送重置密码邮件
func SendResetPasswordEmail(toEmail, resetLink string) error {
	if emailConfig == nil {
		return fmt.Errorf("邮箱配置未初始化")
	}

	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", emailConfig.FromName, emailConfig.Username)
	e.To = []string{toEmail}
	e.Subject = "博客系统 - 重置密码"
	e.HTML = []byte(fmt.Sprintf(`尊敬的用户：

您请求重置密码，请点击以下链接重置密码：

<a href="%s" style="display:inline-block;padding:10px 20px;background:#007bff;color:white;text-decoration:none;border-radius:5px;">重置密码</a>

或者复制以下链接到浏览器：
%s

该链接有效期为 1 小时。

如非本人操作，请忽略此邮件。

此致
敬礼
博客系统团队`, resetLink, resetLink))

	auth := smtp.PlainAuth("", emailConfig.Username, emailConfig.Password, emailConfig.SMTPServer)
	addr := fmt.Sprintf("%s:%s", emailConfig.SMTPServer, emailConfig.SMTPPort)
	return e.Send(addr, auth)
}

// SendWelcomeEmail 发送欢迎邮件
func SendWelcomeEmail(toEmail, username string) error {
	if emailConfig == nil {
		return fmt.Errorf("邮箱配置未初始化")
	}

	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", emailConfig.FromName, emailConfig.Username)
	e.To = []string{toEmail}
	e.Subject = "欢迎加入博客系统！"
	e.HTML = []byte(fmt.Sprintf(`亲爱的 %s：

欢迎加入博客系统！🎉

您已经成功注册，现在可以开始发布文章、管理分类和标签了。

立即登录：<a href="http://localhost:8080/login" style="color:#007bff;">登录</a>

如有任何问题，请随时联系我们。

此致
敬礼
博客系统团队`, username))

	auth := smtp.PlainAuth("", emailConfig.Username, emailConfig.Password, emailConfig.SMTPServer)
	addr := fmt.Sprintf("%s:%s", emailConfig.SMTPServer, emailConfig.SMTPPort)
	return e.Send(addr, auth)
}
