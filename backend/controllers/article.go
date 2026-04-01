package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"blog-system/models"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
)

// 文章控制器
type ArticleController struct {
	web.Controller
}

// 查看文章详情
func (c *ArticleController) View() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	o := orm.NewOrm()

	var article models.Article
	err := o.QueryTable(new(models.Article)).Filter("id", id).One(&article)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    404,
			"message": "文章不存在",
		}
		c.ServeJSON()
		return
	}

	// 增加浏览量
	article.ViewCount++
	o.Update(&article, "ViewCount")

	// 加载分类
	_, err = o.LoadRelated(&article, "Category")
	if err != nil {
		article.Category = nil
	}

	// 加载作者信息
	_, err = o.LoadRelated(&article, "Author")
	if err != nil {
		log.Printf("[ERROR] [Article.View] 加载作者信息失败 - article_id: %d, error: %v", article.Id, err)
		article.Author = nil
	}

	// 加载标签
	var tags []*models.Tag
	_, err = o.LoadRelated(&article, "Tags")
	if err == nil && article.Tags != nil {
		tags = article.Tags
	} else {
		tags = make([]*models.Tag, 0)
	}
	article.Tags = tags

	// 加载评论
	var comments []*models.Comment
	_, err = o.QueryTable(new(models.Comment)).Filter("article_id", id).Filter("status", 1).OrderBy("-created_at").All(&comments)
	if err != nil {
		// 如果评论表不存在或其他错误，返回空列表
		comments = make([]*models.Comment, 0)
	} else {
		// 加载每条评论的用户信息
		for _, comment := range comments {
			if comment != nil && comment.Id > 0 {
				_, err := o.LoadRelated(comment, "User")
				if err != nil {
					comment.User = nil
					continue
				}
			}
		}
	}

	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"article":  article,
			"comments": comments,
		},
	}
	c.ServeJSON()
}

// 文章列表
func (c *ArticleController) List() {
	page, _ := strconv.Atoi(c.GetString("page", "1"))
	pageSize, _ := strconv.Atoi(c.GetString("page_size", "10"))
	categorySlug := c.GetString("category", "")
	tagSlug := c.GetString("tag", "")
	search := c.GetString("search", "")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	o := orm.NewOrm()
	qs := o.QueryTable(new(models.Article)).Filter("status", 1)

	// 分类筛选
	if categorySlug != "" {
		var category models.Category
		err := o.QueryTable(new(models.Category)).Filter("slug", categorySlug).One(&category)
		if err == nil {
			qs = qs.Filter("category_id", category.Id)
		}
	}

	// 标签筛选
	if tagSlug != "" {
		var tag models.Tag
		err := o.QueryTable(new(models.Tag)).Filter("slug", tagSlug).One(&tag)
		if err == nil {
			qs = qs.Filter("tags__id", tag.Id)
		}
	}

	// 搜索
	if search != "" {
		qs = qs.Filter("title__icontains", search)
	}

	total, _ := qs.Count()

	var articles []*models.Article
	qs.OrderBy("-created_at").Limit(pageSize, (page-1)*pageSize).All(&articles)

	// 加载每篇文章的分类和标签
	for _, article := range articles {
		if article == nil || article.Id <= 0 {
			continue
		}
		_, err := o.LoadRelated(article, "Category")
		if err != nil {
			// 忽略加载分类错误，可能是没有关联分类
			article.Category = nil
		}
		_, err = o.LoadRelated(article, "Tags")
		if err != nil {
			// 忽略加载标签错误，可能是没有关联标签
			article.Tags = nil
		}
	}

	// 加载所有分类和标签
	var categories []*models.Category
	o.QueryTable(new(models.Category)).All(&categories)

	var tags []*models.Tag
	o.QueryTable(new(models.Tag)).All(&tags)

	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"articles":   articles,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"categories": categories,
			"tags":       tags,
		},
	}
	c.ServeJSON()
}

