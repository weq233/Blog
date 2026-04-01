package models

import (
	"time"
)

// 用户模型
type User struct {
	Id             int       `orm:"pk;auto" json:"id"`
	Username       string    `orm:"size(50);unique" json:"username"`
	Password       string    `orm:"size(255)" json:"-"`
	Email          string    `orm:"size(100);unique" json:"email"`
	Nickname       string    `orm:"size(100);null" json:"nickname"`
	Avatar         string    `orm:"size(500);null" json:"avatar"`
	Role           int       `orm:"default(1)" json:"role"`
	Status         int       `orm:"default(1)" json:"status"`
	FailedAttempts int       `orm:"default(0)" json:"failed_attempts"`
	LockedUntil    time.Time `orm:"null;type(datetime)" json:"locked_until"`
	CreatedAt      time.Time `orm:"auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt      time.Time `orm:"auto_now;type(datetime)" json:"updated_at"`
}

// 分类模型
type Category struct {
	Id          int       `orm:"pk;auto" json:"id"`
	Name        string    `orm:"size(100)" json:"name"`
	Slug        string    `orm:"size(100);unique" json:"slug"`
	ArticleCount int      `orm:"-" json:"article_count"` // 文章数量（非数据库字段）
	CreatedAt   time.Time `orm:"auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt   time.Time `orm:"auto_now;type(datetime)" json:"updated_at"`
}

// 标签模型
type Tag struct {
	Id          int       `orm:"pk;auto" json:"id"`
	Name        string    `orm:"size(50);unique" json:"name"`
	Slug        string    `orm:"size(50);unique" json:"slug"`
	ArticleCount int     `orm:"-" json:"article_count"` // 文章数量（非数据库字段）
	CreatedAt   time.Time `orm:"auto_now_add;type(datetime)" json:"created_at"`
}

// 用户自定义分类模型（私有分类）
type UserCategory struct {
	Id        int       `orm:"pk;auto" json:"id"`
	UserID    int       `orm:"column(user_id);index" json:"user_id"` // 所属用户 ID
	Name      string    `orm:"size(100)" json:"name"` // 分类名称
	Color     string    `orm:"size(20);default(#409EFF)" json:"color"` // 分类颜色标识
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime)" json:"updated_at"`
}

// 文章模型
type Article struct {
	Id            int           `orm:"pk;auto" json:"id"`
	Title         string        `orm:"size(200)" json:"title"`
	Slug          string        `orm:"size(200);unique" json:"slug"`
	Summary       string        `orm:"size(500)" json:"summary"`
	Content       string        `orm:"type(text)" json:"content"`
	CoverImage    string        `orm:"size(500);null" json:"cover_image"`
	ViewCount     int           `orm:"default(0)" json:"view_count"`
	Status        int           `orm:"default(1)" json:"status"`
	Author        *User         `orm:"rel(fk);column(author_id)" json:"author"`
	Category      *Category     `orm:"rel(fk);column(category_id);null" json:"category"` // 系统分类（管理员管理）
	UserCategory  *UserCategory `orm:"rel(fk);column(user_category_id);null" json:"user_category"` // 用户自定义分类
	Tags          []*Tag        `orm:"rel(m2m)" json:"tags"`
	CreatedAt     time.Time     `orm:"auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt     time.Time     `orm:"auto_now;type(datetime)" json:"updated_at"`
}

// 评论模型
type Comment struct {
	Id        int       `orm:"pk;auto" json:"id"`
	Article   *Article  `orm:"rel(fk);index" json:"article"`
	User      *User     `orm:"rel(fk);column(user_id);null" json:"user"`
	Author    string    `orm:"size(100);null" json:"author"`
	Email     string    `orm:"size(100);null" json:"email"`
	Content   string    `orm:"type(text)" json:"content"`
	Status    int       `orm:"default(1)" json:"status"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)" json:"created_at"`
}

// 点赞模型
type Like struct {
	Id        int       `orm:"pk;auto" json:"id"`
	Article   *Article  `orm:"rel(fk);index" json:"article"`
	User      *User     `orm:"rel(fk);column(user_id);index" json:"user"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)" json:"created_at"`
}

// 关注模型
type Follow struct {
	Id         int       `orm:"pk;auto" json:"id"`
	Follower   *User     `orm:"rel(fk);column(follower_id);index" json:"follower"` // 关注者
	Followee   *User     `orm:"rel(fk);column(followee_id);index" json:"followee"` // 被关注者
	CreatedAt  time.Time `orm:"auto_now_add;type(datetime)" json:"created_at"`
}
