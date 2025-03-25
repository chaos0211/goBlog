package handler

import (
	"scsPro/internal/handler/blog"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/", HomeHandler)
	r.GET("/blog", blog.ListHandler)
	r.GET("/blog/detail", blog.DetailHandler)
	r.POST("/blog/comment", blog.CommentHandler)
	r.GET("/blog/edit", blog.EditHandler)
	r.POST("/blog/save", blog.SaveHandler)
	r.POST("/blog/like", blog.LikeHandler)
	r.POST("/blog/comment/like", blog.LikeCommentHandler)
	r.POST("/blog/comment/dislike", blog.DislikeCommentHandler)
}