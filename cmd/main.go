package main

import (
	"fmt"
	"scsPro/internal/config"
	"scsPro/internal/handler"
	"scsPro/internal/handler/blog"
	"scsPro/internal/model"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()
	if err := model.InitDB(); err != nil {
		fmt.Printf("初始化数据库失败: %v\n", err)
		return
	}




	r := gin.Default()


	blog.RegisterTemplateFuncs(r)

	r.LoadHTMLGlob("templates/*.html")
	fmt.Println("Loaded root templates successfully")
	r.LoadHTMLGlob("templates/blog/*.html")
	fmt.Println("Loaded blog templates successfully")

	r.Static("/static", "./static")
	// 调整模板加载，确保包含所有文件
	r.LoadHTMLGlob("templates/*.html")           // 加载根目录下的模板
	r.LoadHTMLGlob("templates/blog/*.html")      // 加载 blog 子目录下的模板
	handler.RegisterRoutes(r)

	fmt.Printf("服务器启动于 %s\n", config.AppConfig.Port)
	if err := r.Run(config.AppConfig.Port); err != nil {
		fmt.Printf("启动失败: %v\n", err)
	}



}