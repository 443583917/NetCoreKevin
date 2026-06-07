# API 文档

## 认证

所有需要认证的 API 都需要在请求头中携带 JWT Token:

```
Authorization: Bearer <your-token>
```

### 获取 Token

```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "userName": "admin",
  "password": "123456"
}
```

响应:
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "expiresIn": 86400
  }
}
```

## 用户管理

### 获取用户列表

```bash
GET /api/v1/user?page=1&pageSize=10
Authorization: Bearer <token>
```

### 获取用户详情

```bash
GET /api/v1/user/:id
Authorization: Bearer <token>
```

### 创建用户

```bash
POST /api/v1/user
Authorization: Bearer <token>
Content-Type: application/json

{
  "userName": "newuser",
  "password": "password123",
  "realName": "新用户",
  "email": "new@example.com"
}
```

### 更新用户

```bash
PUT /api/v1/user/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "realName": "更新后的名字",
  "email": "updated@example.com"
}
```

### 删除用户

```bash
DELETE /api/v1/user/:id
Authorization: Bearer <token>
```

## 统一响应格式

所有 API 返回统一格式:

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {},
  "timestamp": 1234567890
}
```

### 错误码

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |
