package blog

import (
	"net/http"
	"strconv"
	"time"

	"scsPro/internal/common"
	"scsPro/internal/model"
	"github.com/gin-gonic/gin"
)

func DetailHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	article, err := model.GetArticleByID(uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "文章不存在: " + err.Error()})
		return
	}

	commonData, err := common.GetCommonData()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "获取公共数据失败: " + err.Error()})
		return
	}

	data := struct {
		Title           string
		Article         model.Article
		NavItems        []common.Module
		PopularArticles []model.Article
		UserStatus      string
		CurrentTime     time.Time
		IsHomePage      bool
	}{
		Title:           article.Title,
		Article:         article,
		NavItems:        commonData["NavItems"].([]common.Module),
		PopularArticles: commonData["PopularArticles"].([]model.Article),
		UserStatus:      commonData["UserStatus"].(string),
		CurrentTime:     commonData["CurrentTime"].(time.Time),
		IsHomePage:      false,
	}

	c.HTML(http.StatusOK, "base.html", data) // 改为 "base.html"
}

func CommentHandler(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "仅支持POST请求"})
		return
	}

	articleID, err := strconv.Atoi(c.PostForm("article_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID: " + err.Error()})
		return
	}

	parentIDStr := c.PostForm("parent_id")
	content := c.PostForm("content")
	username := c.PostForm("username")
	if username == "" {
		username = "匿名用户"
	}
	if content == "" {
		c.Redirect(http.StatusSeeOther, "/blog/detail?id="+strconv.Itoa(articleID))
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "添加评论失败: " + err.Error()})
		return
	}

	c.Redirect(http.StatusSeeOther, "/blog/detail?id="+strconv.Itoa(articleID))
}

func LikeHandler(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "仅支持POST请求"})
		return
	}

	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID: " + err.Error()})
		return
	}

	var article model.Article
	if err := model.DB.First(&article, id).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "文章不存在: " + err.Error()})
		return
	}

	article.Likes++
	if err := model.DB.Save(&article).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "点赞失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "likes": article.Likes})
}

func LikeCommentHandler(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "仅支持POST请求"})
		return
	}

	commentID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "无效的评论ID: " + err.Error()})
		return
	}

	var comment model.Comment
	if err := model.DB.First(&comment, commentID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "评论不存在: " + err.Error()})
		return
	}

	comment.Likes++
	if err := model.DB.Save(&comment).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "点赞失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "likes": comment.Likes})
}

func DislikeCommentHandler(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "仅支持POST请求"})
		return
	}

	commentID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "无效的评论ID: " + err.Error()})
		return
	}

	var comment model.Comment
	if err := model.DB.First(&comment, commentID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "评论不存在: " + err.Error()})
		return
	}

	comment.Dislikes++
	if err := model.DB.Save(&comment).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "点踩失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "dislikes": comment.Dislikes})
}