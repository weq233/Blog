package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"blog-system/models"
	"strconv"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
)

// 用户自定义分类控制器
type UserCategoryController struct {
	web.Controller
}

// 获取用户的分类列表（需要 JWT 认证）
func (c *UserCategoryController) List() {
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
	_, err = o.QueryTable(new(models.UserCategory)).Filter("user_id", userID).OrderBy("-created_at").All(&categories)
	if err != nil {
		log.Printf("[ERROR] [UserCategory.List] 获取用户分类失败 - userID: %d, error: %v", userID, err)
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "获取分类列表失败：" + err.Error(),
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

// 创建用户分类（需要 JWT 认证）
func (c *UserCategoryController) Create() {
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
		Name  string `json:"name"`
		Color string `json:"color"`
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
	if req.Name == "" {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "分类名称不能为空",
		}
		c.ServeJSON()
		return
	}

	// 检查是否已存在同名分类
	o := orm.NewOrm()
	var existing models.UserCategory
	err = o.QueryTable(new(models.UserCategory)).Filter("user_id", userID).Filter("name", req.Name).One(&existing)
	if err == nil {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "该名称的分类已存在",
		}
		c.ServeJSON()
		return
	}

	// 设置默认颜色
	if req.Color == "" {
		req.Color = "#409EFF"
	}

	// 创建分类
	category := &models.UserCategory{
		UserID:    userID,
		Name:      req.Name,
		Color:     req.Color,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := o.Insert(category)
	if err != nil {
		log.Printf("[ERROR] [UserCategory.Create] 创建用户分类失败 - userID: %d, name: %s, color: %s, error: %v", category.UserID, category.Name, category.Color, err)
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "创建分类失败",
		}
		c.ServeJSON()
		return
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

// 更新用户分类（需要 JWT 认证）
func (c *UserCategoryController) Update() {
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

	categoryID, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if categoryID <= 0 {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "分类 ID 错误",
		}
		c.ServeJSON()
		return
	}

	// 解析请求体
	var req struct {
		Name  string `json:"name"`
		Color string `json:"color"`
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
	
	// 查询分类并验证所有权
	var category models.UserCategory
	err = o.QueryTable(new(models.UserCategory)).Filter("id", categoryID).Filter("user_id", userID).One(&category)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    404,
			"message": "分类不存在或无权操作",
		}
		c.ServeJSON()
		return
	}

	// 验证名称
	if req.Name != "" {
		// 检查是否有其他分类使用了该名称
		var existing models.UserCategory
		err = o.QueryTable(new(models.UserCategory)).
			Filter("user_id", userID).
			Filter("name", req.Name).
			Exclude("id", categoryID).
			One(&existing)
		if err == nil {
			c.Data["json"] = map[string]interface{}{
				"code":    400,
				"message": "该名称的分类已存在",
			}
			c.ServeJSON()
			return
		}
		category.Name = req.Name
	}

	if req.Color != "" {
		category.Color = req.Color
	}

	category.UpdatedAt = time.Now()

	_, err = o.Update(&category)
	if err != nil {
		log.Printf("[ERROR] [UserCategory.Update] 更新用户分类失败 - id: %d, userID: %d, name: %s, error: %v", category.Id, category.UserID, category.Name, err)
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "更新分类失败",
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

// 删除用户分类（需要 JWT 认证）
func (c *UserCategoryController) Delete() {
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

	categoryID, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if categoryID <= 0 {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": "分类 ID 错误",
		}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	
	// 查询分类并验证所有权
	var category models.UserCategory
	err = o.QueryTable(new(models.UserCategory)).Filter("id", categoryID).Filter("user_id", userID).One(&category)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    404,
			"message": "分类不存在或无权操作",
		}
		c.ServeJSON()
		return
	}

	// 检查是否有文章使用该分类
	var count int64
	count, _ = o.QueryTable(new(models.Article)).Filter("user_category_id", categoryID).Count()
	if count > 0 {
		c.Data["json"] = map[string]interface{}{
			"code":    400,
			"message": fmt.Sprintf("有 %d 篇文章正在使用该分类，无法删除", count),
		}
		c.ServeJSON()
		return
	}

	_, err = o.Delete(&category)
	if err != nil {
		log.Printf("[ERROR] [UserCategory.Delete] 删除用户分类失败 - id: %d, userID: %d, name: %s, error: %v", category.Id, category.UserID, category.Name, err)
		c.Data["json"] = map[string]interface{}{
			"code":    500,
			"message": "删除分类失败",
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
