package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	"blog-system/models"
	"log"
	"strconv"
	"time"
)

// 管理后台控制器
type AdminController struct {
	web.Controller
}

// 管理后台首页 - 数据概览
func (c *AdminController) Index() {
	o := orm.NewOrm()
	
	// 统计数据
	var articleCount int64
	var userCount int64
	var categoryCount int64
	var tagCount int64
	var commentCount int64
	
	qsArticle := o.QueryTable(new(models.Article))
	articleCount, _ = qsArticle.Count()
	
	qsUser := o.QueryTable(new(models.User))
	userCount, _ = qsUser.Count()
	
	qsCategory := o.QueryTable(new(models.Category))
	categoryCount, _ = qsCategory.Count()
	
	qsTag := o.QueryTable(new(models.Tag))
	tagCount, _ = qsTag.Count()
	
	qsComment := o.QueryTable(new(models.Comment))
	commentCount, _ = qsComment.Count()
	
	// 获取最新文章
	var recentArticles []*models.Article
	o.QueryTable(new(models.Article)).OrderBy("-created_at").Limit(5).All(&recentArticles)
	
	// 加载关联数据
	for _, article := range recentArticles {
		if article.Author != nil {
			o.LoadRelated(article, "Author")
		}
		if article.Category != nil {
			o.LoadRelated(article, "Category")
		}
		o.LoadRelated(article, "Tags")
	}
	
	// 获取热门用户 (按文章数)
	var users []*models.User
	o.QueryTable(new(models.User)).OrderBy("-id").Limit(5).All(&users)
	
	c.Data["ArticleCount"] = articleCount
	c.Data["UserCount"] = userCount
	c.Data["CategoryCount"] = categoryCount
	c.Data["TagCount"] = tagCount
	c.Data["CommentCount"] = commentCount
	c.Data["RecentArticles"] = recentArticles
	c.Data["ActiveUsers"] = users
	c.Data["PageTitle"] = "管理后台"
	c.TplName = "admin/index.tpl"
}

// 管理后台统计数据 API
func (c *AdminController) Stats() {
	o := orm.NewOrm()
	
	// 统计数据
	var articleCount int64
	var userCount int64
	var categoryCount int64
	var tagCount int64
	var commentCount int64
	
	qsArticle := o.QueryTable(new(models.Article))
	articleCount, _ = qsArticle.Count()
	
	qsUser := o.QueryTable(new(models.User))
	userCount, _ = qsUser.Count()
	
	qsCategory := o.QueryTable(new(models.Category))
	categoryCount, _ = qsCategory.Count()
	
	qsTag := o.QueryTable(new(models.Tag))
	tagCount, _ = qsTag.Count()
	
	qsComment := o.QueryTable(new(models.Comment))
	commentCount, _ = qsComment.Count()
	
	// 获取最新文章
	var recentArticles []*models.Article
	o.QueryTable(new(models.Article)).OrderBy("-created_at").Limit(5).All(&recentArticles)
	
	// 获取活跃用户
	var users []*models.User
	o.QueryTable(new(models.User)).OrderBy("-id").Limit(5).All(&users)
	
	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"stats": map[string]int64{
				"articleCount":  articleCount,
				"userCount":     userCount,
				"categoryCount": categoryCount,
				"tagCount":      tagCount,
				"commentCount":  commentCount,
			},
			"recentArticles": recentArticles,
			"activeUsers":    users,
		},
	}
	c.ServeJSON()
}

// 创建文章
func (c *AdminController) CreateArticle() {
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
		Title      string `json:"title"`
		Slug       string `json:"slug"`
		Summary    string `json:"summary"`
		Content    string `json:"content"`
		CoverImage string `json:"cover_image"`
		Status     int    `json:"status"`
		CategoryID int    `json:"category_id"`
		TagIDs     []int  `json:"tag_ids"`
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
		UpdatedAt:  time.Now(),
	}

	// 设置作者
	var author models.User
	err = o.QueryTable(new(models.User)).Filter("id", userID).One(&author)
	if err == nil {
		article.Author = &author
	}

	// 设置分类
	if req.CategoryID > 0 {
		var category models.Category
		err := o.QueryTable(new(models.Category)).Filter("id", req.CategoryID).One(&category)
		if err == nil {
			article.Category = &category
		}
	}

	// 插入文章
	_, err = o.Insert(article)
	if err != nil {
		log.Printf("[ERROR] [Admin.CreateArticle] 创建文章失败 - title: %s, author_id: %d, error: %v", article.Title, article.Author.Id, err)
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "创建文章失败：" + err.Error(),
		}
		c.ServeJSON()
		return
	}

	// 关联标签
	if len(req.TagIDs) > 0 {
		m2m := o.QueryM2M(article, "Tags")
		for _, tagID := range req.TagIDs {
			var tag models.Tag
			err := o.QueryTable(new(models.Tag)).Filter("id", tagID).One(&tag)
			if err == nil {
				m2m.Add(&tag)
			}
		}
	}

	c.Data["json"] = map[string]interface{}{
		"code":    0,
		"message": "创建成功",
		"data":    article,
	}
	c.ServeJSON()
}

