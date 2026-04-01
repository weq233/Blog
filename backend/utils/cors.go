package utils

import (
	"github.com/beego/beego/v2/server/web/context"
	"net/http"
)

// CORSMiddleware CORS 跨域中间件
func CORSMiddleware(ctx *context.Context) {
	// 设置允许的源
	ctx.Output.Header("Access-Control-Allow-Origin", "*")
	
	// 允许的方法
	ctx.Output.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
	
	// 允许的头部
	ctx.Output.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	
	// 允许携带凭证
	ctx.Output.Header("Access-Control-Allow-Credentials", "true")
	
	// 预检请求缓存时间
	ctx.Output.Header("Access-Control-Max-Age", "86400")
	
	// 处理预检请求（OPTIONS）
	if ctx.Request.Method == http.MethodOptions {
		ctx.Output.SetStatus(http.StatusNoContent) // 204 No Content
		return
	}
}
