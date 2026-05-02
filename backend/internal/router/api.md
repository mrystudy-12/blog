# API 接口文档

## 基础信息

- **Base URL**: `http://localhost:8080/api/v1`
- **认证方式**: JWT Token (在请求头中携带 `Authorization: Bearer <token>`)
- **响应格式**: JSON

---

## 1. 身份认证模块

### 1.1 用户登录

- **路径**: `POST /auth/login`
- **描述**: 用户登录获取 Token
- **请求体**:
```json
{
    "username": "string (用户名)",
    "password": "string (密码)"
}
```
- **成功响应** (200):
```json
{
    "code": 200,
    "message": "登录成功",
    "token": "string (JWT Token)",
    "data": {
        "id": 1,
        "username": "admin",
        "email": "admin@example.com",
        "role": "admin"
    }
}
```

### 1.2 用户注册

- **路径**: `POST /auth/register`
- **描述**: 新用户注册
- **请求体**:
```json
{
    "username": "string (用户名)",
    "password": "string (密码)",
    "email": "string (邮箱)"
}
```
- **成功响应** (200):
```json
{
    "code": 200,
    "message": "注册成功",
    "data": {
        "id": 1,
        "username": "user",
        "email": "user@example.com"
    }
}
```

---

## 2. 前台门户模块

### 2.1 文章相关

#### 2.1.1 获取文章列表

- **路径**: `GET /portal/articles`
- **描述**: 获取已发布文章列表（带分页）
- **请求参数**:
  - `page`: int (页码，默认1)
  - `page_size`: int (每页数量，默认10)
  - `keyword`: string (搜索关键词，可选)
- **成功响应** (200):
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "list": [...],
        "total": 100,
        "page": 1,
        "page_size": 10
    }
}
```

#### 2.1.2 获取文章详情

- **路径**: `GET /portal/article/:id`
- **描述**: 获取单篇文章详情
- **路径参数**:
  - `id`: int (文章ID)
- **成功响应** (200):
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "id": 1,
        "title": "文章标题",
        "content": "文章内容",
        "category_id": 1,
        "author_id": 1,
        "created_at": "2024-01-01 10:00:00"
    }
}
```

### 2.2 分类相关

#### 2.2.1 获取分类列表

- **路径**: `GET /portal/categories`
- **描述**: 获取所有分类（用于导航）
- **成功响应** (200):
```json
{
    "code": 200,
    "message": "success",
    "data": [
        {"id": 1, "name": "技术", "created_at": "2024-01-01"},
        {"id": 2, "name": "生活", "created_at": "2024-01-02"}
    ]
}
```

### 2.3 评论相关

#### 2.3.1 获取文章评论

- **路径**: `GET /portal/comments/:aid`
- **描述**: 获取指定文章的评论列表
- **路径参数**:
  - `aid`: int (文章ID)
- **请求参数**:
  - `page`: int (页码，默认1)
  - `pageSize`: int (每页数量，默认10)
- **成功响应** (200):
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "list": [...],
        "total": 10
    }
}
```

### 2.4 用户互动（需登录）

#### 2.4.1 发表评论

- **路径**: `POST /portal/comment/add`
- **描述**: 登录用户发表评论
- **请求头**: `Authorization: Bearer <token>`
- **请求体**:
```json
{
    "article_id": 1,
    "content": "评论内容"
}
```
- **成功响应** (200):
```json
{
    "code": 200,
    "message": "评论成功"
}
```

#### 2.4.2 头像上传

- **路径**: `POST /portal/upload`
- **描述**: 上传用户头像
- **请求头**: `Authorization: Bearer <token>`
- **请求体**: `multipart/form-data`
  - `avatar`: file (图片文件)
- **成功响应** (200):
```json
{
    "code": 200,
    "message": "头像上传成功",
    "data": {
        "avatar_url": "/static/uploads/avatar/xxx.jpg"
    }
}
```

#### 2.4.3 获取当前用户信息

- **路径**: `GET /portal/user/me`
- **描述**: 获取当前登录用户信息
- **请求头**: `Authorization: Bearer <token>`
- **成功响应** (200):
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "id": 1,
        "username": "user",
        "email": "user@example.com",
        "avatar": "/static/uploads/avatar/xxx.jpg"
    }
}
```

#### 2.4.4 获取用户资料

- **路径**: `GET /portal/user/profile`
- **描述**: 获取当前用户详细资料
- **请求头**: `Authorization: Bearer <token>`

#### 2.4.5 更新用户资料

- **路径**: `PUT /portal/user/profile`
- **描述**: 更新当前用户资料
- **请求头**: `Authorization: Bearer <token>`
- **请求体**:
```json
{
    "username": "新用户名",
    "email": "new@example.com"
}
```

---

## 3. 管理后台模块（需管理员权限）

### 3.1 仪表盘

#### 3.1.1 获取统计数据

