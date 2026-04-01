package main

import (
	"blog-system/models"
	_ "blog-system/routers"
	"log"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	// 从配置文件读取数据库配置
	dbDriver, _ := web.AppConfig.String("dbdriver")
	dbConn, _ := web.AppConfig.String("dbconn")

	log.Printf("[INFO] [Main.Init] 正在初始化数据库连接 - driver: %s", dbDriver)
	
	// 注册数据库驱动和连接
	orm.RegisterDriver(dbDriver, orm.DRMySQL)
	err := orm.RegisterDataBase("default", dbDriver, dbConn)
	if err != nil {
		log.Fatalf("[FATAL] [Main.Init] 数据库连接失败 - driver: %s, error: %v", dbDriver, err)
	}

	// 设置最大连接数（从配置文件读取）
	maxIdle := web.AppConfig.DefaultInt("dbmaxidle", 10)
	maxOpen := web.AppConfig.DefaultInt("dbmaxopen", 100)
	orm.SetMaxIdleConns("default", maxIdle)
	orm.SetMaxOpenConns("default", maxOpen)
	log.Printf("[INFO] [Main.Init] 数据库连接池配置 - maxIdle: %d, maxOpen: %d", maxIdle, maxOpen)

	// 注册模型
	orm.RegisterModel(new(models.Category), new(models.Tag), new(models.Article), new(models.Comment), new(models.User), new(models.UserCategory), new(models.Like), new(models.Follow))
	log.Printf("[INFO] [Main.Init] 数据模型已注册")
}

func main() {
	log.Println("[INFO] [Main] 开始检查数据库状态...")

	o := orm.NewOrm()
	
	// 测试数据库连接是否正常
	var result int
	err := o.Raw("SELECT 1").QueryRow(&result)
	if err != nil {
		log.Fatalf("[FATAL] [Main] ❌ 数据库连接失败 - error: %v", err)
	}
	
	log.Println("[INFO] [Main] ✅ 数据库连接正常")

	// 启动 Web 服务
	log.Println("[INFO] [Main] 正在启动 Web 服务...")
	
	// 配置静态文件服务 - Beego v2 的正确方式
	web.BConfig.WebConfig.StaticDir["/uploads"] = "uploads"
	log.Println("[INFO] [Main] 静态文件服务已配置：/uploads -> ./uploads")
	
	web.Run()
}