// 文章管理列表
func (c *AdminController) Articles() {
	o := orm.NewOrm()
	
	// 获取所有文章
	var articles []*models.Article
	qs := o.QueryTable(new(models.Article))
	
	// 搜索过滤
	title := c.GetString("title")
	if title != "" {
		qs = qs.Filter("title__icontains", title)
	}
	
	// 状态过滤
	status := c.GetString("status")
	if status != "" {
		statusInt, _ := strconv.Atoi(status)
		qs = qs.Filter("status", statusInt)
	}
	
	// 排序
	order := c.GetString("order", "-created_at")
	qs = qs.OrderBy(order)
	
	// 分页
	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("page_size", 20)
	
	total, _ := qs.Count()
	qs.Limit(pageSize).Offset((page - 1) * pageSize).All(&articles)
	
	// 加载关联数据 (作者、分类、标签)
	for _, article := range articles {
		// 加载作者信息
		if article.Author != nil {
			o.LoadRelated(article, "Author")
		}
		
		// 加载分类信息
		if article.Category != nil {
			o.LoadRelated(article, "Category")
		}
		
		// 加载标签信息
		o.LoadRelated(article, "Tags")
	}
	
	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"articles": articles,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	}
	c.ServeJSON()
}

// 编辑文章页面
func (c *AdminController) EditArticle() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	o := orm.NewOrm()
	
	var article models.Article
	err := o.QueryTable(new(models.Article)).Filter("id", id).One(&article)
	if err != nil {
		c.Ctx.WriteString("文章不存在")
		return
	}
	
	// 加载关联数据
	o.LoadRelated(&article, "Author")
	o.LoadRelated(&article, "Category")
	o.LoadRelated(&article, "Tags")
	
	var categories []*models.Category
	o.QueryTable(new(models.Category)).All(&categories)
	
	var tags []*models.Tag
	o.QueryTable(new(models.Tag)).All(&tags)
	
	c.Data["Article"] = article
	c.Data["Categories"] = categories
	c.Data["Tags"] = tags
	c.Data["PageTitle"] = "编辑文章"
	c.TplName = "admin/edit_article.tpl"
}

// 更新文章
func (c *AdminController) UpdateArticle() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	
	// 解析请求体
	var req struct {
		Title      string `json:"title"`
		Slug       string `json:"slug"`
		Summary    string `json:"summary"`
		Content    string `json:"content"`
		CoverImage string `json:"cover_image"`
		Status     int    `json:"status"`
		CategoryID int    `json:"category_id"`
		TagIDs     []int  `json:"tag_ids"`
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
	
	// 更新分类
	if req.CategoryID > 0 {
		var category models.Category
		err := o.QueryTable(new(models.Category)).Filter("id", req.CategoryID).One(&category)
		if err == nil {
			article.Category = &category
		}
	}
	
	// 更新文章
	_, err = o.Update(&article)
	if err != nil {
		log.Printf("[ERROR] [Admin.UpdateArticle] 更新文章失败 - id: %d, title: %s, error: %v", article.Id, article.Title, err)
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "更新文章失败：" + err.Error(),
		}
		c.ServeJSON()
		return
	}
	
	// 更新标签关联
	if len(req.TagIDs) > 0 {
		m2m := o.QueryM2M(&article, "Tags")
		m2m.Clear()
		
		for _, tagID := range req.TagIDs {
			var tag models.Tag
			err := o.QueryTable(new(models.Tag)).Filter("id", tagID).One(&tag)
			if err == nil {
				m2m.Add(&tag)
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
func (c *AdminController) DeleteArticle() {
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
	
	// 删除关联的标签
	m2m := o.QueryM2M(&article, "Tags")
	m2m.Clear()
	
	// 删除文章
	_, err = o.Delete(&article)
	if err != nil {
		log.Printf("[ERROR] [Admin.DeleteArticle] 删除文章失败 - id: %d, title: %s, error: %v", article.Id, article.Title, err)
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "删除文章失败：" + err.Error(),
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

// 分类管理列表
func (c *AdminController) Categories() {
	o := orm.NewOrm()
	
	var categories []*models.Category
	o.QueryTable(new(models.Category)).All(&categories)
	
	// 统计每个分类的文章数
	for _, cat := range categories {
		count, _ := o.QueryTable(new(models.Article)).Filter("category", cat).Count()
		cat.ArticleCount = int(count)
	}
	
	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"categories": categories,
		},
	}
	c.ServeJSON()
}

// 创建分类
func (c *AdminController) CreateCategory() {
	var req struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
	
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
		}
		c.ServeJSON()
		return
	}
	
	if req.Name == "" || req.Slug == "" {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "分类名称和标识不能为空",
		}
		c.ServeJSON()
		return
	}
	
	o := orm.NewOrm()
	
	// 检查 slug 是否已存在
	var existCategory models.Category
	err := o.QueryTable(new(models.Category)).Filter("slug", req.Slug).One(&existCategory)
	if err == nil {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "分类标识已存在",
		}
		c.ServeJSON()
		return
	}
	
	category := &models.Category{
		Name:      req.Name,
		Slug:      req.Slug,
	}
	
	_, err = o.Insert(category)
	if err != nil {
		log.Printf("[ERROR] [Admin.CreateCategory] 创建分类失败 - name: %s, slug: %s, error: %v", category.Name, category.Slug, err)
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "创建分类失败：" + err.Error(),
		}
		c.ServeJSON()
		return
	}
	
	c.Data["json"] = map[string]interface{}{
		"code":    0,
		"message": "创建成功",
		"data":    category,
	}
	c.ServeJSON()
}