- **路径**: `GET /admin/dashboard`
- **描述**: 获取后台仪表盘统计数据
- **请求头**: `Authorization: Bearer <token>`
- **成功响应** (200):
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "article_count": 100,
        "category_count": 10,
        "comment_count": 500,
        "user_count": 200
    }
}
```

### 3.2 文章管理

#### 3.2.1 创建文章

- **路径**: `POST /admin/articles/create`
- **描述**: 发布新文章
- **请求头**: `Authorization: Bearer <token>`
- **请求体**:
```json
{
    "title": "文章标题",
    "content": "文章内容",
    "category_id": 1,
    "status": 1
}
```

#### 3.2.2 文章图片上传

- **路径**: `POST /admin/articles/upload`
- **描述**: 上传文章图片
- **请求头**: `Authorization: Bearer <token>`
- **请求体**: `multipart/form-data`
  - `image`: file (图片文件)

#### 3.2.3 获取文章管理列表

- **路径**: `GET /admin/articles/list`
- **描述**: 获取文章管理列表（含分页、搜索）
- **请求头**: `Authorization: Bearer <token>`
- **请求参数**:
  - `page`: int (页码)
  - `page_size`: int (每页数量)
  - `keyword`: string (搜索关键词)

#### 3.2.4 获取文章详情（编辑用）

- **路径**: `GET /admin/articles/get/:id`
- **描述**: 获取文章详情用于编辑回显
- **请求头**: `Authorization: Bearer <token>`
- **路径参数**: `id` (文章ID)

#### 3.2.5 更新文章

- **路径**: `PUT /admin/articles/update/:id`
- **描述**: 修改文章内容
- **请求头**: `Authorization: Bearer <token>`
- **路径参数**: `id` (文章ID)
- **请求体**:
```json
{
    "title": "新标题",
    "content": "新内容",
    "category_id": 1,
    "status": 1
}
```

#### 3.2.6 删除文章

- **路径**: `DELETE /admin/articles/delete/:id`
- **描述**: 删除文章
- **请求头**: `Authorization: Bearer <token>`
- **路径参数**: `id` (文章ID)

### 3.3 分类管理

#### 3.3.1 获取分类列表

- **路径**: `GET /admin/categories/list`
- **描述**: 获取所有分类
- **请求头**: `Authorization: Bearer <token>`

#### 3.3.2 创建分类

- **路径**: `POST /admin/categories/create`
- **描述**: 新增分类
- **请求头**: `Authorization: Bearer <token>`
- **请求体**:
```json
{
    "name": "分类名称"
}
```

#### 3.3.3 更新分类

- **路径**: `PUT /admin/categories/update/:id`
- **描述**: 修改分类名称
- **请求头**: `Authorization: Bearer <token>`
- **路径参数**: `id` (分类ID)

#### 3.3.4 删除分类

- **路径**: `DELETE /admin/categories/delete/:id`
- **描述**: 删除分类
- **请求头**: `Authorization: Bearer <token>`
- **路径参数**: `id` (分类ID)

### 3.4 评论管理

#### 3.4.1 获取评论列表

- **路径**: `GET /admin/comments/list`
- **描述**: 获取全站评论审核列表
- **请求头**: `Authorization: Bearer <token>`
- **请求参数**:
  - `page`: int (页码)
  - `pageSize`: int (每页数量)
  - `keyword`: string (搜索关键词)

#### 3.4.2 审核评论

- **路径**: `PUT /admin/comments/audit/:id`
- **描述**: 审核评论（通过/隐藏）
- **请求头**: `Authorization: Bearer <token>`
- **路径参数**: `id` (评论ID)
- **请求参数**: `pass` (true: 通过, false: 屏蔽)

#### 3.4.3 删除评论

- **路径**: `DELETE /admin/comments/delete/:id`
- **描述**: 物理删除违规评论
- **请求头**: `Authorization: Bearer <token>`
- **路径参数**: `id` (评论ID)

### 3.5 用户管理

#### 3.5.1 获取用户列表

- **路径**: `GET /admin/users/list`
- **描述**: 获取用户列表（带分页）
- **请求头**: `Authorization: Bearer <token>`
- **请求参数**:
  - `page`: int (页码)
  - `page_size`: int (每页数量)

#### 3.5.2 更新用户状态

- **路径**: `PUT /admin/users/:id/status`
- **描述**: 更新用户状态（启用/禁用）
- **请求头**: `Authorization: Bearer <token>`
- **路径参数**: `id` (用户ID)
- **请求体**:
```json
{
    "status": 0
}
```
- **status 值**: 0=禁用, 1=正常

---

## 响应码说明

| 响应码 | 含义 |
|--------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 401 | 未授权/登录失效 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 通用响应格式

### 成功响应
```json
{
    "code": 200,
    "message": "success",
    "data": "interface{}"
}
```

### 失败响应
```json
{
    "code": 500,
    "message": "错误描述",
    "data": null
}
```