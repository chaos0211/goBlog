package model

import (
	"fmt"
	"sort"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Article 文章模型
type Article struct {
	ID          uint           `gorm:"primaryKey"`
	Title       string         `gorm:"type:varchar(255);not null"`
	Summary     string         `gorm:"type:text;not null"`
	Content     string         `gorm:"type:text;not null"`
	Author      string         `gorm:"type:varchar(100);not null;default:'佚名'"`
	Views       uint           `gorm:"default:0"` // 阅读量
	Likes       uint           `gorm:"default:0"` // 点赞数
	CreatedAt   time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Comments    []Comment      `gorm:"foreignKey:ArticleID"`
}

// Comment 模型
type Comment struct {
	ID        uint      `gorm:"primaryKey"`
	ArticleID uint
	ParentID  *uint
	Content   string
	Username  string
	Likes     int
	Dislikes  int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
	Replies   []Comment `gorm:"-"` // 忽略数据库映射，动态填充
}

// 确保 Replies 字段初始化
func (c *Comment) InitReplies() {
	if c.Replies == nil {
		c.Replies = make([]Comment, 0)
	}
}

var DB *gorm.DB

// InitDB 初始化 GORM 数据库连接
func InitDB() error {
	dsn := "root:123456@tcp(127.0.0.1:33309)/goBlog?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	// 自动迁移表结构
	if err := db.AutoMigrate(&Article{}, &Comment{}); err != nil {
		return fmt.Errorf("自动迁移失败: %v", err)
	}

	// 赋值给全局变量
	DB = db
	fmt.Println("数据库连接成功")
	return nil
}


// GetArticlesWithSort 获取文章列表，支持分页和排序
func GetArticlesWithSort(page, perPage int, sortBy, order string) ([]Article, int64, error) {
	if DB == nil {
		return nil, 0, fmt.Errorf("数据库未初始化")
	}

	var articles []Article
	var total int64

	// 计算偏移量
	offset := (page - 1) * perPage

	// 构建查询
	query := DB.Model(&Article{})

	// 根据 sortBy 和 order 调整排序
	switch sortBy {
	case "created_at":
		query = query.Order("created_at " + order)
	case "updated_at":
		query = query.Order("updated_at " + order)
	case "latest_comment":
		// 使用子查询获取每篇文章的最新评论时间
		query = query.Select("articles.*, (SELECT MAX(created_at) FROM comments WHERE comments.article_id = articles.id AND comments.deleted_at IS NULL) as latest_comment_time").
			Order("latest_comment_time " + order + ", created_at DESC")
	case "views":
		query = query.Order("views " + order + ", created_at DESC")
	case "likes":
		query = query.Order("likes " + order + ", created_at DESC")
	case "comment_count":
		// 使用子查询获取每篇文章的评论数量
		query = query.Select("articles.*, (SELECT COUNT(*) FROM comments WHERE comments.article_id = articles.id AND comments.deleted_at IS NULL) as comment_count").
			Order("comment_count " + order + ", created_at DESC")
	default:
		// 默认按创建时间降序
		query = query.Order("created_at DESC")
	}

	// 获取总数
	if err := DB.Model(&Article{}).Where("deleted_at IS NULL").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := query.Where("deleted_at IS NULL").Offset(offset).Limit(perPage).Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	// 预加载评论
	for i := range articles {
		if err := DB.Model(&articles[i]).Where("deleted_at IS NULL").Association("Comments").Find(&articles[i].Comments); err != nil {
			return nil, 0, err
		}
	}

	return articles, total, nil
}

// GetArticles 获取文章列表
// GetArticles 获取文章列表（调用 GetArticlesWithSort）
func GetArticles(page, perPage int) ([]Article, error) {
	articles, _, err := GetArticlesWithSort(page, perPage, "created_at", "DESC")
	return articles, err
}

func GetArticleByID(id uint) (Article, error) {
	if DB == nil {
		return Article{}, fmt.Errorf("数据库未初始化")
	}

	var article Article
	if err := DB.First(&article, id).Error; err != nil {
		return Article{}, err
	}

	// 加载所有评论，排除软删除的记录
	var allComments []Comment
	if err := DB.Where("article_id = ? AND deleted_at IS NULL", id).Find(&allComments).Error; err != nil {
		return Article{}, err
	}
	fmt.Printf("GetArticleByID: Loaded %d comments for article ID=%d\n", len(allComments), id)

	// 打印每条评论的 ID 和 ParentID
	//for _, comment := range allComments {
	//	parentID := "NULL"
	//	if comment.ParentID != nil {
	//		parentID = fmt.Sprintf("%d", *comment.ParentID)
	//	}
	//	fmt.Printf("Comment ID=%d, ParentID=%s\n", comment.ID, parentID)
	//}

	// 构建评论树
	commentMap := make(map[uint]*Comment)
	var topLevelComments []*Comment // 修改为指针切片

	// 第一遍：收集所有评论到 commentMap，并初始化 Replies
	for i := range allComments {
		comment := &allComments[i]
		comment.InitReplies()
		commentMap[comment.ID] = comment
	}

	// 第二遍：构建评论树，直接使用 commentMap 中的引用
	for _, comment := range allComments {
		if comment.ParentID == nil {
			topLevelComments = append(topLevelComments, commentMap[comment.ID])
		} else {
			parent, exists := commentMap[*comment.ParentID]
			if exists {
				parent.Replies = append(parent.Replies, *commentMap[comment.ID])
				//fmt.Printf("Added comment ID=%d to parent ID=%d, parent now has %d replies\n", comment.ID, *comment.ParentID, len(parent.Replies))
			} else {
				//fmt.Printf("Parent comment %d not found for comment %d\n", *comment.ParentID, comment.ID)
			}
		}
	}

	// 打印顶级评论和回复数量
	//fmt.Printf("GetArticleByID: Loaded %d top-level comments for article ID=%d\n", len(topLevelComments), id)
	//for _, comment := range topLevelComments {
	//	fmt.Printf("Comment ID=%d has %d replies\n", comment.ID, len(comment.Replies))
	//}

	// 按点赞数和创建时间排序回复
	for _, comment := range topLevelComments {
		sortReplies(comment)
	}

	// 将顶级评论赋值给 article.Comments
	articleComments := make([]Comment, len(topLevelComments))
	for i, comment := range topLevelComments {
		articleComments[i] = *comment
	}
	article.Comments = articleComments

	// 增加阅读量
	article.Views++
	if err := DB.Save(&article).Error; err != nil {
		return Article{}, err
	}

	return article, nil
}

func sortReplies(comment *Comment) {
	if len(comment.Replies) == 0 {
		return
	}
	sort.Slice(comment.Replies, func(i, j int) bool {
		if comment.Replies[i].Likes == comment.Replies[j].Likes {
			return comment.Replies[i].CreatedAt.Before(comment.Replies[j].CreatedAt) // 按时间升序（最早的在前）
		}
		return comment.Replies[i].Likes > comment.Replies[j].Likes // 按点赞数降序
	})
	for i := range comment.Replies {
		sortReplies(&comment.Replies[i]) // 递归排序嵌套回复
	}
}

// AddComment 添加评论
func AddComment(articleID uint, parentID *uint, content, username string) error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	comment := Comment{
		ArticleID: articleID,
		ParentID:  parentID,
		Content:   content,
		Username:  username,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return DB.Create(&comment).Error
}

// GetTotalArticles 获取文章总数
func GetTotalArticles() (int64, error) {
	if DB == nil {
		return 0, fmt.Errorf("数据库未初始化")
	}

	var count int64
	if err := DB.Model(&Article{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetPopularArticles 获取热门文章（按阅读量排序）
func GetPopularArticles(limit int) ([]Article, error) {
	if DB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	var articles []Article
	if err := DB.Order("views DESC").Limit(limit).Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

// GenerateArticles 批量生成文章
func GenerateArticles(count int) error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	// 伪造数据
	articles := make([]Article, count)
	for i := 0; i < count; i++ {
		articles[i] = Article{
			Title:     fmt.Sprintf("文章标题 %d", i+1),
			Summary:   fmt.Sprintf("这是第 %d 篇文章的摘要", i+1),
			Content:   fmt.Sprintf("这是第 %d 篇文章的详细内容", i+1),
			Author:    fmt.Sprintf("作者%d", (i%3)+1), // 随机分配3个作者
			CreatedAt: time.Now().Add(-time.Duration(i) * time.Hour), // 按时间倒序
			UpdatedAt: time.Now().Add(-time.Duration(i) * time.Hour),
		}
	}

	// 批量插入
	if err := DB.Create(&articles).Error; err != nil {
		return fmt.Errorf("批量插入文章失败: %v", err)
	}
	fmt.Printf("成功生成了 %d 篇文章\n", count)
	return nil
}