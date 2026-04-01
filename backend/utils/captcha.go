package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// 验证码存储器（简单的内存存储）
type CaptchaStore struct {
	codes map[string]string
}

var captchaStore = &CaptchaStore{
	codes: make(map[string]string),
}

// CaptchaConfig 验证码配置
type CaptchaConfig struct {
	CaptchaLen    int
	MathChallenge bool // 是否使用数学计算题
}

// generateRandomString 生成随机字符串
func generateRandomString(length int) (string, error) {
	const charset = "234567890abcdefghjkmnpqrstuvwxyz"
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

// GenerateCaptcha 生成验证码（返回 captcha_id 和答案）
func GenerateCaptcha(id string, config *CaptchaConfig) (captchaId, answer string, err error) {
	if config == nil {
		config = &CaptchaConfig{
			CaptchaLen:    4,
			MathChallenge: false,
		}
	}
	
	// 如果没有提供 id，生成一个新的
	if id == "" {
		idBytes := make([]byte, 16)
		rand.Read(idBytes)
		id = fmt.Sprintf("%x", idBytes)
	}
	
	// 生成验证码内容
	if config.MathChallenge {
		// 数学计算题
		a, _ := rand.Int(rand.Reader, big.NewInt(10))
		b, _ := rand.Int(rand.Reader, big.NewInt(10))
		answer = fmt.Sprintf("%d", a.Int64()+b.Int64())
		captchaStore.codes[id] = answer
	} else {
		// 字符验证码
		answer, err = generateRandomString(config.CaptchaLen)
		if err != nil {
			return "", "", err
		}
		captchaStore.codes[id] = answer
	}
	
	return id, answer, nil
}

// VerifyCaptcha 验证验证码
func VerifyCaptcha(id, answer string, clear bool) bool {
	storedAnswer, exists := captchaStore.codes[id]
	if !exists {
		return false
	}
	
	if clear {
		delete(captchaStore.codes, id)
	}
	
	return storedAnswer == answer
}

// GenerateMathCaptcha 生成数学计算题验证码
func GenerateMathCaptcha() (id, question, answer string, err error) {
	idBytes := make([]byte, 16)
	rand.Read(idBytes)
	id = fmt.Sprintf("%x", idBytes)
	
	// 生成简单的加法题
	a, _ := rand.Int(rand.Reader, big.NewInt(10))
	b, _ := rand.Int(rand.Reader, big.NewInt(10))
	answer = fmt.Sprintf("%d", a.Int64()+b.Int64())
	question = fmt.Sprintf("%d + %d = ?", a.Int64(), b.Int64())
	
	captchaStore.codes[id] = answer
	
	return id, question, answer, nil
}