// 按分类查看
func (c *ArticleController) ByCategory() {
	slug := c.Ctx.Input.Param(":slug")
	o := orm.NewOrm()

	var category models.Category
	err := o.QueryTable(new(models.Category)).Filter("slug", slug).One(&category)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    404,
			"message": "分类不存在",
		}
		c.ServeJSON()
		return
	}

	var articles []*models.Article
	o.QueryTable(new(models.Article)).Filter("status", 1).Filter("category", &category).OrderBy("-created_at").All(&articles)

	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"category": category,
			"articles": articles,
		},
	}
	c.ServeJSON()
}

// 按标签查看
func (c *ArticleController) ByTag() {
	slug := c.Ctx.Input.Param(":slug")
	o := orm.NewOrm()

	var tag models.Tag
	err := o.QueryTable(new(models.Tag)).Filter("slug", slug).One(&tag)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    404,
			"message": "标签不存在",
		}
		c.ServeJSON()
		return
	}

	var articles []*models.Article
	o.QueryTable(new(models.Article)).Filter("status", 1).Filter("tags__id", tag.Id).Distinct().OrderBy("-created_at").All(&articles)

	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"tag":      tag,
			"articles": articles,
		},
	}
	c.ServeJSON()
}

// 获取用户自己的文章列表（用于创作中心）
func (c *ArticleController) MyArticles() {
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

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    401,
			"message": "用户 ID 格式错误",
		}
		c.ServeJSON()
		return
	}

	page, _ := strconv.Atoi(c.GetString("page", "1"))
	pageSize, _ := strconv.Atoi(c.GetString("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	o := orm.NewOrm()

	// 查询用户的文章（包括已发布和未发布的）
	qs := o.QueryTable(new(models.Article)).Filter("author_id", userID)

	total, _ := qs.Count()

	var articles []*models.Article
	qs.OrderBy("-created_at").Limit(pageSize, (page-1)*pageSize).All(&articles)

	// 加载每篇文章的分类和标签
	for _, article := range articles {
		if article == nil || article.Id <= 0 {
			continue
		}
		_, err := o.LoadRelated(article, "Category")
		if err != nil {
			article.Category = nil
		}
		_, err = o.LoadRelated(article, "Tags")
		if err != nil {
			article.Tags = nil
		}
	}

	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"articles":  articles,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	}
	c.ServeJSON()
}

// 创建文章
func (c *ArticleController) CreateArticle() {
	// 从 JWT 中间件中获取用户 ID
	userIDStr := c.Ctx.Input.Param("userID")
	if userIDStr == "" {
		c.Data["json"] = map[string]interface{}{
			"code":    401,
			"message": "未登录",
		}
		c.ServeJSON()
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "用户 ID 格式错误",
		}
		c.ServeJSON()
		return
	}

	// 解析请求体
	var req struct {
		Title            string `json:"title"`
		Slug             string `json:"slug"`
		Summary          string `json:"summary"`
		Content          string `json:"content"`
		CoverImage       string `json:"cover_image"`
		Status           int    `json:"status"`
		CategoryID       int    `json:"category_id"`        // 系统分类 ID
		UserCategoryID   int    `json:"user_category_id"`   // 用户自定义分类 ID
		TagIDs           []int  `json:"tag_ids"`
	}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
		}
		c.ServeJSON()
		return
	}

	// 验证必填字段
	if req.Title == "" {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "标题不能为空",
		}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()

	// 获取作者
	var author models.User
	err = o.QueryTable(new(models.User)).Filter("id", userID).One(&author)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    404,
			"message": "用户不存在",
		}
		c.ServeJSON()
		return
	}

	// 如果 slug 为空，根据标题自动生成
	slug := req.Slug
	if slug == "" {
		slug = generateUniqueSlug(req.Title, o)
	}
	
	// 创建文章
	article := &models.Article{
		Title:      req.Title,
		Slug:       slug,
		Summary:    req.Summary,
		Content:    req.Content,
		CoverImage: req.CoverImage,
		Status:     req.Status,
		Author:     &author,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// 设置分类（支持系统分类和用户自定义分类）
	if req.CategoryID > 0 {
		var category models.Category
		err := o.QueryTable(new(models.Category)).Filter("id", req.CategoryID).One(&category)
		if err == nil {
			article.Category = &category
			article.UserCategory = nil // 清除用户分类
		}
	} else if req.UserCategoryID > 0 {
		var userCategory models.UserCategory
		err := o.QueryTable(new(models.UserCategory)).Filter("id", req.UserCategoryID).Filter("user_id", userID).One(&userCategory)
		if err == nil {
			article.UserCategory = &userCategory
			article.Category = nil // 清除系统分类
		}
	} else {
		// 都没有设置，清空分类
		article.Category = nil
		article.UserCategory = nil
	}

	// 插入文章
	id, err := o.Insert(article)
	if err != nil {
		log.Printf("[ERROR] [Article.Create] 创建文章失败 - title: %s, author_id: %d, error: %v", article.Title, article.Author.Id, err)
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "创建文章失败：" + err.Error(),
		}
		c.ServeJSON()
		return
	}

	// 关联标签
	if len(req.TagIDs) > 0 {
		var tags []*models.Tag
		for _, tagID := range req.TagIDs {
			var tag models.Tag
			err := o.QueryTable(new(models.Tag)).Filter("id", tagID).One(&tag)
			if err == nil {
				tags = append(tags, &tag)
			}
		}
		if len(tags) > 0 {
			o.QueryM2M(&article, "Tags").Add(tags)
		}
	}

	c.Data["json"] = map[string]interface{}{
		"code":    0,
		"message": "创建成功",
		"data": map[string]interface{}{
			"id": id,
		},
	}
	c.ServeJSON()
}

