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
	//log.Printf("ID=======%d", articleId)
	valid := validation.Validation{}
	valid.Min(articleId, 1, "article_id").Message("文章ID必须大于0")
	var state = -1
	if arg := c.Query("state"); arg != "" {
		valid.Range(arg, 0, 1, "state").Message("状态只能为0或者1")
		state = com.StrTo(arg).MustInt()
	}
	var tagId = -1
	if arg := c.Query("tag_id"); arg != "" {
		//valid.Min(tagId, 1, "tag_id").Message("标签ID不能为空")
		tagId = com.StrTo(arg).MustInt()
	}
	var cover string
	if arg := c.Query("cover"); arg != "0" {
		cover = arg
	}
	var title, desc, detail string
	if arg := c.Query("title"); arg != "" {
		title = arg
		valid.MaxSize(title, 100, "title").Message("标题最长100个字")
		valid.MinSize(title, 1, "title").Message("文章标题不能为空")
	}
	if arg := c.Query("desc"); arg != "" {
		desc = arg
		valid.MinSize(desc, 20, "desc").Message("文章简述最少20个字符")
		valid.MaxSize(desc, 200, "desc").Message("文章简述最多200个字符")
	}
	if arg := c.Query("detail"); arg != "" {
		detail = arg
		valid.MinSize(detail, 300, "detail").Message("文章详情最少300个字符")
		valid.MaxSize(detail, 65535, "detail").Message("文章详情最少65535个字符")
	}

	code := e.INVALID_PARAMS
	data := make(map[string]interface{})
	if !valid.HasErrors() {
		if !models.ExistArticleById(articleId) {
			code = e.ERROR_NOT_EXIST_ARTICLE
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})
		}
		if tagId <= 0 && !models.ExistTagByTagId(tagId) {
			code = e.ERROR_NOT_EXIST_TAG
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})
		}
		if title != "" {
			data["title"] = title
		}
		if desc != "" {
			data["desc"] = desc
		}
		if detail != "" {
			data["detail"] = detail
		}
		if state >= 0 {
			data["state"] = state
		}
		data["cover"] = cover
		models.EditArticle(articleId, data)
		code = e.SUCCESS
	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key: %s, err.msg: %s", err.Key, err.Message)
		}
	}

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
