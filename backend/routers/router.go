package routers

import (
	"blog-system/controllers"
	"github.com/beego/beego/v2/server/web"
	"blog-system/utils"
)

func init() {
    // CORS 中间件（全局启用）
    web.InsertFilter("*", web.BeforeRouter, utils.CORSMiddleware)
    
    // 认证接口（不需要验证）- 放在 JWT 中间件之前
    web.Router("/api/login", &controllers.AuthController{}, "post:Login")
    web.Router("/api/register", &controllers.AuthController{}, "post:Register")
    
    // 验证码接口（不需要验证）- 放在 JWT 中间件之前
    web.Router("/api/captcha", &controllers.AuthController{}, "get:GetCaptcha")
    web.Router("/api/captcha/math", &controllers.AuthController{}, "get:GetMathCaptcha")
    
    // 文章相关接口（公开访问，不需要 JWT 验证）- 放在 JWT 中间件之前
    web.Router("/api/articles", &controllers.ArticleController{}, "get:List")
    web.Router("/api/article/:id", &controllers.ArticleController{}, "get:View")
    web.Router("/api/articles/category/:slug", &controllers.ArticleController{}, "get:ByCategory")
    web.Router("/api/articles/tag/:slug", &controllers.ArticleController{}, "get:ByTag")
    web.Router("/api/categories", &controllers.ArticleController{}, "get:GetCategories")
    web.Router("/api/tags", &controllers.ArticleController{}, "get:GetTags")
    
    // 点赞和关注接口（部分需要 JWT）
    web.Router("/api/article/:id/like", &controllers.ArticleController{}, "post:Like")
    web.Router("/api/article/:id/like-status", &controllers.ArticleController{}, "get:GetLikeStatus")
    web.Router("/api/user/:id/follow", &controllers.ArticleController{}, "post:Follow")
    web.Router("/api/user/:id/follow-status", &controllers.ArticleController{}, "get:GetFollowStatus")
    
    // API 路由（需要 JWT 验证）- 放在认证接口之后
    web.InsertFilter("/api/*", web.BeforeRouter, utils.JWTMiddleware)
    
    // 需要验证的 API 接口
    web.Router("/api/user/me", &controllers.AuthController{}, "get:GetCurrentUser")
    web.Router("/api/user/stats", &controllers.AuthController{}, "get:GetUserStats")
    
    // 用户自定义分类接口（需要登录）
    web.Router("/api/user-categories", &controllers.UserCategoryController{}, "get:List;post:Create")
    web.Router("/api/user-categories/:id", &controllers.UserCategoryController{}, "put:Update;delete:Delete")
    
    // 获取当前用户的分类列表
    web.Router("/api/my-categories", &controllers.ArticleController{}, "get:GetUserCategories")
    
    // 文章管理接口（需要登录）
    web.Router("/api/my-articles", &controllers.ArticleController{}, "get:MyArticles")
    web.Router("/api/articles/create", &controllers.ArticleController{}, "post:CreateArticle")
    web.Router("/api/articles/:id", &controllers.ArticleController{}, "put:UpdateArticle;delete:DeleteArticle")
    web.Router("/api/upload/cover", &controllers.ArticleController{}, "post:UploadCover")
    
    // 评论接口（需要登录）
    web.Router("/api/article/:id/comments", &controllers.ArticleController{}, "post:CreateComment;get:GetComments")
    web.Router("/api/comment/:id", &controllers.ArticleController{}, "delete:DeleteComment")
    
    // 用户资料接口（公开访问）
    web.Router("/api/user/:id/profile", &controllers.ArticleController{}, "get:GetUserProfile")
    web.Router("/api/user/:id/stats", &controllers.ArticleController{}, "get:GetUserStats")
    web.Router("/api/user/:id/articles", &controllers.ArticleController{}, "get:GetUserArticles")
    web.Router("/api/user/:id/comments", &controllers.ArticleController{}, "get:GetUserComments")
    
    // 修改密码（需要登录）
    web.Router("/api/change-password", &controllers.AuthController{}, "post:ChangePassword")
    
    // 首页路由
    web.Router("/", &controllers.MainController{}, "*:Index")
    
    // 管理后台路由（需要 JWT 验证和管理员权限）
    web.InsertFilter("/admin/*", web.BeforeRouter, utils.JWTMiddleware)
    web.InsertFilter("/admin/*", web.BeforeRouter, utils.AdminOnlyMiddleware)
    
    web.Router("/admin", &controllers.AdminController{}, "get:Index")
    web.Router("/admin/stats", &controllers.AdminController{}, "get:Stats") // 统计数据 API
    
    // 文章管理
    web.Router("/admin/articles", &controllers.AdminController{}, "get:Articles")
    web.Router("/admin/article/create", &controllers.AdminController{}, "post:CreateArticle")
    web.Router("/admin/article/edit/:id", &controllers.AdminController{}, "get:EditArticle")
    web.Router("/admin/article/update/:id", &controllers.AdminController{}, "put:UpdateArticle")
    web.Router("/admin/article/delete/:id", &controllers.AdminController{}, "delete:DeleteArticle")
    
    // 分类管理
    web.Router("/admin/categories", &controllers.AdminController{}, "get:Categories")
    web.Router("/admin/category/create", &controllers.AdminController{}, "post:CreateCategory")
    web.Router("/admin/category/update/:id", &controllers.AdminController{}, "put:UpdateCategory")
    web.Router("/admin/category/delete/:id", &controllers.AdminController{}, "delete:DeleteCategory")
    
    // 标签管理
    web.Router("/admin/tags", &controllers.AdminController{}, "get:Tags")
    web.Router("/admin/tag/create", &controllers.AdminController{}, "post:CreateTag")
    web.Router("/admin/tag/update/:id", &controllers.AdminController{}, "put:UpdateTag")
    web.Router("/admin/tag/delete/:id", &controllers.AdminController{}, "delete:DeleteTag")
    
    // 用户管理
    web.Router("/admin/users", &controllers.AdminController{}, "get:Users")
    web.Router("/admin/user/delete/:id", &controllers.AdminController{}, "delete:DeleteUser")
}
