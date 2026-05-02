/*
 Navicat Premium Data Transfer

 Source Server         : ryh
 Source Server Type    : MySQL
 Source Server Version : 80031
 Source Host           : localhost:3306
 Source Schema         : blog

 Target Server Type    : MySQL
 Target Server Version : 80031
 File Encoding         : 65001

 Date: 02/05/2026 11:24:53
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for articles
-- ----------------------------
DROP TABLE IF EXISTS `articles`;
CREATE TABLE `articles`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` bigint(0) UNSIGNED NULL DEFAULT NULL,
  `category_id` bigint(0) UNSIGNED NULL DEFAULT NULL,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `status` tinyint(0) NULL DEFAULT 0,
  `view_count` bigint(0) UNSIGNED NULL DEFAULT 0,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_id`(`user_id`) USING BTREE,
  INDEX `idx_category_id`(`category_id`) USING BTREE,
  INDEX `idx_status`(`status`) USING BTREE,
  INDEX `idx_articles_deleted_at`(`deleted_at`) USING BTREE,
  INDEX `idx_articles_author_id`(`user_id`) USING BTREE,
  INDEX `idx_articles_category_id`(`category_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of articles
-- ----------------------------
INSERT INTO `articles` VALUES (1, 3, 0, '为什么世界是不存在的', '世界是虚假的，人类是一个又一个的虚拟生物，这个是在高维生物的眼中，但是为什么人类的体内有这么多的精密仪器，为什么会保持姓朱，而不是细胞团', 1, 0, '2026-04-14 20:17:14.418', '2026-04-14 20:54:04.343', NULL);
INSERT INTO `articles` VALUES (2, 3, 0, '人类', '回家看了', 1, 0, '2026-04-14 20:19:28.397', '2026-04-14 20:54:40.391', NULL);
INSERT INTO `articles` VALUES (3, 3, 1, '将开放的国家萨拉', '78970', 1, 0, '2026-04-14 20:34:41.220', '2026-04-14 20:53:53.517', NULL);
INSERT INTO `articles` VALUES (4, 3, 6, 'Gin简介', '# Gin简介：\n\nGin是一个非常受欢迎的Golang Web框架，常用来**快速、高效、稳定地搭建 Web 服务、RESTful API 和微服务**，大幅简化 Go Web 开发流程\n\n路由是由HTTP方法+路径规则+处理函数的映射\n\n一个简单的案例\n\n```go\n//使用Gin框架的时候一定要进行导入\nimport \"github.com/gin-gonic/gin\"\n\nfunc Index(context *gin.Context) {\n	context.String(200, \"Hello ryh!\")\n}\n\nfunc main() {\n	// 创建一个默认的路由\n	router := gin.Default()\n\n	// 绑定路由规则和路由函数，访问/index的路由，将有对应的函数去处理\n	router.GET(\"/index\", Index)\n\n	// 启动监听，gin会把web服务运行在本机的0.0.0.0:8090端口上\n	router.Run(\":8090\")\n}\n```\n\n两种启动方式\n\n```go\n//启动方式一\nrouter.Run(\":8090\")\n//启动方式二\nhttp.ListenAndServe(\":8090\",router)\n```\n\n# 响应\n\n## 返回字符串\n\n```go\nrouter.GET(\"/txt\",func(c *gin.Context){\n    c.String(http.StatusOK,\"返回txt\")\n})\n```\n\n## 返回json\n\n```go\nrouter.GET(\"/json\",func(c *gin.Context){\n    c.JSON(http.StatusOK,gin.H{\"message\":\"hey\",\"status\":http.StatusOK})\n})\n\nrouter.GET(\"/moreJSON\",func(c *gin.Context){\n   type Msg struct {\n		Name    string `json:\"user\"`\n		Message string `json:\"message\"`\n		Number  int    `json:\"number\"`\n	}\n    msg := Msg{\"fencing\", \"hey\", 21}\n    c.JSON(http.StatusOK,msg})\n})\n```\n\n## 文件响应\n\n```go\nfunc _html(context *gin.Context) {\n	context.HTML(200, \"index.html\", gin.H{\n		\"username\": \"lzh\",\n	})\n}\n\n//网页请求这个静态目录的前缀，第二个参数是一个目录，注意，前缀不要重复\nrouter.StaticFS(\"/static\",http.Dir(\"static/static\"))\n//配置单个文件，网页请求的路由，文件的路径\nrouter。StaticFile(\"/titian.png\",\"static/titian.png\")\n```\n\n## 重定向\n\n```go\nfunc _redirect(c *gin.Context) {\n    //支持内部和外部的重定向\n	c.Redirect(302, \"http://www.baidu.com\")\n}\n```\n\n301 Moved Permanently\n\n被请求的资源已永久移动到新位置，并且将来任何对此资源的引用都应该使用本响应返回的若干个URI之一。如果可能，拥有链接编辑功能的客户端应当自动把请求的地址修改从服务器反馈回来的地址，除非额外指定，否则这个响应也是可缓存的\n\n302 Found\n\n请求的资源现在临时从不同的URI响应请求，由于这样的重定向是临时的，客户端应当继续像原有地址发送以后的请求，只有在Cache-Control或Expires中进行了指定的情况下，这个响应才是可缓存的\n\n# 请求\n\n## 请求参数\n\n### 查询参数Query\n\n```go\nfunc _query(c *gin.Context) {\n	fmt.Println(c.Query(\"user\"))\n	fmt.Println(c.GetQuery(\"user\"))\n	fmt.Println(c.QueryArray(\"user\")) //拿到多个相同的查询参数\n}\n```\n\n### 动态参数Param\n\n```go\nfunc _param(c *gin.Context){\n    fmt.Println(c.Param(\"user_id\"))\n    fmt.Println(c.Param(\"book_id\"))\n}\n\nrouter.GET(\"/param/:user_id/\",_param)\nrouter.GET(\"/param/:user_di/book_id\",_param)\n```\n\n### 表单PostForm\n\n```go\n// 表单参数\nfunc _form(c *gin.Context) {\n	fmt.Println(c.PostForm(\"name\"))\n	fmt.Println(c.PostFormArray(\"name\"))\n	fmt.Println(c.DefaultPostForm(\"addr\", \"四川省\"))\n	forms, err := c.MultipartForm()\n	fmt.Println(forms, err)\n}\n//用post请求\n```\n\n### 原始参数 GetRawData\n\n```go\nfunc bindJSON(c *gin.Context, obj any) (err error) {\n	body, _ := c.GetRawData()\n	contentType := c.GetHeader(\"content-Type\")\n	switch contentType {\n	case \"application/json\":\n		err = json.Unmarshal(body, &obj)\n		if err != nil {\n			fmt.Println(err.Error())\n			return err\n		}\n	}\n	return nil\n}\n\n// 原始参数\nfunc _raw(c *gin.Context) {\n	type User struct {\n		Name string `json:\"name\"`\n		Age  string `json:\"age\"`\n	}\n	var user User\n	err := bindJSON(c, user)\n	if err != nil {\n		fmt.Println(err)\n	}\n	fmt.Println(user)\n\n	//body, _ := c.GetRawData()\n	//contentType := c.GetHeader(\"Content-Type\")\n	//switch contentType {\n	//case \"application/json\":\n	//	type User struct {\n	//		Name string `json:\"name\"`\n	//		Age  string `json:\"age\"`\n	//	}\n	//	var user User\n	//	err := json.Unmarshal(body, &user)\n	//	if err != nil {\n	//		fmt.Println(err.Error())\n	//	}\n	//	fmt.Println(string(body))\n	//}\n}\n```\n\n## 四大请求方式\n\n`GET` 	`POST`	 `PUT`	`DELETE`\n\nRestful风格指的是网络应用中就是资源定位和资源操作的风格，不是标准也不是协议\n\n(这四个不是全部，还有其他的，但是这四个是经常使用的，要牢记)\n\n**GET**：从服务器取出资源（一项或多项）\n\n**POST**：在服务器新建一个资源\n\n**PUT**:在服务器更新资源（客户端提供完整资源数据）\n\n**PATCH**:在服务器更新资源（客户端提供需要修改的资源数据）\n\n**DELETE**:从服务器删除数据\n\n```go\n//案例\n//GET		/articles		文章列表\n//GET		/articles/:id	文章详情\n//POST	/articles		添加文章\n//PUT	/articles/:id	修改某一篇文章\n//DELETE/articles/:id	删除某一篇文章\n```\n\n## 请求头相关\n\n### 请求头参数获取\n\n```go\nrouter.GET(\"/\", func(c *gin.Context) {\n		//首字母大小写不区分，单词与单词之间用“-”连接\n		//用于获取一个请求头\n		fmt.Println(c.GetHeader(\"User-Agent\"))\n		fmt.Println(c.GetHeader(\"user-agent\"))\n		fmt.Println(c.GetHeader(\"user-Agent\"))\n		//Header是一个普通的map[string][]string\n		//如果使用get的方法，可以不用区分大小写，并且返回第一个value\n		fmt.Println(c.Request.Header.Get(\"User-Agent\"))\n		//如果使用map的取值的方式，请注意大小写问题\n		fmt.Println(c.Request.Header[\"User-Agent\"])\n		//自定义的请求头，用Get方法也是免大小写的\n		fmt.Println(c.Request.Header.Get(\"Token\"))\n		c.JSON(200, gin.H{\"msg\": \"成功\"})\n	})\n```\n\n## 响应头相关\n\n### 设置响应头\n\n```go\n//设置响应头\n	router.GET(\"/res\", func(c *gin.Context) {\n		c.Header(\"Token\", \"fjal;kjd.shjlafjl;.jalskdjf.aj;lsdkfj\")\n		c.Header(\"Content-Type\", \"application/json; charset=utf-8\")\n		c.JSON(0, gin.H{\"data\": \"看看响应头\"})\n	})\n```\n\n\n\n# 参数绑定\n\n## bind 绑定参数\n\ngin中的bind可以很方便的将前段传递来的数据与结构体进行参数绑定一级参数校验\n\n在使用这个功能的时候，需要给结构体加上Tag  `json`  `form` `url` `xml` `yaml` \n\n### Must bind 系列\n\n可以绑定json、query、param、yaml、xml 如果校验不通过会返回错误\n\n函数是bind（） 、BindJSON（）\n\n### Should bind系列\n\n#### ShouldBindJSON\n\n绑定json格式的数据，tag对应json\n\n```go\n\n```\n\n#### ShouldBindQuery\n\n绑定查询参数\n\ntag对应form\n\n```go\n\n```\n\n#### ShouldBindUri\n\n绑定动态参数\n\ntag对应为uri\n\n```go\n\n```\n\n#### ShouldBind\n\n会根据请求头中的content-type去自动绑定\n\nform-data的参数也使用这个，tag用form\n\n默认的tag就是form\n\n绑定form-data、x-www-form-urlencode\n\n```go\n\n```\n\n## bind绑定器\n\n需要使用参数验证功能，需要加binding tag\n\n### 常用的验证器\n\n| 标签     | 说明                                     | 示例                                               |\n| -------- | ---------------------------------------- | -------------------------------------------------- |\n| required | 必填字段，不能为空，并且不能没有这个字段 | binding:”required“                                 |\n| min=值   | 最小长度                                 | binding:”min=5“                                    |\n| max=值   | 最大长度                                 | binding:”max=18“                                   |\n| len=值   | 长度                                     | binding:”len=6“                                    |\n| gte=值   | 大于等于                                 | binding:”gte=10“                                   |\n| gt=值    | 大于                                     | binding:”gt=10“                                    |\n| lte=值   | 小于等于                                 | binding:”lte=10“                                   |\n| lt=值    | 小于                                     | binding:”lt=10“                                    |\n| eq       | 等于                                     | binding:”eq=3“                                     |\n| ne       | 不等于                                   | binding:”ne=12“                                    |\n| eqfield  | 等于其他字段的值                         | Password string`binding:”eqfirld=ConfirmPassword“` |\n| nefield  | 不等于其他字段的值                       | binding:”“                                         |\n| -        | 忽略字段                                 | binging“-“                                         |\n\n### 自定义验证的错误信息\n\n当验证不通过时，会给出错误的信息，但是原始的错误信息不太好，不利于用户查看\n\n只需要给结构体加一个msg的tag\n\n```go\ntype UserInfo Struct{\n    Username string	`json:\"username\" binding:\"required\" msg:\"用户名不能为空\"`\n    Password string `json:\"password\" binding:\"min=3,max=6\" msg:\"密码长度不能小于3大于6\"`\n    Email	 string `json:\"email\" bingding:\"email\" msg:\"邮箱地址不正确\"`\n}\n```\n\n当出现错误时，就可以来获取出错字段上的msg。\n\n- `err:`这个参数为`ShouldBindJSON`返回的错误信息\n- `obj:`这个参数为绑定的结构体\n- 还有一点要注意的是，vaildator这个包要引用v10这个版本的，否则会出错\n\n### 自定义验证器\n\n1. 注册验证器函数\n\n   ```go\n   //\"github.com/go-playground/validator/v10\"	\n   // 注册自定义验证器\"sign\"，用于 JSON 数据绑定的验证\n   if v, ok := binding.Validator.Engine().(*validator.Validate); ok {\n   	v.RegisterValidation(\"sign\", signValid)\n   }\n   ```\n\n2. 编写函数\n\n   ```go\n   func signValid(fl validator.FieldLevel) bool {\n   	// 定义禁止使用的名称列表\n   	var nameList []string = []string{\"lzh\", \"张三\", \"renlei\"}\n   \n   	// 遍历禁止名单，检查字段值是否与名单中的任何名称匹配\n   	for _, nameStr := range nameList {\n   		// 获取当前字段的值并转换为字符串类型\n   		name := fl.Field().Interface().(string)\n   		// 如果字段值与禁止名单中的某个名称相同，则验证失败\n   		if name == nameStr {\n   			return false\n   		}\n   	}\n   	// 所有检查通过，字段值不在禁止名单中\n   	return true\n   }\n   ```\n\n3. 使用\n\n   ```go\n   // 注册 POST 请求处理函数，路径为根路径\"/\"\n   // 该处理函数接收 JSON 格式的用户数据并进行验证\n   router.POST(\"/\", func(c *gin.Context) {\n   		// 声明 User 结构体变量用于接收绑定数据\n   		var user User\n   		// 将请求的 JSON 数据绑定到 user 变量\n   		err := c.ShouldBindJSON(&user)\n   		// 如果数据绑定或验证失败，返回错误信息\n   		if err != nil {\n   			c.JSON(200, gin.H{\"msg\": GetValidMsg(err, &user)})\n   			return\n   		}\n   		// 验证成功，返回用户数据\n   		c.JSON(200, gin.H{\"data\": user})\n   		return\n   	})\n   ```\n\n# 文件上传和下载\n\n## 文件下载\n\n### 单文件\n\n```go\nfunc main(){\n    router := gin.Default()\n    \n    router.POST(\"/upload\",func(c *gin.Context){\n        file,_ := c.FormFile(\"file\")\n        log.Println(file.Filename)\n        \n        dst :=\"./\" + file.Filename\n        //上传文件至指定的完整文件路径\n        c.SaveUploadedFile(file,dst)\n        c.String(http.StatusOK,fmt.Sprintf(\"\'%s\' uploaded!\",file.Filename))\n    })\n    router.Run(\":80\")\n}\n```\n\n### 服务端保存文件的几种方式\n\n#### SaveUploadedFile\n\n```go\nc.SaveUPloadedFile(file,dst) //文件对象 文件路径,注意要从项目根路径开始写\n```\n\n#### Create+Copy\n\nfile.Open的第一个返回值就是我们讲文件对象中的那个文件（只读的），我们可以使用这个去直接读取文件内容\n\n```go\nfile,_:=c.FormFile(\"file\")\nlog.Println(file.Filename)\n//读取文件中的数据，返回文件对象\nfileRead，——：=file，Open()\ndst := \"./\" + file.Filename\n//创建一个文件\nout,err:=os.Create(dst)\nif err != nil{\n    fmt.Println(err)\n}\ndefer out.Close()\n//拷贝文件对象到out中\nio.Copy(out,fileRead)\n```\n\n#### 读取上传的文件\n\n```go\n	file, _ := c.FormFile(\"file\")\n		readerFile, _ := file.Open()\n		writerFile, _ := os.Create(\"./gin_study/uploads/13.png\")\n		defer writerFile.Close()\n		n, _ := io.Copy(writerFile, readerFile)\n		c.JSON(200, gin.H{\"msg\": \"上传成功\"})\n```\n\n### 多文件\n\n## 文件下载\n\n直接响应一个路径下的文件\n\n```go\nc.File(\"uploads/12.png\")\n```\n\n有些响应，比如图片，浏览器就会显示这个图片，而不是下载，所以我们需要使浏览器唤起下载行为\n\n```go\n\n```\n\n注意、文件下载浏览器可能会有缓存，这个要注意一下\n\n解决办法就是添加查询参数\n\n### 前后端模式下的文件下载\n\n如果是前后端模式下，后端就只需要响应一个文件数据\n\n文件名和其他信息就写在请求头中\n\n```go\nc.Header(\"Content-Type\",\"application/octet-stream\")\nc.File(\"uploads/12.png\")\n```\n\n# gin中间件和路由\n\nGin框架允许开发者在处理请求的过程中，加入用户自己的钩子（Hook）函数，这个钩子函数就叫中间件，中间件适合处理一些公开的业务逻辑，比如登录认证、权限校验、数据分项、记录日志、耗时统计等，即比如、如果访问一个网页的话，不管访问什么路径都需要进行登录、此时就需要为所有路径的处理函数进行统一一个中间件\n\nGin中间件必须是一个gin.Handlerfunc 类型\n\n## 单独注册中间件\n\n```go\nfunc indexHandler(c *Context){\n    fmt.Println(\"index---\")\n    c.JSON(http.StatusOK,gin.H{\n        \"msg\":\"index\",\n    })\n}\n//定义一个中间件\n```\n\n![图片](http://localhost:8080/static/images/articles/1776919975740049100_3.jpg)\n\n全局中间件和中间件传参', 1, 0, '2026-04-23 12:53:06.410', '2026-04-23 12:53:06.410', NULL);

-- ----------------------------
-- Table structure for categories
-- ----------------------------
DROP TABLE IF EXISTS `categories`;
CREATE TABLE `categories`  (
  `id` int(0) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `sort` bigint(0) NULL DEFAULT 0,
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_sort`(`sort`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 13 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of categories
-- ----------------------------
INSERT INTO `categories` VALUES (1, 'go', 11, NULL);
INSERT INTO `categories` VALUES (2, 'mysql', 10, NULL);
INSERT INTO `categories` VALUES (3, 'redis', 9, NULL);
INSERT INTO `categories` VALUES (4, 'docker', 8, NULL);
INSERT INTO `categories` VALUES (5, 'web', 7, NULL);
INSERT INTO `categories` VALUES (6, 'gin', 6, NULL);
INSERT INTO `categories` VALUES (7, 'gorm', 5, NULL);
INSERT INTO `categories` VALUES (8, 'grpc', 4, NULL);
INSERT INTO `categories` VALUES (9, 'ai', 3, NULL);
INSERT INTO `categories` VALUES (10, '分布式', 2, NULL);
INSERT INTO `categories` VALUES (11, '云原生', 1, NULL);
INSERT INTO `categories` VALUES (12, '查看历史搭街坊卡拉', 0, 'JFK拉萨大家发了飞机阿拉就是大佬');

-- ----------------------------
-- Table structure for comments
-- ----------------------------
DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT,
  `article_id` bigint(0) UNSIGNED NOT NULL,
  `user_id` bigint(0) UNSIGNED NOT NULL,
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `status` tinyint(0) NULL DEFAULT 0,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_article_id`(`article_id`) USING BTREE,
  INDEX `idx_user_id`(`user_id`) USING BTREE,
  INDEX `idx_comments_article_id`(`article_id`) USING BTREE,
  INDEX `idx_comments_user_id`(`user_id`) USING BTREE,
  INDEX `idx_comments_deleted_at`(`deleted_at`) USING BTREE,
  INDEX `idx_article_status_deleted`(`article_id`, `status`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 11 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of comments
-- ----------------------------
INSERT INTO `comments` VALUES (1, 3, 2001, '这篇文章写得太好了，学到很多！', 0, '2026-04-26 09:24:08.316', NULL);
INSERT INTO `comments` VALUES (2, 2, 3, '这篇文章的观点我不太认同，有几点想和作者探讨一下', 1, '2026-04-26 09:24:08.316', NULL);
INSERT INTO `comments` VALUES (3, 4, 2003, '垃圾文章，完全没营养', 0, '2026-04-26 09:24:08.316', '2026-04-26 09:24:08.316');
INSERT INTO `comments` VALUES (4, 1, 2004, '我对这篇文章的技术细节有几个补充：首先，关于你提到的性能优化方案，我在实际项目中遇到过类似场景，发现用XXX方式实现比YYY方案能提升约30%的响应速度；其次，文中提到的数据库索引设计有个小误区，在高并发场景下可能会导致死锁问题，建议调整一下索引顺序。希望我的补充能帮到更多读者！', 0, '2026-04-26 09:24:08.316', NULL);
INSERT INTO `comments` VALUES (5, 2, 2005, '', 0, '2026-04-26 09:24:08.316', NULL);
INSERT INTO `comments` VALUES (6, 3, 3, '123456千问', 0, '2026-04-26 10:19:39.856', NULL);
INSERT INTO `comments` VALUES (7, 2, 3, '对话框罚款老师电话发客户', 0, '2026-04-26 10:42:40.255', NULL);
INSERT INTO `comments` VALUES (8, 2, 3, '交罚款上帝就发', 0, '2026-04-26 10:42:54.366', NULL);
INSERT INTO `comments` VALUES (9, 2, 3, '到付哈科室领导反馈', 1, '2026-04-26 10:44:18.598', '2026-04-26 16:14:51.922');
INSERT INTO `comments` VALUES (10, 1, 3, '交罚款螺丝钉解放了\n', 0, '2026-05-01 17:05:42.162', NULL);

-- ----------------------------
-- Table structure for images
-- ----------------------------
DROP TABLE IF EXISTS `images`;
CREATE TABLE `images`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT,
  `article_id` bigint(0) UNSIGNED NULL DEFAULT NULL,
  `user_id` bigint(0) UNSIGNED NULL DEFAULT NULL,
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_article_id`(`article_id`) USING BTREE,
  INDEX `idx_url`(`url`) USING BTREE,
  INDEX `idx_images_article_id`(`article_id`) USING BTREE,
  INDEX `idx_images_user_id`(`user_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 34 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of images
-- ----------------------------
INSERT INTO `images` VALUES (1, 0, 3, '/static/images/articles/1776151836951825800_img.jpg', '2026-04-14 15:30:36.954');
INSERT INTO `images` VALUES (2, 0, 3, '/static/images/articles/1776151870558370000_img.jpg', '2026-04-14 15:31:10.561');
INSERT INTO `images` VALUES (3, 0, 3, '/static/images/articles/1776152038548526200_img.jpg', '2026-04-14 15:33:58.551');
INSERT INTO `images` VALUES (4, 0, 3, '/static/images/articles/1776152182270761200_img.jpg', '2026-04-14 15:36:22.273');
INSERT INTO `images` VALUES (5, 0, 3, '/static/images/articles/1776152446863322500_img.jpg', '2026-04-14 15:40:46.866');
INSERT INTO `images` VALUES (6, 0, 3, '/static/images/articles/1776152468486306400_img.jpg', '2026-04-14 15:41:08.489');
INSERT INTO `images` VALUES (7, 0, 3, '/static/images/articles/1776152711374741700_img.jpg', '2026-04-14 15:45:11.376');
INSERT INTO `images` VALUES (8, 0, 3, '/static/images/articles/1776153019090856400_img.jpg', '2026-04-14 15:50:19.093');
INSERT INTO `images` VALUES (9, 0, 3, '/static/images/articles/1776153386988104900_img.jpg', '2026-04-14 15:56:26.990');
INSERT INTO `images` VALUES (10, 0, 3, '/static/images/articles/1776154069016558900_img.jpg', '2026-04-14 16:07:49.019');
INSERT INTO `images` VALUES (11, 0, 3, '/static/images/articles/1776167062083669100_img.jpg', '2026-04-14 19:44:22.086');
INSERT INTO `images` VALUES (12, 0, 3, '/static/images/articles/1776167390442630600_img.jpg', '2026-04-14 19:49:50.445');
INSERT INTO `images` VALUES (13, 0, 3, '/static/images/articles/1776167891015947400_img.jpg', '2026-04-14 19:58:11.022');
INSERT INTO `images` VALUES (14, 0, 3, '/static/images/articles/1776168293720462600_img.jpg', '2026-04-14 20:04:53.722');
INSERT INTO `images` VALUES (15, 0, 3, '/static/images/articles/1776168546628715500_img.jpg', '2026-04-14 20:09:06.631');
INSERT INTO `images` VALUES (16, 0, 3, '/static/images/articles/1776168657217958500_img.jpg', '2026-04-14 20:10:57.219');
INSERT INTO `images` VALUES (17, 0, 3, '/static/images/articles/1776168693681164600_img.jpg', '2026-04-14 20:11:33.683');
INSERT INTO `images` VALUES (18, 0, 3, '/static/images/articles/1776168710943786000_img.jpg', '2026-04-14 20:11:50.945');
INSERT INTO `images` VALUES (19, 0, 3, '/static/images/articles/1776168843973777800_img.jpg', '2026-04-14 20:14:03.975');
INSERT INTO `images` VALUES (20, 0, 3, '/static/images/articles/1776168859814926100_img.jpg', '2026-04-14 20:14:19.817');
INSERT INTO `images` VALUES (21, 0, 3, '/static/images/articles/1776169030997779100_img.jpg', '2026-04-14 20:17:11.000');
INSERT INTO `images` VALUES (22, 0, 3, '/static/images/articles/1776169162103754500_img.jpg', '2026-04-14 20:19:22.106');
INSERT INTO `images` VALUES (23, 0, 3, '/static/images/articles/1776169233987470900_img.jpg', '2026-04-14 20:20:33.989');
INSERT INTO `images` VALUES (24, 0, 3, '/static/images/articles/1776170041748634400_img.webp', '2026-04-14 20:34:01.749');
INSERT INTO `images` VALUES (25, 0, 3, '/static/images/articles/1776170065849912300_img.jpg', '2026-04-14 20:34:25.851');
INSERT INTO `images` VALUES (26, 0, 3, '/static/images/articles/1776904529222658900_3.jpg', '2026-04-23 08:35:29.226');
INSERT INTO `images` VALUES (27, 0, 3, '/static/images/articles/1776904759288757100_3.jpg', '2026-04-23 08:39:19.291');
INSERT INTO `images` VALUES (28, 0, 3, '/static/images/articles/1776905180600898500_3.jpg', '2026-04-23 08:46:20.604');
INSERT INTO `images` VALUES (29, 0, 3, '/static/images/articles/1776905311946221100_3.jpg', '2026-04-23 08:48:31.948');
INSERT INTO `images` VALUES (30, 0, 3, '/static/images/articles/1776905331819471600_3.jpg', '2026-04-23 08:48:51.822');
INSERT INTO `images` VALUES (31, 0, 3, 'http://localhost:8080/static/images/articles/1776917506877133400_3.jpg', '2026-04-23 12:11:46.881');
INSERT INTO `images` VALUES (32, 0, 3, 'http://localhost:8080/static/images/articles/1776919681055606300_3.jpg', '2026-04-23 12:48:01.059');
INSERT INTO `images` VALUES (33, 4, 3, 'http://localhost:8080/static/images/articles/1776919975740049100_3.jpg', '2026-04-23 12:52:55.744');

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `id` int(0) UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `avatar_url` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `role` enum('admin','user') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT 'user',
  `status` tinyint(0) NULL DEFAULT 1 COMMENT '\'用户状态: 0-禁用, 1-正常\'',
  `create_at` datetime(3) NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_users_username`(`username`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of users
-- ----------------------------
INSERT INTO `users` VALUES (1, 'testuser_86749', '$2a$10$uaWVMb/Nwpb9B.FISqN6z.WGnuDFdESOfYZBuxPvwfnB1femoOT1y', '', 'test@example.com', 'user', 0, '2026-05-02 09:41:35.098');
INSERT INTO `users` VALUES (2, 'testuser_99046', '$2a$10$Ah6Dzl/p2lQPOIlZVXXKjeFTCv1/exMApuGJ8f07aHGf9.Mdw3XKe', '', 'test@example.com', 'user', 1, '2026-05-02 09:41:35.098');
INSERT INTO `users` VALUES (3, 'lzh', '$2a$10$JqCBX1x7E.o6jfZB2TyVzuoNfrNvR1zY9lQ381NyRb76Ztb1iBgMe', 'http://localhost:8080/static/images/avatars/1777633497898092900_3.png', '3641525319@qq.com', 'admin', 1, '2026-05-02 09:41:35.098');
INSERT INTO `users` VALUES (4, 'ryh', '$2a$10$MhEJYNVmoqRTruDqgpYrVeVEPPRIN52fHGb7aTTbak6tdbLYboS3W', '', 'lzh123456@2925.com', 'user', 1, '2026-05-02 09:41:35.098');

SET FOREIGN_KEY_CHECKS = 1;
