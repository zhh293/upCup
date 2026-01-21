# 短信拦截与AI分析前后端协作接口文档

本文档详细说明了前端（APP端）在拦截到短信后，如何与后端交互进行AI风险分析以及数据上报的流程。

## 1. 业务流程说明

当APP在用户手机上拦截到新的短信时，应遵循以下流程：

1.  **实时分析**：APP调用 **AI文本分析接口** (`/ai/run`)，将短信内容发送给服务器。服务器调用AI模型判断该短信是否存在诈骗风险，并返回风险类型。
2.  **数据上报**：APP调用 **数据上报接口** (`/data/add`)，将拦截到的短信内容、发送方号码等元数据保存到服务器，用于后续的历史记录查询和取证。
    *   *注：建议APP根据用户设置决定是否自动上报。*

---

## 2. 接口详情

### 2.1 AI 文本分析接口 (喂给AI)

该接口用于将短信文本喂给AI模型，获取风险评估结果。

*   **URL**: `/ai/run`
*   **Method**: `POST`
*   **Content-Type**: `application/json` (推荐) 或 `application/x-www-form-urlencoded`
*   **鉴权**: 需要 `Authorization` Header

#### 请求头 (Headers)

| 参数名 | 值 | 说明 |
| :--- | :--- | :--- |
| Authorization | `Bearer <access_token>` | 登录后获取的 Token |
| Content-Type | `application/json` | 数据格式声明 |

#### 请求参数 (Body)

| 参数名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `text` | string | 是 | 需要分析的短信内容 |

#### 响应结果 (JSON)

**成功响应 (Code: 0)**

```json
{
  "code": 0,
  "message": "",
  "type": "网络交易及兼职诈骗",  // 风险类型描述
  "type_id": 1             // 风险类型ID (0:中性, 1:网络交易..., 2:虚假金融..., 3:身份冒充...)
}
```

**失败响应**

```json
{
  "code": 1,
  "message": "ai error"
}
```

---

### 2.2 拦截数据上报接口

该接口用于将拦截到的短信详细信息保存到云端数据库。

*   **URL**: `/data/add`
*   **Method**: `POST`
*   **Content-Type**: `application/json` (推荐) 或 `application/x-www-form-urlencoded`
*   **鉴权**: 需要 `Authorization` Header

#### 请求头 (Headers)

| 参数名 | 值 | 说明 |
| :--- | :--- | :--- |
| Authorization | `Bearer <access_token>` | 登录后获取的 Token |
| Content-Type | `application/json` | 数据格式声明 |

#### 请求参数 (Body)

| 参数名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `telephone` | string | 是 | 短信所属用户的手机号（通常是当前登录用户的手机号） |
| `text` | string | 是 | 短信内容 |
| `package` | string | 是 | 来源信息（通常填写入短信发送方号码，如 10086；或者是来源APP包名） |
| `type` | string | 是 | 短信类型标签（例如："验证码", "推广", "诈骗" 或 AI返回的类型） |

#### 响应结果 (JSON)

**成功响应 (Code: 0)**

```json
{
  "code": 0,
  "message": ""
}
```

**失败响应**

```json
{
  "code": -1,
  "message": "参数不全"
}
```
或
```json
{
  "code": 1,
  "message": "权限错误" // 当尝试为非关联账户上传数据时
}
```

---

## 3. 前端开发建议

1.  **异步处理**：建议在后台线程中执行上述网络请求，避免阻塞UI主线程。
2.  **Token管理**：确保请求时 Token 有效。如果 Token 过期（返回 401），应引导用户重新登录或自动刷新 Token。
3.  **用户隐私**：在上报数据前，最好向用户展示隐私提示，明确告知短信内容将被上传用于反诈分析。
4.  **关联账户支持**：
    *   如果是监护人模式，APP可能会拦截被监护人的短信。
    *   在调用 `/data/add` 时，`telephone` 参数应填入**被监护人（即短信接收者）**的手机号。
    *   后端会自动校验当前登录用户是否有权限为该手机号上传数据。