// 更新文章
func (c *ArticleController) UpdateArticle() {
	// 从 JWT 中间件中获取用户 ID
	userIDStr := c.Ctx.Input.Param("userID")
	if userIDStr == "" {
		c.Data["json"] = map[string]interface{}{
			"code":    401,
			"message": "未登录",
		}
		c.ServeJSON()
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "用户 ID 格式错误",
		}
		c.ServeJSON()
		return
	}

	articleID, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if articleID <= 0 {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "文章 ID 错误",
		}
		c.ServeJSON()
		return
	}

	// 解析请求体
	var req struct {
		Title            string `json:"title"`
		Slug             string `json:"slug"`
		Summary          string `json:"summary"`
		Content          string `json:"content"`
		CoverImage       string `json:"cover_image"`
		Status           int    `json:"status"`
		CategoryID       int    `json:"category_id"`        // 系统分类 ID
		UserCategoryID   int    `json:"user_category_id"`   // 用户自定义分类 ID
		TagIDs           []int  `json:"tag_ids"`
	}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
		}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()

	// 查询文章并验证作者权限
	var article models.Article
	err = o.QueryTable(new(models.Article)).Filter("id", articleID).Filter("author_id", userID).One(&article)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    404,
			"message": "文章不存在或无权操作",
		}
		c.ServeJSON()
		return
	}

	// 如果 slug 为空，根据标题自动生成
	slug := req.Slug
	if slug == "" {
		slug = generateUniqueSlug(req.Title, o)
	}
	
	// 更新文章字段
	article.Title = req.Title
	article.Slug = slug
	article.Summary = req.Summary
	article.Content = req.Content
	article.CoverImage = req.CoverImage
	article.Status = req.Status
	article.UpdatedAt = time.Now()

	// 更新分类（支持系统分类和用户自定义分类）
	if req.CategoryID > 0 {
		var category models.Category
		err := o.QueryTable(new(models.Category)).Filter("id", req.CategoryID).One(&category)
		if err == nil {
			article.Category = &category
			article.UserCategory = nil // 清除用户分类
		}
	} else if req.UserCategoryID > 0 {
		var userCategory models.UserCategory
		err := o.QueryTable(new(models.UserCategory)).Filter("id", req.UserCategoryID).Filter("user_id", userID).One(&userCategory)
		if err == nil {
			article.UserCategory = &userCategory
			article.Category = nil // 清除系统分类
		}
	}
	// 注意：如果两个 ID 都为 0，保持原有分类不变，不强制清空

	_, err = o.Update(&article)
	if err != nil {
		log.Printf("[ERROR] [Article.Update] 更新文章失败 - id: %d, title: %s, error: %v", article.Id, article.Title, err)
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "更新文章失败：" + err.Error(),
		}
		c.ServeJSON()
		return
	}

	// 更新标签
	if req.TagIDs != nil {
		o.QueryM2M(&article, "Tags").Clear()
		if len(req.TagIDs) > 0 {
			var tags []*models.Tag
			for _, tagID := range req.TagIDs {
				var tag models.Tag
				err := o.QueryTable(new(models.Tag)).Filter("id", tagID).One(&tag)
				if err == nil {
					tags = append(tags, &tag)
				}
			}
			if len(tags) > 0 {
				o.QueryM2M(&article, "Tags").Add(tags)
			}
		}
	}

	c.Data["json"] = map[string]interface{}{
		"code":    0,
		"message": "更新成功",
	}
	c.ServeJSON()
}

