package main

import (
	"fmt"
	"html/template"
	"net/http"
	"scsPro/internal/handler"
	"scsPro/internal/handler/blog"
	"scsPro/internal/model"
	"time"
	//"scsPro/internal/model"
)

type PageData struct {
	Title    string
	Modules  []Module
	NavItems []string
	Time     string
}

type Module struct {
	Name        string
	Description string
	Link        string
}


func main() {
	// 初始化数据库
	if err := model.InitDB(); err != nil {
		fmt.Printf("初始化数据库失败: %v\n", err)
		return
	}

	// 批量生成30篇文章（仅在数据库为空时执行）
	var count int64
	if err := model.DB.Model(&model.Article{}).Count(&count).Error; err != nil {
		fmt.Printf("检查文章数量失败: %v\n", err)
		return
	}
	if count == 0 {
		if err := model.GenerateArticles(30); err != nil {
			fmt.Printf("生成文章失败: %v\n", err)
			return
		}
	}

	// 静态文件服务
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// 路由注册
	http.HandleFunc("/", handler.HomeHandler)
	http.HandleFunc("/blog", blog.ListHandler)
	http.HandleFunc("/blog/detail", blog.DetailHandler)
	http.HandleFunc("/blog/comment", blog.CommentHandler)
	http.HandleFunc("/blog/edit", blog.EditHandler)    // 新增编辑路由
	http.HandleFunc("/blog/save", blog.SaveHandler)
	http.HandleFunc("/blog/like", blog.LikeHandler)  // 新增点赞路由
	http.HandleFunc("/blog/comment/like", blog.LikeCommentHandler)   //评论区点赞和踩
	http.HandleFunc("/blog/comment/dislike", blog.DislikeCommentHandler)

	//http.HandleFunc("/about", placeholderHandler)     // 占位
	//http.HandleFunc("/skills", placeholderHandler)    // 占位
	//http.HandleFunc("/projects", placeholderHandler)  // 占位
	//http.HandleFunc("/hobbies", placeholderHandler)   // 占位
	//http.HandleFunc("/timeline", placeholderHandler)  // 占位
	//http.HandleFunc("/resources", placeholderHandler) // 占位
	//http.HandleFunc("/contact", placeholderHandler)   // 占位

	fmt.Println("服务器启动于 :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("启动失败:", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	modules := []Module{
		{Name: "个人简介", Description: "关于我的基本信息", Link: "/profile"},
		{Name: "技能展示", Description: "我的技术栈", Link: "/skills"},
	}
	navItems := []string{"首页", "博客"}
	data := PageData{
		Title:    "欢迎来到我的个人网站",
		Modules:  modules,
		NavItems: navItems,
		Time:     time.Now().Format("2006-01-02 15:04:05"),
	}

	tmpl, err := template.ParseFiles("../templates/base.html", "../templates/home.html")
	if err != nil {
		fmt.Println("模板加载错误:", err)
		http.Error(w, "模板加载失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		fmt.Println("模板渲染错误:", err)
		http.Error(w, "模板渲染失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
}