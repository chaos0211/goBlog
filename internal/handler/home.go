package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"scsPro/internal/model"
	"time"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	modules := []Module{
		{Name: "个人简介", Description: "关于我的基本信息", Link: "/about"},
		{Name: "博客", Description: "我的文章", Link: "/blog"},
		{Name: "技能", Description: "我的技术栈", Link: "/skills"},
		{Name: "项目", Description: "我的作品集", Link: "/projects"},
		{Name: "爱好", Description: "我的兴趣爱好", Link: "/hobbies"},
		{Name: "时间线", Description: "成长记录", Link: "/timeline"},
		{Name: "资源", Description: "分享资源", Link: "/resources"},
		{Name: "联系我", Description: "与我联系", Link: "/contact"},
	}

	commonData, err := GetCommonData()
	if err != nil {
		fmt.Println("获取公共数据失败:", err)
		http.Error(w, "获取公共数据失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Title           string
		Modules         []Module
		NavItems        []Module
		PopularArticles []model.Article
		UserStatus      string
		CurrentTime     time.Time
		IsHomePage      bool
	}{
		Title:           "首页",
		Modules:         modules,
		NavItems:        commonData["NavItems"].([]Module),
		PopularArticles: commonData["PopularArticles"].([]model.Article),
		UserStatus:      commonData["UserStatus"].(string),
		CurrentTime:     commonData["CurrentTime"].(time.Time),
		IsHomePage:      true, // 标记为首页
	}

	//fmt.Printf("Modules: %+v\n", modules)

	tmpl, err := template.ParseFiles("templates/base.html", "templates/home.html")
	if err != nil {
		fmt.Println("模板加载失败:", err)
		http.Error(w, "模板加载失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "base", data)
	if err != nil {
		fmt.Println("模板渲染失败:", err)
		http.Error(w, "模板渲染失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(buf.Bytes())
}