// 删除文章
func (c *ArticleController) DeleteArticle() {
	// 从 JWT 中间件中获取用户 ID
	userIDStr := c.Ctx.Input.Param("userID")
	if userIDStr == "" {
		c.Data["json"] = map[string]interface{}{
			"code":    401,
			"message": "未登录",
		}
		c.ServeJSON()
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "用户 ID 格式错误",
		}
		c.ServeJSON()
		return
	}

	articleID, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if articleID <= 0 {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "文章 ID 错误",
		}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()

	// 查询文章并验证作者权限
	var article models.Article
	err = o.QueryTable(new(models.Article)).Filter("id", articleID).Filter("author_id", userID).One(&article)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    404,
			"message": "文章不存在或无权操作",
		}
		c.ServeJSON()
		return
	}

	// 删除文章（先删除多对多关系）
	o.QueryM2M(&article, "Tags").Clear()
	_, err = o.Delete(&article)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "删除文章失败",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"code":    0,
		"message": "删除成功",
	}
	c.ServeJSON()
}

// 上传封面图片
func (c *ArticleController) UploadCover() {
	// 从 JWT 中间件中获取用户 ID（可选，用于权限验证）
	userIDStr := c.Ctx.Input.Param("userID")
	if userIDStr != "" {
		_, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.Data["json"] = map[string]interface{}{
				"code":    401,
				"message": "未登录或 token 无效",
			}
			c.ServeJSON()
			return
		}
	}

	// 获取上传的文件
	file, header, err := c.GetFile("cover")
	if err != nil {
		log.Printf("[ERROR] [Article.UploadCover] 获取上传文件失败 - filename: %s, error: %v", header.Filename, err)
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "获取上传文件失败：" + err.Error(),
		}
		c.ServeJSON()
		return
	}
	defer file.Close()

	// 验证文件大小（最大 5MB）
	const maxFileSize = 5 << 20 // 5MB
	if header.Size > maxFileSize {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "图片大小不能超过 5MB",
		}
		c.ServeJSON()
		return
	}

	// 验证文件类型
	ext := filepath.Ext(header.Filename)
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	valid := false
	for _, allowed := range allowedExts {
		if strings.ToLower(ext) == allowed {
			valid = true
			break
		}
	}
	if !valid {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "只支持 jpg/png/gif/webp 格式的图片",
		}
		c.ServeJSON()
		return
	}

	// 创建 uploads/covers 目录
	uploadDir := "uploads/covers"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Printf("[ERROR] [Article.UploadCover] 创建目录失败 - path: %s, error: %v", uploadDir, err)
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "服务器错误：" + err.Error(),
		}
		c.ServeJSON()
		return
	}

	// 生成唯一文件名：timestamp_random.ext
	filename := fmt.Sprintf("%d_%d%s", time.Now().UnixNano(), time.Now().Unix(), ext)
	savePath := filepath.Join(uploadDir, filename)

	// 保存文件
	f, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Printf("[ERROR] [Article.UploadCover] 打开文件失败 - path: %s, error: %v", savePath, err)
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "保存文件失败：" + err.Error(),
		}
		c.ServeJSON()
		return
	}

	defer f.Close()
	if _, err := io.Copy(f, file); err != nil {
		log.Printf("[ERROR] [Article.UploadCover] 写入文件失败 - path: %s, size: %d, error: %v", savePath, header.Size, err)
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "保存文件失败：" + err.Error(),
		}
		c.ServeJSON()
		return
	}

	// 返回图片 URL（相对路径）
	imageURL := "/" + savePath

	c.Data["json"] = map[string]interface{}{
		"code":    0,
		"message": "上传成功",
		"data": map[string]string{
			"url": imageURL,
		},
	}
	c.ServeJSON()
}