// 更新分类
func (c *AdminController) UpdateCategory() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	
	var req struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
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
	
	var category models.Category
	err := o.QueryTable(new(models.Category)).Filter("id", id).One(&category)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    404,
			"message": "分类不存在",
		}
		c.ServeJSON()
		return
	}
	
	// 检查 slug 是否被其他分类使用
	var existCategory models.Category
	err = o.QueryTable(new(models.Category)).Filter("slug", req.Slug).Exclude("id", id).One(&existCategory)
	if err == nil {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "分类标识已存在",
		}
		c.ServeJSON()
		return
	}
	
	category.Name = req.Name
	category.Slug = req.Slug
	
	_, err = o.Update(&category)
	if err != nil {
		log.Printf("[ERROR] [Admin.UpdateCategory] 更新分类失败 - id: %d, name: %s, error: %v", category.Id, category.Name, err)
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "更新分类失败：" + err.Error(),
		}
		c.ServeJSON()
		return
	}
	
	c.Data["json"] = map[string]interface{}{
		"code":    0,
		"message": "更新成功",
	}
	c.ServeJSON()
}

// 删除分类
func (c *AdminController) DeleteCategory() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	o := orm.NewOrm()
	
	var category models.Category
	err := o.QueryTable(new(models.Category)).Filter("id", id).One(&category)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    404,
			"message": "分类不存在",
		}
		c.ServeJSON()
		return
	}
	
	// 检查是否有文章使用该分类
	articleCount, _ := o.QueryTable(new(models.Article)).Filter("category", &category).Count()
	if articleCount > 0 {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": fmt.Sprintf("该分类下有 %d 篇文章，无法删除", articleCount),
		}
		c.ServeJSON()
		return
	}
	
	_, err = o.Delete(&category)
	if err != nil {
		log.Printf("[ERROR] [Admin.DeleteCategory] 删除分类失败 - id: %d, name: %s, error: %v", category.Id, category.Name, err)
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "删除分类失败：" + err.Error(),
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

// 标签管理列表
func (c *AdminController) Tags() {
	o := orm.NewOrm()
	
	var tags []*models.Tag
	o.QueryTable(new(models.Tag)).All(&tags)
	
	// 统计每个标签的文章数
	for _, tag := range tags {
		count, _ := o.QueryTable(new(models.Article)).Filter("tags", tag).Count()
		tag.ArticleCount = int(count)
	}
	
	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"tags": tags,
		},
	}
	c.ServeJSON()
}

