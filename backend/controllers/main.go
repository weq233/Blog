package controllers

import (
	"blog-system/models"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
)

// 主控制器 - 首页
type MainController struct {
	web.Controller
}

func (c *MainController) Index() {
	var articles []*models.Article
	o := orm.NewOrm()
	o.QueryTable(new(models.Article)).Filter("status", 1).OrderBy("-created_at").Limit(10).All(&articles)

	c.Data["Articles"] = articles
	c.Data["PageTitle"] = "首页"

}