// GetCategories 获取所有分类（公开接口，不需要 JWT）
func (c *ArticleController) GetCategories() {
	o := orm.NewOrm()
	
	var categories []*models.Category
	_, err := o.QueryTable(new(models.Category)).OrderBy("id").All(&categories)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "获取分类失败",
		}
		c.ServeJSON()
		return
	}
	
	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"categories": categories,
		},
	}
	c.ServeJSON()
}

// GetTags 获取所有标签（公开接口，不需要 JWT）
func (c *ArticleController) GetTags() {
	o := orm.NewOrm()
	
	var tags []*models.Tag
	_, err := o.QueryTable(new(models.Tag)).OrderBy("id").All(&tags)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "获取标签失败",
		}
		c.ServeJSON()
		return
	}
	
	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"tags": tags,
		},
	}
	c.ServeJSON()
}

// GetUserCategories 获取当前用户的分类列表（需要 JWT 认证）
func (c *ArticleController) GetUserCategories() {
	// 从 JWT 中间件中获取用户 ID
	userIDStr := c.Ctx.Input.Param("userID")
	if userIDStr == "" {
		c.Data["json"] = map[string]interface{}{
			"code":    401,
			"message": "未登录",
		}
		c.ServeJSON()
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "用户 ID 格式错误",
		}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	
	var categories []*models.UserCategory
	qs := o.QueryTable(new(models.UserCategory)).Filter("user_id", userID)
	num, err := qs.All(&categories)
	if err != nil {
		log.Printf("[ERROR] [Article.GetUserCategories] 获取用户分类失败 - userID: %d, error: %v", userID, err)
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "获取分类失败",
		}
		c.ServeJSON()
		return
	}
	log.Printf("[INFO] [Article.GetUserCategories] 用户 ID: %d, 分类数量：%d", userID, num)

	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"categories": categories,
		},
	}
	c.ServeJSON()
}

// 点赞文章
func (c *ArticleController) Like() {
	// 从 JWT 中间件中获取用户 ID
	userIDStr := c.Ctx.Input.Param("userID")
	if userIDStr == "" {
		c.Data["json"] = map[string]interface{}{
			"code":    401,
			"message": "未登录",
		}
		c.ServeJSON()
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "用户 ID 格式错误",
		}
		c.ServeJSON()
		return
	}

	articleID, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if articleID <= 0 {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "文章 ID 错误",
		}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()

	// 检查文章是否存在
	var article models.Article
	err = o.QueryTable(new(models.Article)).Filter("id", articleID).One(&article)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    404,
			"message": "文章不存在",
		}
		c.ServeJSON()
		return
	}

	// 检查用户是否存在
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

	// 检查是否已经点过赞
	var like models.Like
	err = o.QueryTable(new(models.Like)).Filter("article_id", articleID).Filter("user_id", userID).One(&like)
	if err == nil {
		// 已经点过赞，取消点赞
		o.Delete(&like)
		c.Data["json"] = map[string]interface{}{
			"code":    0,
			"message": "已取消点赞",
			"data": map[string]interface{}{
				"liked": false,
				"count": getLikeCount(articleID, o),
			},
		}
	} else {
		// 未点过赞，添加点赞
		like = models.Like{
			Article:   &article,
			User:      &user,
			CreatedAt: time.Now(),
		}
		_, err = o.Insert(&like)
		if err != nil {
			c.Data["json"] = map[string]interface{}{
				"code":    500,
				"message": "点赞失败",
			}
			c.ServeJSON()
			return
		}
		c.Data["json"] = map[string]interface{}{
			"code":    0,
			"message": "点赞成功",
			"data": map[string]interface{}{
				"liked": true,
				"count": getLikeCount(articleID, o),
			},
		}
	}
	c.ServeJSON()
}

