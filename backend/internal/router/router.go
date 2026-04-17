package router

import (
	"GoWork_9/backend/internal/controller"
	"GoWork_9/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SetupRouter 配置所有路由
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// 使用跨域中间件
	router.Use(middleware.Cors())

	// 1. 托管静态资源 (JS, CSS, Images)
	// 访问路径: http://localhost:8080/static/...
	// 物理路径: ./frontend/static
	router.Static("/static", "./frontend/static")

	// 2. 托管 HTML 页面 (View 目录)
	// 访问路径: http://localhost:8080/view/...
	router.StaticFS("/view", http.Dir("./frontend/view"))

	// 3. 配置根目录重定向或直接访问 (可选)
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/view/admin/go.html")
	})

	v1 := router.Group("/api/v1")
	{
		// ================= 1. 身份认证 (公用) =================
		auth := v1.Group("/auth")
		{
			auth.POST("/login", controller.AuthCtrl.Login)       // 登录
			auth.POST("/register", controller.AuthCtrl.Register) // 注册
		}

		// ================= 2. 前台门户 (Portal) =================
		// 面向读者：只读操作为主，无需管理员权限
		portal := v1.Group("/portal")
		{
			// 文章
			portal.GET("/articles", controller.ArticleCtrl.PortalList)   // 已发布列表(带分页)
			portal.GET("/article/:id", controller.ArticleCtrl.PortalGet) // 文章详情内容

			// 分类
			portal.GET("/categories", controller.CategoryCtrl.AdminList) // 分类导航条

			// 评论
			portal.GET("/comments/:aid", controller.CommentCtrl.GetByArticleID) // 查看文章评论

			// 互动：仅限登录用户
			portalUser := portal.Group("/").Use(middleware.AuthJWT())
			{
				portalUser.POST("/comment/add", controller.CommentCtrl.Create) // 发表评论
				portalUser.GET("/user/me", controller.AuthCtrl.GetMe)          // 获取个人信息
			}
		}

		// ================= 3. 管理后台 (Admin) =================
		// 面向管理：严控权限，支持全量 CRUD
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthJWT(), middleware.IsAdmin())
		{
			// 系统统计
			admin.GET("/dashboard", controller.AdminCtrl.GetStats)
			//图片上传
			admin.POST("/upload", controller.ArticleCtrl.UploadImage)
			// 文章管理
			article := admin.Group("/articles")
			{
				article.POST("/create", controller.ArticleCtrl.Create)       // 发布
				article.POST("/upload", controller.ArticleCtrl.UploadImage)  // 图片上传
				article.GET("/list", controller.ArticleCtrl.AdminList)       // 管理列表(含分页、搜索)
				article.GET("/get/:id", controller.ArticleCtrl.GetByID)      // 用于编辑回显
				article.PUT("/update/:id", controller.ArticleCtrl.Update)    // 修改
				article.DELETE("/delete/:id", controller.ArticleCtrl.Delete) // 删除
			}

			// 分类管理 (补全)
			category := admin.Group("/categories")
			{
				category.GET("/list", controller.CategoryCtrl.AdminList)       // 全量列表
				category.POST("/create", controller.CategoryCtrl.Create)       // 新增分类
				category.PUT("/update/:id", controller.CategoryCtrl.Update)    // 修改分类
				category.DELETE("/delete/:id", controller.CategoryCtrl.Delete) // 删除分类
			}

			// 评论管理 (补全)
			comment := admin.Group("/comments")
			{
				comment.GET("/list", controller.CommentCtrl.AdminList)       // 全站评论审核列表
				comment.PUT("/audit/:id", controller.CommentCtrl.Audit)      // 审核评论(通过/隐藏)
				comment.DELETE("/delete/:id", controller.CommentCtrl.Delete) // 违规物理删除
			}
		}
	}
	return router
}
