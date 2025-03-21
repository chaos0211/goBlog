package blog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"scsPro/internal/handler"
	"scsPro/internal/model"
	"strconv"
	"time"
)

func DetailHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	article, err := model.GetArticleByID(uint(id))
	if err != nil {
		http.Error(w, "文章不存在: "+err.Error(), http.StatusNotFound)
		return
	}

	commonData, err := handler.GetCommonData()
	if err != nil {
		http.Error(w, "获取公共数据失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Title           string
		Article         model.Article
		NavItems        []handler.Module
		PopularArticles []model.Article
		UserStatus      string
		CurrentTime     time.Time
		IsHomePage      bool
	}{
		Title:           article.Title,
		Article:         article,
		NavItems:        commonData["NavItems"].([]handler.Module),
		PopularArticles: commonData["PopularArticles"].([]model.Article),
		UserStatus:      commonData["UserStatus"].(string),
		CurrentTime:     commonData["CurrentTime"].(time.Time),
		IsHomePage:      false, // 非首页
	}

	tmpl := template.New("base").Funcs(template.FuncMap{
		"sub": sub,
		"add": add,
	})
	tmpl, err = tmpl.ParseFiles("templates/base.html", "templates/blog/detail.html")
	if err != nil {
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

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持POST请求", http.StatusMethodNotAllowed)
		return
	}
	articleID, err := strconv.Atoi(r.FormValue("article_id"))
	if err != nil {
		http.Error(w, "无效的文章ID: "+err.Error(), http.StatusBadRequest)
		return
	}
	parentIDStr := r.FormValue("parent_id")
	content := r.FormValue("content")
	username := r.FormValue("username")
	if username == "" {
		username = "匿名用户"
	}
	if content == "" {
		http.Redirect(w, r, "/blog/detail?id="+strconv.Itoa(articleID), http.StatusSeeOther)
		return
	}

	var parentID *uint
	if parentIDStr != "" {
		pid, err := strconv.Atoi(parentIDStr)
		if err == nil && pid > 0 {
			p := uint(pid)
			parentID = &p
		}
	}

	if err := model.AddComment(uint(articleID), parentID, content, username); err != nil {
		http.Error(w, "添加评论失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 直接重新加载文章并渲染
	article, err := model.GetArticleByID(uint(articleID))
	if err != nil {
		http.Error(w, "加载文章失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	commonData, err := handler.GetCommonData()
	if err != nil {
		http.Error(w, "获取公共数据失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Title           string
		NavItems        []handler.Module
		PopularArticles []model.Article
		UserStatus      string
		CurrentTime     time.Time
		IsHomePage      bool
		Article         model.Article
	}{
		Title:           article.Title,
		NavItems:        commonData["NavItems"].([]handler.Module),
		PopularArticles: commonData["PopularArticles"].([]model.Article),
		UserStatus:      commonData["UserStatus"].(string),
		CurrentTime:     commonData["CurrentTime"].(time.Time),
		IsHomePage:      false,
		Article:         article,
	}

	tmpl := template.New("base").Funcs(template.FuncMap{
		"sub": sub,
		"add": add,
	})
	tmpl, err = tmpl.ParseFiles("templates/base.html", "templates/blog/detail.html")
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

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持POST请求", http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "无效的文章ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	var article model.Article
	if err := model.DB.First(&article, id).Error; err != nil {
		http.Error(w, "文章不存在: "+err.Error(), http.StatusNotFound)
		return
	}

	article.Likes++
	if err := model.DB.Save(&article).Error; err != nil {
		http.Error(w, "点赞失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回 JSON 响应
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success": true,
		"likes":   article.Likes,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "响应编码失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持POST请求", http.StatusMethodNotAllowed)
		return
	}
	commentID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "无效的评论ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	var comment model.Comment
	if err := model.DB.First(&comment, commentID).Error; err != nil {
		http.Error(w, "评论不存在: "+err.Error(), http.StatusNotFound)
		return
	}

	comment.Likes++
	if err := model.DB.Save(&comment).Error; err != nil {
		http.Error(w, "点赞失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回 JSON 响应
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success": true,
		"likes":   comment.Likes,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "响应编码失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func DislikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持POST请求", http.StatusMethodNotAllowed)
		return
	}
	commentID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "无效的评论ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	var comment model.Comment
	if err := model.DB.First(&comment, commentID).Error; err != nil {
		http.Error(w, "评论不存在: "+err.Error(), http.StatusNotFound)
		return
	}

	comment.Dislikes++
	if err := model.DB.Save(&comment).Error; err != nil {
		http.Error(w, "点踩失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回 JSON 响应
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success":  true,
		"dislikes": comment.Dislikes,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "响应编码失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

