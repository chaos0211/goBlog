package blog

import (
	"net/http"
	"strconv"
	"time"

	"scsPro/internal/common"
	"scsPro/internal/model"
	"github.com/gin-gonic/gin"
)

func ListHandler(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}
	perPage := 10

	// 解析排序参数
	sortBy := c.Query("sort_by")
	if sortBy == "" {
		sortBy = "created_at"
	}
	order := c.Query("order")
	if order == "" || (order != "ASC" && order != "DESC") {
		order = "DESC"
	}

	// 获取文章列表
	articles, total, err := model.GetArticlesWithSort(page, perPage, sortBy, order)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "获取文章失败: " + err.Error()})
		return
	}

	// 计算总页数
	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	// 获取公共数据
	commonData, err := common.GetCommonData()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "获取公共数据失败: " + err.Error()})
		return
	}

	// 准备模板数据
	data := struct {
		Title           string
		Articles        []model.Article
		NavItems        []common.Module
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
		NavItems:        commonData["NavItems"].([]common.Module),
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
	c.HTML(http.StatusOK, "blog/list.html", data)
}