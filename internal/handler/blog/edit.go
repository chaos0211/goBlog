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

func EditHandler(w http.ResponseWriter, r *http.Request) {
	commonData, err := handler.GetCommonData()
	if err != nil {
		http.Error(w, "获取公共数据失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 获取文章 ID
	idStr := r.URL.Query().Get("id")
	fmt.Printf("EditHandler: idStr=%s\n", idStr) // 调试：打印 id 参数
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		fmt.Println("EditHandler: No valid ID, creating new article")
		data := struct {
			Title           string
			NavItems        []handler.Module
			PopularArticles []model.Article
			UserStatus      string
			CurrentTime     time.Time
			IsHomePage      bool
			Article         model.Article
		}{
			Title:           "新建文章",
			NavItems:        commonData["NavItems"].([]handler.Module),
			PopularArticles: commonData["PopularArticles"].([]model.Article),
			UserStatus:      commonData["UserStatus"].(string),
			CurrentTime:     commonData["CurrentTime"].(time.Time),
			IsHomePage:      false,
			Article:         model.Article{},
		}
		renderEditTemplate(w, data)
		return
	}

	// 获取文章详情
	article, err := model.GetArticleByID(uint(id))
	if err != nil {
		http.Error(w, "文章不存在: "+err.Error(), http.StatusNotFound)
		return
	}

	fmt.Printf("EditHandler: Loaded article ID=%d, Title=%s\n", article.ID, article.Title) // 调试：打印文章信息
	data := struct {
		Title           string
		NavItems        []handler.Module
		PopularArticles []model.Article
		UserStatus      string
		CurrentTime     time.Time
		IsHomePage      bool
		Article         model.Article
	}{
		Title:           "编辑文章",
		NavItems:        commonData["NavItems"].([]handler.Module),
		PopularArticles: commonData["PopularArticles"].([]model.Article),
		UserStatus:      commonData["UserStatus"].(string),
		CurrentTime:     commonData["CurrentTime"].(time.Time),
		IsHomePage:      false,
		Article:         article,
	}

	renderEditTemplate(w, data)
}

func SaveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil || id <= 0 {
		// 新建文章
		title := r.FormValue("title")
		summary := r.FormValue("summary")
		content := r.FormValue("content")

		if title == "" || summary == "" || content == "" {
			http.Error(w, "标题、摘要和内容不能为空", http.StatusBadRequest)
			return
		}

		article := model.Article{
			Title:     title,
			Summary:   summary,
			Content:   content,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := model.DB.Create(&article).Error; err != nil {
			http.Error(w, "保存文章失败: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// 编辑已有文章
		var article model.Article
		if err := model.DB.First(&article, id).Error; err != nil {
			http.Error(w, "文章不存在: "+err.Error(), http.StatusNotFound)
			return
		}

		article.Title = r.FormValue("title")
		article.Summary = r.FormValue("summary")
		article.Content = r.FormValue("content")
		article.UpdatedAt = time.Now() // 更新时间刷新
		if err := model.DB.Save(&article).Error; err != nil {
			http.Error(w, "更新文章失败: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/blog", http.StatusSeeOther)
}

	// 辅助函数，渲染编辑模板
	func renderEditTemplate(w http.ResponseWriter, data interface{}) {
		tmpl := template.New("base").Funcs(template.FuncMap{
			"sub": sub,
			"add": add,
		})
		tmpl, err := tmpl.ParseFiles("templates/base.html", "templates/blog/edit.html")
		if err != nil {
			http.Error(w, "模板加载失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var buf bytes.Buffer
		err = tmpl.ExecuteTemplate(&buf, "base", data)
		if err != nil {
			http.Error(w, "模板渲染失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(buf.Bytes())
	}

