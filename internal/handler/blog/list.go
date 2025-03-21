package blog

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"scsPro/internal/handler"
	"scsPro/internal/model"
	"strconv"
	"time"
)

func ListHandler(w http.ResponseWriter, r *http.Request) {
	// 解析分页参数
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage := 10

	// 解析排序参数
	sortBy := r.URL.Query().Get("sort_by")
	if sortBy == "" {
		sortBy = "created_at" // 默认按创建时间排序
	}
	order := r.URL.Query().Get("order")
	if order == "" {
		order = "DESC" // 默认降序
	}
	if order != "ASC" && order != "DESC" {
		order = "DESC"
	}

	// 获取文章列表
	articles, total, err := model.GetArticlesWithSort(page, perPage, sortBy, order)
	if err != nil {
		fmt.Println("获取文章失败:", err)
		http.Error(w, "获取文章失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 计算总页数
	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	// 获取公共数据
	commonData, err := handler.GetCommonData()
	if err != nil {
		fmt.Println("获取公共数据失败:", err)
		http.Error(w, "获取公共数据失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := struct {
		Title           string
		Articles        []model.Article
		NavItems        []handler.Module
		Page            int
		Total           int
		PopularArticles []model.Article
		UserStatus      string
		CurrentTime     time.Time
		IsHomePage      bool
		SortBy          string
		Order           string
	}{
		Title:           "博客列表",
		Articles:        articles,
		NavItems:        commonData["NavItems"].([]handler.Module),
		Page:            page,
		Total:           totalPages,
		PopularArticles: commonData["PopularArticles"].([]model.Article),
		UserStatus:      commonData["UserStatus"].(string),
		CurrentTime:     commonData["CurrentTime"].(time.Time),
		IsHomePage:      false,
		SortBy:          sortBy,
		Order:           order,
	}

	// 渲染模板
	tmpl := template.New("base").Funcs(template.FuncMap{
		"sub": sub,
		"add": add,
	})
	tmpl, err = tmpl.ParseFiles("templates/base.html", "templates/blog/list.html")
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