package api

import (
	"log"
	"net/http"

	"gin-blog/src/gin-blog/models"
	"gin-blog/src/gin-blog/pkg/e"
	"gin-blog/src/gin-blog/pkg/setting"
	"gin-blog/src/gin-blog/pkg/util"

	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

//获取单个文章
func GetArticle(c *gin.Context) {
	articleId := com.StrTo(c.Param("article_id")).MustInt()

	valid := validation.Validation{}
	valid.Min(articleId, 1, "article_id").Message("article_id必须大于0")

	code := e.INVALID_PARAMS
	var data interface{}
	if !valid.HasErrors() {
		if models.ExistArticleById(articleId) {
			data = models.GetArticle(articleId)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key:%s, err.message:%s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

//获取文章列表
func GetArticles(c *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state

		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	var tagId int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		maps["tag_id"] = tagId
		valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	}

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS
		data["lists"] = models.GetArticles(util.GetPage(c), setting.PageSize, maps)
		data["total"] = models.GetArticleTotal(maps)
	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

//新增文章
func AddArticle(c *gin.Context) {
	tagId := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	cover := c.Query("cover")
	detail := c.Query("detail")
	createBy := com.StrTo(c.Query("create_by")).MustInt()
	state := com.StrTo(c.DefaultQuery("state", "1")).MustInt()

	valid := validation.Validation{}
	valid.Min(tagId, 1, "tag_id").Message("标签id必须大于0")
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(detail, "detail").Message("详情不能为空")
	valid.Required(createBy, "create_by").Message("创建人不能为空")
	valid.Range(state, 0, 1, "state").Message("状态只能为0或1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistTagByTagId(tagId) {
			data := make(map[string]interface{})
			data["tag_id"] = tagId
			data["title"] = title
			data["cover"] = cover
			data["desc"] = desc
			data["detail"] = detail
			data["create_by"] = createBy
			data["state"] = state

			models.AddArticle(data)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}

//修改文章
func EditArticle(c *gin.Context) {
	articleId := com.StrTo(c.Param("article_id")).MustInt()
	code := e.INVALID_PARAMS
	if articleId <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": make(map[string]interface{}),
		})
		return
	}
	if !models.ExistArticleById(articleId) {
		code := e.ERROR_NOT_EXIST_ARTICLE
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": make(map[string]interface{}),
		})
		return
	}

	tagId := com.StrTo(c.PostForm("tag_id")).MustInt()
	title := c.PostForm("title")
	cover := c.PostForm("cover")
	desc := c.PostForm("desc")
	detail := c.PostForm("detail")
	state := com.StrTo(c.PostForm("state")).MustInt()

	data := make(map[string]interface{})
	if tagId > 0 {
		if !models.ExistTagByTagId(tagId) {
			code := e.ERROR_NOT_EXIST_TAG
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": make(map[string]interface{}),
			})
			return
		}
		data["tag_id"] = tagId
	}
	if len(title) > 0 {
		if len(title) > 100 || len(title) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": make(map[string]interface{}),
			})
			return
		} else {
			data["title"] = title
		}
	}
	// 没有好办法判断是否传如cover参数，所以如果不想修改数据库中cover的值，前端需传入"-1"
	if cover != "-1" {
		data["cover"] = cover
	}
	if len(desc) > 0 {
		if len(desc) > 200 || len(desc) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": make(map[string]interface{}),
			})
			return
		} else {
			data["desc"] = desc
		}
	}
	if len(detail) > 0 {
		if len(detail) > 200 || len(detail) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": make(map[string]interface{}),
			})
			return
		} else {
			data["detail"] = detail
		}
	}
	//同cover
	if state != -1 {
		data["state"] = state
	}

	models.EditArticle(articleId, data)
	code = e.SUCCESS
	//code = 1
	//msg := ""
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

//删除文章
func DeleteArticle(c *gin.Context) {
	articleId := com.StrTo(c.Query("article_id")).MustInt()

	valid := validation.Validation{}
	valid.Min(articleId, 1, "article_id").Message("文章id必须大于0")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistArticleById(articleId) {
			models.DeleteArticle(articleId)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}
