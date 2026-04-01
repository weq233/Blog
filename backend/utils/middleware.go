package utils

import (
	"fmt"
	"github.com/beego/beego/v2/server/web/context"
	"strings"
	"net/http"
)

// JWTMiddleware JWT 验证中间件
func JWTMiddleware(ctx *context.Context) {
	// 设置 CORS 头部（确保所有响应都包含）
	ctx.Output.Header("Access-Control-Allow-Origin", "*")
	ctx.Output.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
	ctx.Output.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	ctx.Output.Header("Access-Control-Allow-Credentials", "true")
	
	// 跳过登录、注册、验证码接口
	if ctx.Request.URL.Path == "/api/login" || 
	   ctx.Request.URL.Path == "/api/register" ||
	   strings.HasPrefix(ctx.Request.URL.Path, "/api/captcha") ||
	   strings.HasPrefix(ctx.Request.URL.Path, "/admin/login") {
		return
	}
	
	// 跳过公开的文章接口（列表、详情、分类、标签）
	// 注意：评论接口（POST /api/article/:id/comments）需要认证，不能跳过
	isPublicArticle := (strings.HasPrefix(ctx.Request.URL.Path, "/api/articles") && 
	                   !strings.HasPrefix(ctx.Request.URL.Path, "/api/articles/create") &&
	                   !strings.HasPrefix(ctx.Request.URL.Path, "/api/articles/my")) ||
	                   (strings.HasPrefix(ctx.Request.URL.Path, "/api/article/") && 
	                    ctx.Request.Method == "GET" &&
	                    !strings.Contains(ctx.Request.URL.Path, "/comments")) ||
	                   strings.HasPrefix(ctx.Request.URL.Path, "/api/categories") ||
	                   strings.HasPrefix(ctx.Request.URL.Path, "/api/tags")
	
	if isPublicArticle {
		return
	}

	// 从 Header 中获取 token
	tokenString := ctx.Input.Header("Authorization")
	
	if tokenString == "" {
		// 尝试从 Cookie 中获取
		cookie := ctx.GetCookie("jwt_token")
		if cookie != "" {
			tokenString = cookie
		}
	}
	
	// 去除 Bearer 前缀
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	
	if tokenString == "" {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]interface{}{
			"code":    401,
			"message": "未提供认证令牌",
		}, false, false)
		return
	}
	
	// 验证 token
	claims, err := ParseToken(tokenString)
	if err != nil {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]interface{}{
			"code":    401,
			"message": "无效的认证令牌",
		}, false, false)
		return
	}
	
	// 将用户信息存入上下文
	ctx.Input.SetParam("userID", fmt.Sprintf("%d", claims.UserId))
	ctx.Input.SetParam("username", claims.Username)
	ctx.Input.SetParam("user_role", fmt.Sprintf("%d", claims.Role))
}

// AdminOnlyMiddleware 仅管理员访问的中间件
func AdminOnlyMiddleware(ctx *context.Context) {
	// 设置 CORS 头部
	ctx.Output.Header("Access-Control-Allow-Origin", "*")
	ctx.Output.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
	ctx.Output.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	
	roleStr := ctx.Input.Param("user_role")
	if roleStr != "2" {
		ctx.Output.SetStatus(http.StatusForbidden)
		ctx.Output.JSON(map[string]interface{}{
			"code":    403,
			"message": "需要管理员权限",
		}, false, false)
		return
	}
}
