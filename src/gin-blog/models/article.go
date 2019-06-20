package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Article struct {
	Model

	TagId int `json:"tag_id" gorm:"index"`
	Tag   Tag `json:"tag"`

	ArticleId int    `json: "article_id" gorm:"index"`
	Title     string `json: "title"`
	Cover     string `json: "cover"`
	Desc      string `json: "desc"`
	Detail    string `json: "detail"`
	CreateBy  int    `json: "create_by"`
	State     int    `json:"state"`
}

func (article *Article) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreateAt", time.Now().Unix())

	return nil
}

func (article *Article) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdateAt", time.Now().Unix())

	return nil
}

func ExistArticleById(articleId int) bool {
	var article Article
	db.Select("article_id").Where("article_id=?", articleId).First(&article)

	if article.ArticleId > 0 {
		return true
	}
	return false
}

func GetArticleTotal(maps interface{}) (count int) {
	db.Model(&Article{}).Where(maps).Count(&count)

	return
}

func GetArticles(pageNum int, pageSize int, maps interface{}) (articles []Article) {
	db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles)

	return
}

func GetArticle(articleId int) (article Article) {
	db.Where("article_id=?", articleId).First(&article)
	db.Model(&article).Related(&article.Tag)

	return
}

func EditArticle(articleId int, data interface{}) bool {
	db.Model(&Article{}).Where("article_id=?", articleId).Update(data)

	return true
}

/**
func genId() (uuid int) {
	uuid = int(time.Now().Unix())
	return
}
*/

func AddArticle(data map[string]interface{}) bool {
	articleId := genId()
	db.Create(&Article{
		ArticleId: articleId,
		TagId:     data["tag_id"].(int),
		Title:     data["title"].(string),
		Cover:     data["cover"].(string),
		Desc:      data["desc"].(string),
		Detail:    data["detail"].(string),
		CreateBy:  data["create_by"].(int),
		State:     data["state"].(int),
	})

	return true
}

func DeleteArticle(articleId int) bool {
	db.Where("article_id=?", articleId).Delete(&Article{})

	return true
}