// 创建标签
func (c *AdminController) CreateTag() {
	var req struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
	
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
		}
		c.ServeJSON()
		return
	}
	
	if req.Name == "" || req.Slug == "" {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "标签名称和标识不能为空",
		}
		c.ServeJSON()
		return
	}
	
	o := orm.NewOrm()
	
	// 检查 slug 是否已存在
	var existTag models.Tag
	err := o.QueryTable(new(models.Tag)).Filter("slug", req.Slug).One(&existTag)
	if err == nil {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "标签标识已存在",
		}
		c.ServeJSON()
		return
	}
	
	tag := &models.Tag{
		Name:      req.Name,
		Slug:      req.Slug,
	}
	
	_, err = o.Insert(tag)
	if err != nil {
		log.Printf("[ERROR] [Admin.CreateTag] 创建标签失败 - name: %s, slug: %s, error: %v", tag.Name, tag.Slug, err)
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "创建标签失败：" + err.Error(),
		}
		c.ServeJSON()
		return
	}
	
	c.Data["json"] = map[string]interface{}{
		"code":    0,
		"message": "创建成功",
		"data":    tag,
	}
	c.ServeJSON()
}

// 更新标签
func (c *AdminController) UpdateTag() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	
	var req struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
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
	
	var tag models.Tag
	err := o.QueryTable(new(models.Tag)).Filter("id", id).One(&tag)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    404,
			"message": "标签不存在",
		}
		c.ServeJSON()
		return
	}
	
	// 检查 slug 是否被其他标签使用
	var existTag models.Tag
	err = o.QueryTable(new(models.Tag)).Filter("slug", req.Slug).Exclude("id", id).One(&existTag)
	if err == nil {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "标签标识已存在",
		}
		c.ServeJSON()
		return
	}
	
	tag.Name = req.Name
	tag.Slug = req.Slug
	
	_, err = o.Update(&tag)
	if err != nil {
		log.Printf("[ERROR] [Admin.UpdateTag] 更新标签失败 - id: %d, name: %s, error: %v", tag.Id, tag.Name, err)
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "更新标签失败：" + err.Error(),
		}
		c.ServeJSON()
		return
	}
	
	c.Data["json"] = map[string]interface{}{
		"code":    0,
		"message": "更新成功",
	}
	c.ServeJSON()
}

// 删除标签
func (c *AdminController) DeleteTag() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	o := orm.NewOrm()
	
	var tag models.Tag
	err := o.QueryTable(new(models.Tag)).Filter("id", id).One(&tag)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    404,
			"message": "标签不存在",
		}
		c.ServeJSON()
		return
	}
	
	// 检查是否有文章使用该标签
	articleCount, _ := o.QueryTable(new(models.Article)).Filter("tags", &tag).Count()
	if articleCount > 0 {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": fmt.Sprintf("该标签被 %d 篇文章使用，无法删除", articleCount),
		}
		c.ServeJSON()
		return
	}
	
	_, err = o.Delete(&tag)
	if err != nil {
		log.Printf("[ERROR] [Admin.DeleteTag] 删除标签失败 - id: %d, name: %s, error: %v", tag.Id, tag.Name, err)
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "删除标签失败：" + err.Error(),
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

// 用户管理列表
func (c *AdminController) Users() {
	o := orm.NewOrm()
	
	var users []*models.User
	qs := o.QueryTable(new(models.User))
	
	// 搜索过滤 - 注意：必须从 URL 查询参数获取，而不是从 Input Param 获取
	// 因为中间件设置了 username 参数，会干扰搜索逻辑
	username := c.Ctx.Request.URL.Query().Get("username")
	if username != "" {
		qs = qs.Filter("username__icontains", username)
	}
	
	email := c.Ctx.Request.URL.Query().Get("email")
	if email != "" {
		qs = qs.Filter("email__icontains", email)
	}
	
	// 分页
	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("page_size", 20)
	
	// 确保 page 至少为 1
	if page < 1 {
		page = 1
	}
	
	total, _ := qs.Count()
	
	_, err := qs.Limit(pageSize).Offset((page - 1) * pageSize).All(&users)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "查询失败：" + err.Error(),
		}
		c.ServeJSON()
		return
	}
	
	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"users":    users,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	}
	c.ServeJSON()
}

// 删除用户
func (c *AdminController) DeleteUser() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	o := orm.NewOrm()
	
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("id", id).One(&user)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    404,
			"message": "用户不存在",
		}
		c.ServeJSON()
		return
	}
	
	// 检查是否有文章
	articleCount, _ := o.QueryTable(new(models.Article)).Filter("author", &user).Count()
	if articleCount > 0 {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": fmt.Sprintf("该用户发布了 %d 篇文章，无法删除", articleCount),
		}
		c.ServeJSON()
		return
	}
	
	_, err = o.Delete(&user)
	if err != nil {
		log.Printf("[ERROR] [Admin.DeleteUser] 删除用户失败 - id: %d, username: %s, email: %s, error: %v", user.Id, user.Username, user.Email, err)
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "删除用户失败：" + err.Error(),
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