// 获取点赞状态和数量
func (c *ArticleController) GetLikeStatus() {
	articleID, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if articleID <= 0 {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "文章 ID 错误",
		}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	userID := 0
	userIDStr := c.Ctx.Input.Param("userID")
	if userIDStr != "" {
		userID, _ = strconv.Atoi(userIDStr)
	}

	liked := false
	if userID > 0 {
		var like models.Like
		err := o.QueryTable(new(models.Like)).Filter("article_id", articleID).Filter("user_id", userID).One(&like)
		if err == nil {
			liked = true
		}
	}

	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"liked": liked,
			"count": getLikeCount(articleID, o),
		},
	}
	c.ServeJSON()
}

// 获取点赞数量
func getLikeCount(articleID int, o orm.Ormer) int {
	count, _ := o.QueryTable(new(models.Like)).Filter("article_id", articleID).Count()
	return int(count)
}

// 关注用户
func (c *ArticleController) Follow() {
	// 从 JWT 中间件中获取用户 ID（关注者）
	followerIDStr := c.Ctx.Input.Param("userID")
	if followerIDStr == "" {
		c.Data["json"] = map[string]interface{}{
			"code":    401,
			"message": "未登录",
		}
		c.ServeJSON()
		return
	}

	followerID, err := strconv.Atoi(followerIDStr)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "用户 ID 格式错误",
		}
		c.ServeJSON()
		return
	}

	followeeID, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if followeeID <= 0 {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "被关注用户 ID 错误",
		}
		c.ServeJSON()
		return
	}

	// 不能关注自己
	if followerID == followeeID {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "不能关注自己",
		}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()

	// 检查关注者是否存在
	var follower models.User
	err = o.QueryTable(new(models.User)).Filter("id", followerID).One(&follower)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    404,
			"message": "用户不存在",
		}
		c.ServeJSON()
		return
	}

	// 检查被关注者是否存在
	var followee models.User
	err = o.QueryTable(new(models.User)).Filter("id", followeeID).One(&followee)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    404,
			"message": "被关注用户不存在",
		}
		c.ServeJSON()
		return
	}

	// 检查是否已经关注
	var follow models.Follow
	err = o.QueryTable(new(models.Follow)).Filter("follower_id", followerID).Filter("followee_id", followeeID).One(&follow)
	if err == nil {
		// 已经关注，取消关注
		o.Delete(&follow)
		c.Data["json"] = map[string]interface{}{
			"code":    0,
			"message": "已取消关注",
			"data": map[string]interface{}{
				"following": false,
				"count":     getFollowerCount(followeeID, o),
			},
		}
	} else {
		// 未关注，添加关注
		follow = models.Follow{
			Follower:  &follower,
			Followee:  &followee,
			CreatedAt: time.Now(),
		}
		_, err = o.Insert(&follow)
		if err != nil {
			c.Data["json"] = map[string]interface{}{
				"code":    500,
				"message": "关注失败",
			}
			c.ServeJSON()
			return
		}
		c.Data["json"] = map[string]interface{}{
			"code":    0,
			"message": "关注成功",
			"data": map[string]interface{}{
				"following": true,
				"count":     getFollowerCount(followeeID, o),
			},
		}
	}
	c.ServeJSON()
}

// 获取关注状态和粉丝数量
func (c *ArticleController) GetFollowStatus() {
	followeeID, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if followeeID <= 0 {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "用户 ID 错误",
		}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	followerID := 0
	followerIDStr := c.Ctx.Input.Param("userID")
	if followerIDStr != "" {
		followerID, _ = strconv.Atoi(followerIDStr)
	}

	following := false
	if followerID > 0 && followerID != followeeID {
		var follow models.Follow
		err := o.QueryTable(new(models.Follow)).Filter("follower_id", followerID).Filter("followee_id", followeeID).One(&follow)
		if err == nil {
			following = true
		}
	}

	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"following": following,
			"count":     getFollowerCount(followeeID, o),
		},
	}
	c.ServeJSON()
}

// 获取粉丝数量
func getFollowerCount(userID int, o orm.Ormer) int {
	count, _ := o.QueryTable(new(models.Follow)).Filter("followee_id", userID).Count()
	return int(count)
}

