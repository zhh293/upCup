# 后端完整接口文档 (API Reference)

本文档包含了系统所有的后端接口定义、请求方式、参数说明及响应示例。

**Base URL**: `http://localhost:7000` (默认端口，可配置)

---

## 目录

1.  [系统状态接口](#1-系统状态接口)
2.  [用户认证接口](#2-用户认证接口)
3.  [反诈音频分析接口](#3-反诈音频分析接口)
4.  [关联账户管理接口](#4-关联账户管理接口)
5.  [短信拦截与AI数据接口](#5-短信拦截与ai数据接口)

---

## 1. 系统状态接口

### 1.1 服务健康检查

用于检查服务器是否正常运行。

*   **URL**: `/`
*   **Method**: `GET`
*   **鉴权**: 无需鉴权

#### 响应结果 (JSON)

```json
{
  "code": 0,
  "message": "running...",
  "time": "2024-03-20 10:00:00.00"
}
```

---

## 2. 用户认证接口

### 2.1 用户注册

注册新用户。

*   **URL**: `/api/auth/register`
*   **Method**: `POST`
*   **Content-Type**: `application/x-www-form-urlencoded` 或 `multipart/form-data`

#### 请求参数 (Form Data)

| 参数名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `telephone` | string | 是 | 手机号 (作为唯一用户标识) |
| `password` | string | 是 | 密码 |

#### 响应结果 (JSON)

```json
{
  "code": "200",
  "msg": "注册成功",
  "data": {
    "user_id": "user_13800138000",
    "telephone": "13800138000",
    "access_token": "eyJhbGciOiJIUzI1Ni...",
    "token_type": "Bearer",
    "expires_in": 7200
  }
}
```

### 2.2 用户登录

用户登录并获取 Access Token。

*   **URL**: `/api/auth/login`
*   **Method**: `POST`
*   **Content-Type**: `application/x-www-form-urlencoded` 或 `multipart/form-data`

#### 请求参数 (Form Data)

| 参数名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `telephone` | string | 是 | 手机号 |
| `password` | string | 是 | 密码 |

#### 响应结果 (JSON)

```json
{
  "code": "200",
  "msg": "登录成功",
  "data": {
    "user_id": "user_13800138000",
    "telephone": "13800138000",
    "nickname": "",
    "avatar": "",
    "email": "",
    "access_token": "eyJhbGciOiJIUzI1Ni...",
    "token_type": "Bearer",
    "expires_in": 7200
  }
}
```

### 2.3 用户登出

*   **URL**: `/user/logout`
*   **Method**: `POST`
*   **鉴权**: 需要 Header `Authorization: Bearer <token>`

#### 响应结果 (JSON)

```json
{
  "code": 200,
  "message": "登出成功",
  "data": null
}
```

### 2.4 修改密码

*   **URL**: `/user/password/change`
*   **Method**: `PUT`
*   **Content-Type**: `application/json`
*   **鉴权**: 需要 Header `Authorization: Bearer <token>`

#### 请求参数 (Body)

```json
{
  "currentPassword": "oldpassword123",
  "newPassword": "newpassword456"
}
```

#### 响应结果 (JSON)

```json
{
  "code": 200,
  "message": "密码修改成功，请重新登录",
  "data": null
}
```

---

## 3. 反诈音频分析接口

### 3.1 上传通话音频

上传通话录音文件到服务器。

*   **URL**: `/api/auth/anti-scam/upload-audio`
*   **Method**: `POST`
*   **Content-Type**: `multipart/form-data`
*   **鉴权**: 需要 Header `Authorization: Bearer <token>`

#### 请求参数 (Form Data)

| 参数名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `audio` | File | 是 | 音频文件 (支持 .mp3, .wav, .m4a, .aac, .flac; 最大 50MB) |

#### 响应结果 (JSON)

```json
{
  "code": "200",
  "msg": "上传成功",
  "data": {
    "audioId": "audio_1710900000000",
    "fileName": "call_record.mp3",
    "fileSize": 1024000,
    "uploadTime": "2024-03-20T10:00:00Z",
    "duration": 0,
    "format": "mp3",
    "audioUrl": "/static/audio/audio_1710900000000.mp3"
  }
}
```

### 3.2 分析通话风险

调用 AI 模型对已上传的音频进行风险分析。

*   **URL**: `/api/auth/anti-scam/analyze`
*   **Method**: `POST`
*   **Content-Type**: `application/json`
*   **鉴权**: 需要 Header `Authorization: Bearer <token>`

#### 请求参数 (Body)

```json
{
  "audioId": "audio_1710900000000" // 上传接口返回的 audioId
}
```

#### 响应结果 (JSON)

```json
{
  "code": "200",
  "msg": "分析成功",
  "data": {
    "analysisId": "analysis_1710900100000",
    "riskLevel": "high", // high, medium, low
    "riskScore": 85,     // 0-100
    "summary": "该通话疑似冒充公检法诈骗...",
    "riskTags": ["冒充公检法", "资金核查"],
    "suggestions": ["立即挂断", "拨打96110"],
    "confidence": 0.95
  }
}
```

---

## 4. 关联账户管理接口

### 4.1 获取关联列表

*   **URL**: `/user/link/list`
*   **Method**: `GET`
*   **鉴权**: 需要 Header `Authorization: Bearer <token>`

#### 响应结果 (JSON)

```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "孩子",
      "targetTelephone": "13800138000"
    }
  ]
}
```

### 4.2 添加关联账户

*   **URL**: `/user/link/add`
*   **Method**: `POST`
*   **Content-Type**: `application/json`
*   **鉴权**: 需要 Header `Authorization: Bearer <token>`

#### 请求参数 (Body)

```json
{
  "name": "孩子",
  "telephone": "13800138000"
}
```

### 4.3 移除关联账户

*   **URL**: `/user/link/remove/:id`
*   **Method**: `DELETE`
*   **鉴权**: 需要 Header `Authorization: Bearer <token>`
*   **URL Params**: `id` (关联记录ID)

### 4.4 获取关联账户短信记录

*   **URL**: `/user/link/sms`
*   **Method**: `GET`
*   **鉴权**: 需要 Header `Authorization: Bearer <token>`
*   **Query Params**: `telephone` (目标手机号)

#### 响应结果 (JSON)

```json
{
  "code": 0,
  "message": "",
  "data": [
    {
      "telephone": "13800138000",
      "text": "诈骗短信内容...",
      "package": "Unknown",
      "type": "网络交易及兼职诈骗"
    }
  ]
}
```

---

## 5. 短信拦截与AI数据接口

### 5.1 AI 文本分析

*   **URL**: `/ai/run`
*   **Method**: `POST`
*   **Content-Type**: `application/json`
*   **鉴权**: 需要 Header `Authorization: Bearer <token>`

#### 请求参数 (Body)

```json
{
  "text": "刷单返利，日赚三百..."
}
```

#### 响应结果 (JSON)

```json
{
  "code": 0,
  "message": "",
  "type": "网络交易及兼职诈骗",
  "type_id": 1
}
```

### 5.2 上报拦截数据

*   **URL**: `/data/add`
*   **Method**: `POST`
*   **Content-Type**: `application/json`
*   **鉴权**: 需要 Header `Authorization: Bearer <token>`

#### 请求参数 (Body)

```json
{
  "telephone": "13800138000",
  "text": "短信内容",
  "type": "诈骗",
  "package": "com.android.mms"
}
```

### 5.3 获取所有拦截数据

*   **URL**: `/data/get`
*   **Method**: `GET`
*   **鉴权**: 需要 Header `Authorization: Bearer <token>`
*   **Query Params**: `telephone` (目标手机号)

### 5.4 分页获取拦截数据

*   **URL**: `/data/cutget`
*   **Method**: `POST`
*   **Content-Type**: `application/json` 或 `Form Data`
*   **鉴权**: 需要 Header `Authorization: Bearer <token>`

#### 请求参数 (Body)

```json
{
  "telephone": "13800138000",
  "page": 1,
  "cut": 10
}
```

### 5.5 获取风险统计

*   **URL**: `/data/count`
*   **Method**: `GET`
*   **鉴权**: 需要 Header `Authorization: Bearer <token>`

#### 响应结果 (JSON)

```json
{
  "code": 0,
  "all": 10,
  "data": {
    "0": 5,
    "1": 3,
    "2": 2
  }
}
```
