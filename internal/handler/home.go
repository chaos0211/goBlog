package handler

import (
	"net/http"
	"time"

	"scsPro/internal/common"
	"scsPro/internal/model"
	"github.com/gin-gonic/gin"
)

func HomeHandler(c *gin.Context) {
	modules := []common.Module{
		{Name: "个人简介", Description: "关于我的基本信息", Link: "/about"},
		{Name: "博客", Description: "我的文章", Link: "/blog"},
		{Name: "技能", Description: "我的技术栈", Link: "/skills"},
		{Name: "项目", Description: "我的作品集", Link: "/projects"},
		{Name: "爱好", Description: "我的兴趣爱好", Link: "/hobbies"},
		{Name: "时间线", Description: "成长记录", Link: "/timeline"},
		{Name: "资源", Description: "分享资源", Link: "/resources"},
		{Name: "联系我", Description: "与我联系", Link: "/contact"},
	}

	commonData, err := common.GetCommonData()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "获取公共数据失败: " + err.Error(),
		})
		return
	}

	data := struct {
		Title           string
		Modules         []common.Module
		NavItems        []common.Module
		PopularArticles []model.Article
		UserStatus      string
		CurrentTime     time.Time
		IsHomePage      bool
	}{
		Title:           "首页",
		Modules:         modules,
		NavItems:        commonData["NavItems"].([]common.Module),
		PopularArticles: commonData["PopularArticles"].([]model.Article),
		UserStatus:      commonData["UserStatus"].(string),
		CurrentTime:     commonData["CurrentTime"].(time.Time),
		IsHomePage:      true,
	}

	// 渲染 base.html，而不是 home.html
	c.HTML(http.StatusOK, "base.html", data)
}