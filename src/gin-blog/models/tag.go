package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Tag struct {
	Model

	TagId    int    `json:"tag_id" gorm:"index"`
	Name     string `json:"name"`
	CreateBy int    `json:"create_by"`
	UpdateBy int    `json:"update_by"`
	State    int    `json:"state"`
}

func (tag *Tag) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreateAt", time.Now().Unix())
	return nil
}

func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdateAt", time.Now().Unix())
	return nil
}

func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag) {
	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)

	return
}

func GetTagTotal(maps interface{}) (count int) {
	db.Model(&Tag{}).Where(maps).Count(&count)

	return
}

func ExistTagByName(name string) bool {
	var tag Tag
	db.Select("id").Where("name=?", name).First(&tag)
	if tag.ID > 0 {
		return true
	}

	return false
}

func ExistTagByTagId(tagId int) bool {
	var tag Tag
	db.Select("tag_id").Where("tag_id=?", tagId).First(&tag)
	if tag.TagId > 0 {
		return true
	}
	return false
}

func AddTag(name string, state int, createBy int) bool {
	tagId := genId()

	db.Create(&Tag{
		TagId:    tagId,
		Name:     name,
		State:    state,
		CreateBy: createBy,
		UpdateBy: createBy,
	})

	return true
}

func DeleteTag(tagId int) bool {
	db.Where("tag_id=?", tagId).Delete(&Tag{})

	return true
}

func EditTag(tagId int, data interface{}) bool {
	db.Model(&Tag{}).Where("tag_id=?", tagId).Update(data)

	return true
}

func genId() (uuid int) {
	uuid = int(time.Now().Unix())
	return
}
