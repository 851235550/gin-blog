package routers

import (
	"gin-blog/src/gin-blog/pkg/setting"
	"gin-blog/src/gin-blog/routers/api"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "test",
		})
	})

	apirouter := r.Group("/api")
	{
		//标签
		apirouter.GET("/tags", api.GetTags)
		apirouter.POST("/tag", api.AddTag)
		apirouter.PUT("/tag/:tag_id", api.EditTag)
		apirouter.DELETE("/tag/:tag_id", api.DeleteTag)

		//文章
		//获取文章列表
		apirouter.GET("/articles", api.GetArticles)
		//获取单个文章
		apirouter.GET("/article/:article_id", api.GetArticle)
		//新建文章
		apirouter.POST("/article", api.AddArticle)
		//更新文章
		apirouter.PUT("/article/:article_id", api.EditArticle)
		//删除文章
		apirouter.DELETE("/article/:article_id", api.DeleteArticle)
	}

	return r
}