// 发表评论
func (c *ArticleController) CreateComment() {
	// 从 JWT 中间件中获取用户 ID
	userIDStr := c.Ctx.Input.Param("userID")
	if userIDStr == "" {
		c.Data["json"] = map[string]interface{}{
			"code":    401,
			"message": "未登录",
		}
		c.ServeJSON()
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "用户 ID 格式错误",
		}
		c.ServeJSON()
		return
	}

	articleID, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if articleID <= 0 {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "文章 ID 错误",
		}
		c.ServeJSON()
		return
	}

	// 解析请求体
	var req struct {
		Content string `json:"content"`
	}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
		}
		c.ServeJSON()
		return
	}

	if req.Content == "" {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "评论内容不能为空",
		}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()

	// 检查文章是否存在
	var article models.Article
	err = o.QueryTable(new(models.Article)).Filter("id", articleID).One(&article)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    404,
			"message": "文章不存在",
		}
		c.ServeJSON()
		return
	}

	// 获取用户信息
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

	// 创建评论
	comment := models.Comment{
		Article:   &article,
		User:      &user,
		Author:    user.Nickname,
		Email:     user.Email,
		Content:   req.Content,
		Status:    1, // 默认审核通过
		CreatedAt: time.Now(),
	}

	_, err = o.Insert(&comment)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "评论失败",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"code":    0,
		"message": "评论成功",
		"data": map[string]interface{}{
			"comment": comment,
		},
	}
	c.ServeJSON()
}

// 删除评论
func (c *ArticleController) DeleteComment() {
	// 从 JWT 中间件中获取用户 ID
	userIDStr := c.Ctx.Input.Param("userID")
	if userIDStr == "" {
		c.Data["json"] = map[string]interface{}{
			"code":    401,
			"message": "未登录",
		}
		c.ServeJSON()
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "用户 ID 格式错误",
		}
		c.ServeJSON()
		return
	}

	commentID, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if commentID <= 0 {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "评论 ID 错误",
		}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()

	// 查询评论并验证权限（只有评论作者或管理员可以删除）
	var comment models.Comment
	err = o.QueryTable(new(models.Comment)).Filter("id", commentID).One(&comment)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    404,
			"message": "评论不存在",
		}
		c.ServeJSON()
		return
	}

	// 检查是否为评论作者或管理员
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

	if comment.User != nil && comment.User.Id != userID && user.Role != 0 {
		c.Data["json"] = map[string]interface{}{
			"code":    403,
			"message": "无权删除该评论",
		}
		c.ServeJSON()
		return
	}

	_, err = o.Delete(&comment)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "删除评论失败",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"code":    0,
		"message": "删除成功",
	}
	c.ServeJSON()
}

// generateUniqueSlug 根据标题生成唯一的 slug
func generateUniqueSlug(title string, o orm.Ormer) string {
	// 将中文标题转换为拼音（简单实现：移除特殊字符，保留字母数字）
	slug := strings.ToLower(title)
	
	// 移除特殊字符，保留字母、数字、空白字符和中文字符
	// 使用双引号字符串以支持 Unicode 转义序列
	reg := regexp.MustCompile(`[^\w\s\p{Han}]`)
	slug = reg.ReplaceAllString(slug, "")
	
	// 将空格和下划线替换为连字符
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	
	// 限制长度
	if len(slug) > 100 {
		slug = slug[:100]
	}
	
	// 如果 slug 为空，使用时间戳
	if slug == "" {
		slug = fmt.Sprintf("article-%d", time.Now().Unix())
	}
	
	// 检查 slug 是否已存在，如果存在则添加随机后缀
	originalSlug := slug
	counter := 1
	for {
		var article models.Article
		err := o.QueryTable(new(models.Article)).Filter("slug", slug).One(&article)
		if err == orm.ErrNoRows {
			// slug 不存在，可以使用
			break
		} else if err != nil {
			// 其他错误，也直接使用
			break
		}
		// slug 已存在，添加后缀
		slug = fmt.Sprintf("%s-%d", originalSlug, counter)
		counter++
	}
	
	return slug
}

