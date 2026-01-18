## 后端启动教程

本项目为 Go + Fiber 后端服务，下面是从零启动的完整步骤。

---

### 1. 环境准备

- 操作系统：Windows / macOS / Linux
- Go：建议 1.20 及以上
- MySQL：任意 5.7+ / 8.x 版本
- Redis：本地或远程实例（默认 `localhost:6379`）

> 说明：项目同时支持 SQLite 和 MySQL。**推荐使用 MySQL**，配置简单且无需 CGO。

---

### 2. 获取代码并安装依赖

在项目根目录（与 `go.mod` 同级）执行：

```bash
go mod tidy
```

该命令会自动拉取 Fiber、GORM 等依赖。

---

### 3. 数据库配置

#### 3.1 推荐：使用 MySQL

1. 在 MySQL 中创建一个空数据库，例如：

   ```sql
   CREATE DATABASE antifraud
     CHARACTER SET utf8mb4
     COLLATE utf8mb4_general_ci;
   ```

2. 打开项目根目录下的 `data/setting.json`，修改 `database` 配置为：

   ```json
   "database": {
     "type": "mysql",
     "source": "root:你的密码@tcp(127.0.0.1:3306)/antifraud?charset=utf8mb4&parseTime=True&loc=Local"
   }
   ```

   - `root`：MySQL 用户名
   - `你的密码`：替换为实际的 MySQL 密码
   - `antifraud`：上面创建的数据库名

3. 首次启动时，程序会通过 GORM 自动在该数据库中创建以下表：

   - `user`
   - `data`
   - `audio`
   - `link`
   - `linked_account`

   不需要手动建表。

#### 3.2 可选：使用 SQLite

如果坚持使用 SQLite，需要：

1. 在本机安装 C 编译器（如 MinGW-w64 / MSYS2 / clang 等）。
2. 启动前在终端设置环境变量（示例）：

   ```bash
   set CGO_ENABLED=1
   set CC=路径到你的gcc.exe
   go run ./main.go
   ```

3. 确认 `data/setting.json` 中配置为：

   ```json
   "database": {
     "type": "sqlite",
     "source": "E:\\antiFraud\\backed\\backed\\upCup\\data\\data.db"
   }
   ```

如果看到错误 `go-sqlite3 requires cgo`，说明 CGO 或 C 编译器配置不正确。

---

### 4. Redis 配置

在 `data/setting.json` 中：

```json
"redisport": 6379
```

默认连接 `localhost:6379`。如果 Redis 不在本机，需要在代码中自行扩展为完整地址（当前版本只配置了端口）。

---

### 5. JWT 与端口配置

同样在 `data/setting.json` 中：

```json
"port": "7000",
"aiport": 6666,
"jwt": {
  "secret": "your-secret-key-change-this-in-production",
  "expires_in": 7200,
  "token_type": "Bearer"
}
```

- `port`：HTTP 服务监听端口
- `aiport`：内部 AI 服务使用的端口
- `jwt.secret`：生产环境请务必修改为复杂随机字符串
- `jwt.expires_in`：Token 有效期（秒）

---

### 6. 启动后端服务

在项目根目录执行：

```bash
go run ./main.go
```

启动成功后，终端会输出类似：

```text
Fiber v2.52.5
http://127.0.0.1:7000
(bound on host 0.0.0.0 and port 7000)
```

此时可以通过浏览器访问：

```text
http://127.0.0.1:7000/
```

---

### 7. 基本接口联调说明

#### 7.1 注册

- URL：`POST /api/auth/register`
- 表单参数：
  - `telephone`
  - `password`

#### 7.2 登录

- URL：`POST /api/auth/login`
- 表单参数：
  - `telephone`
  - `password`
- 返回 `data.access_token`，后续接口需在请求头中携带：

```http
Authorization: Bearer <access_token>
Content-Type: application/json
```

#### 7.3 用户与关联账户相关接口

- 登出：`POST /user/logout`
- 修改密码：`PUT /user/password/change`
- 获取关联账户列表：`GET /user/link/list`
- 添加关联账户：`POST /user/link/add`
- 删除关联账户：`DELETE /user/link/remove/{id}`

以上接口均需要在请求头中携带 JWT Token。

---

### 8. 常见错误排查

1. **`go-sqlite3 requires cgo`**
   - 原因：使用 SQLite 但 `CGO_ENABLED=0` 或未安装 C 编译器。
   - 解决：配置 C 编译器并开启 CGO，或切换为 MySQL。

2. **MySQL 连接失败（`dial tcp ...` 或 `access denied for user`）**
   - 检查：
     - `data/setting.json` 中 DSN 是否正确（用户名 / 密码 / 库名 / 端口）。
     - MySQL 是否监听了 `127.0.0.1:3306`。
     - 用户是否有访问指定数据库的权限。

3. **Redis 连接失败**
   - 确认 Redis 已启动且端口无误。

---

如果需要同时启动前端或 AI 服务，可以在此文档基础上继续补充对应的启动说明。

