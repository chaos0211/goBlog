package blog

import (
	"net/http"
	"strconv"
	"time"

	"scsPro/internal/common"
	"scsPro/internal/model"
	"github.com/gin-gonic/gin"
)

func EditHandler(c *gin.Context) {
	commonData, err := common.GetCommonData()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "获取公共数据失败: " + err.Error()})
		return
	}

	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	var data struct {
		Title           string
		NavItems        []common.Module
		PopularArticles []model.Article
		UserStatus      string
		CurrentTime     time.Time
		IsHomePage      bool
		Article         model.Article
	}

	if err != nil || id <= 0 {
		data = struct {
			Title           string
			NavItems        []common.Module
			PopularArticles []model.Article
			UserStatus      string
			CurrentTime     time.Time
			IsHomePage      bool
			Article         model.Article
		}{
			Title:           "新建文章",
			NavItems:        commonData["NavItems"].([]common.Module),
			PopularArticles: commonData["PopularArticles"].([]model.Article),
			UserStatus:      commonData["UserStatus"].(string),
			CurrentTime:     commonData["CurrentTime"].(time.Time),
			IsHomePage:      false,
			Article:         model.Article{},
		}
	} else {
		article, err := model.GetArticleByID(uint(id))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "文章不存在: " + err.Error()})
			return
		}
		data = struct {
			Title           string
			NavItems        []common.Module
			PopularArticles []model.Article
			UserStatus      string
			CurrentTime     time.Time
			IsHomePage      bool
			Article         model.Article
		}{
			Title:           "编辑文章",
			NavItems:        commonData["NavItems"].([]common.Module),
			PopularArticles: commonData["PopularArticles"].([]model.Article),
			UserStatus:      commonData["UserStatus"].(string),
			CurrentTime:     commonData["CurrentTime"].(time.Time),
			IsHomePage:      false,
			Article:         article,
		}
	}

	c.HTML(http.StatusOK, "blog/edit.html", data)
}

func SaveHandler(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "仅支持POST请求"})
		return
	}

	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil || id <= 0 {
		// 新建文章
		title := c.PostForm("title")
		summary := c.PostForm("summary")
		content := c.PostForm("content")
		if title == "" || summary == "" || content == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "标题、摘要和内容不能为空"})
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
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "保存文章失败: " + err.Error()})
			return
		}
	} else {
		// 编辑文章
		var article model.Article
		if err := model.DB.First(&article, id).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "文章不存在: " + err.Error()})
			return
		}

		article.Title = c.PostForm("title")
		article.Summary = c.PostForm("summary")
		article.Content = c.PostForm("content")
		article.UpdatedAt = time.Now()
		if err := model.DB.Save(&article).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "更新文章失败: " + err.Error()})
			return
		}
	}

	c.Redirect(http.StatusSeeOther, "/blog")
}