// GetComments 获取文章评论列表（公开接口）
func (c *ArticleController) GetComments() {
    articleID, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
    if articleID <= 0 {
        c.Data["json"] = map[string]interface{}{
            "code":    400,
            "message": "文章 ID 错误",
        }
        c.ServeJSON()
        return
    }

    o := orm.NewOrm()
    var comments []*models.Comment
    _, err := o.QueryTable(new(models.Comment)).Filter("article_id", articleID).Filter("status", 1).OrderBy("-created_at").All(&comments)
    
    if err != nil {
        c.Data["json"] = map[string]interface{}{
            "code":    500,
            "message": "获取评论失败",
        }
        c.ServeJSON()
        return
    }

    // 加载每条评论的用户信息（根据规范必须显式加载关联数据）
    for _, comment := range comments {
        if comment != nil && comment.Id > 0 {
            // 使用 LoadRelated 加载 User 关联
            _, err := o.LoadRelated(comment, "User")
            if err != nil {
                // 加载失败也不影响评论显示，只是用户信息为空
                continue
            }
        }
    }

    c.Data["json"] = map[string]interface{}{
        "code": 0,
        "data": map[string]interface{}{
            "comments": comments,
        },
    }
    c.ServeJSON()
}

// 获取用户资料
func (c *ArticleController) GetUserProfile() {
	userID, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if userID <= 0 {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "用户 ID 错误",
		}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("id", userID).One(&user)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    404,
			"message": "用户不存在",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"data": user,
	}
	c.ServeJSON()
}

// 获取用户统计信息
func (c *ArticleController) GetUserStats() {
	userID, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if userID <= 0 {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "用户 ID 错误",
		}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	
	// 统计文章数
	articleCount, _ := o.QueryTable(new(models.Article)).Filter("author_id", userID).Filter("status", 1).Count()
	
	// 统计评论数
	commentCount, _ := o.QueryTable(new(models.Comment)).Filter("user_id", userID).Filter("status", 1).Count()
	
	// 统计获赞数
	likeCount, _ := o.QueryTable(new(models.Like)).Filter("article__author_id", userID).Count()

	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"articleCount": articleCount,
			"commentCount": commentCount,
			"likeCount":    likeCount,
		},
	}
	c.ServeJSON()
}

// 获取用户文章列表
func (c *ArticleController) GetUserArticles() {
	userID, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if userID <= 0 {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "用户 ID 错误",
		}
		c.ServeJSON()
		return
	}

	page, _ := strconv.Atoi(c.Ctx.Input.Param("page"))
	pageSize, _ := strconv.Atoi(c.Ctx.Input.Param("pageSize"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	o := orm.NewOrm()
	var articles []*models.Article
	total, _ := o.QueryTable(new(models.Article)).Filter("author_id", userID).Filter("status", 1).Count()
	_, err := o.QueryTable(new(models.Article)).Filter("author_id", userID).Filter("status", 1).
		OrderBy("-created_at").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		All(&articles)

	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "获取文章失败",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"articles": articles,
			"total":    total,
		},
	}
	c.ServeJSON()
}

// 获取用户评论列表
func (c *ArticleController) GetUserComments() {
	userID, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if userID <= 0 {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "用户 ID 错误",
		}
		c.ServeJSON()
		return
	}

	page, _ := strconv.Atoi(c.Ctx.Input.Param("page"))
	pageSize, _ := strconv.Atoi(c.Ctx.Input.Param("pageSize"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	o := orm.NewOrm()
	var comments []*models.Comment
	total, _ := o.QueryTable(new(models.Comment)).Filter("user_id", userID).Filter("status", 1).Count()
	_, err := o.QueryTable(new(models.Comment)).Filter("user_id", userID).Filter("status", 1).
		OrderBy("-created_at").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		All(&comments)

	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "获取评论失败",
		}
		c.ServeJSON()
		return
	}

	// 加载评论关联的文章信息
	for _, comment := range comments {
		if comment != nil && comment.Id > 0 {
			_, err := o.LoadRelated(comment, "Article")
			if err != nil {
				continue
			}
		}
	}

	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"comments": comments,
			"total":    total,
		},
	}
	c.ServeJSON()
}
