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

### 2.3 历史记录查询接口 (全量)

该接口用于获取指定用户的所有拦截/上报记录。

*   **URL**: `/data/get`
*   **Method**: `GET`
*   **鉴权**: 需要 `Authorization` Header

#### 请求参数 (Query)

| 参数名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `telephone` | string | 是 | 要查询的目标手机号（可以是本人，也可以是已关联的亲友号码） |

#### 响应结果 (JSON)

```json
{
  "code": 0,
  "message": "",
  "data": [
    {
      "telephone": "13800138000",
      "text": "短信内容...",
      "package": "10086",
      "type": "中性"
    },
    ...
  ]
}
```

---

### 2.4 历史记录查询接口 (分页)

该接口用于分页获取指定用户的拦截/上报记录，适合数据量较大时使用。

*   **URL**: `/data/cutget`
*   **Method**: `POST`
*   **Content-Type**: `application/x-www-form-urlencoded`
*   **鉴权**: 需要 `Authorization` Header

#### 请求参数 (Body)

| 参数名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `telephone` | string | 是 | 要查询的目标手机号 |
| `page` | int | 是 | 页码（从1开始） |
| `cut` | int | 是 | 每页显示的条数 |

#### 响应结果 (JSON)

```json
{
  "code": 0,
  "message": "",
  "pages": 10, // 总页数
  "data": [ ... ] // 同上
}
```

---

### 2.5 风险统计接口

该接口用于获取当前用户的风险拦截统计数据（不同风险类型的数量）。

*   **URL**: `/data/count`
*   **Method**: `GET`
*   **鉴权**: 需要 `Authorization` Header

#### 响应结果 (JSON)

```json
{
  "code": 0,
  "all": 15, // 总拦截数
  "data": {
    "0": 10, // 中性 (Type ID: 0)
    "1": 2,  // 网络交易及兼职诈骗 (Type ID: 1)
    "2": 1,  // 虚假金融及投资诈骗 (Type ID: 2)
    "3": 2   // 身份冒充及威胁诈骗 (Type ID: 3)
  }
}
```

---

### 2.6 关联账户管理接口

这些接口用于管理亲情/关联账户。建立关联后，当前用户才有权限查询被关联号码的短信记录（即上述 `/data/get` 接口）。

#### 2.6.1 获取关联列表

获取当前用户已添加的所有关联账户信息。

*   **URL**: `/user/link/list`
*   **Method**: `GET`
*   **鉴权**: 需要 `Authorization` Header

#### 请求头 (Headers)

| 参数名 | 值 | 说明 |
| :--- | :--- | :--- |
| Authorization | `Bearer <access_token>` | 登录后获取的 Token |

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
    },
    {
      "id": 2,
      "name": "父亲",
      "targetTelephone": "13900139000"
    }
  ]
}
```

---

#### 2.6.2 添加关联账户

添加一个新的关联账户。添加成功后，系统会自动建立当前用户与目标手机号之间的权限信任关系。

*   **URL**: `/user/link/add`
*   **Method**: `POST`
*   **Content-Type**: `application/json`
*   **鉴权**: 需要 `Authorization` Header

#### 请求头 (Headers)

| 参数名 | 值 | 说明 |
| :--- | :--- | :--- |
| Authorization | `Bearer <access_token>` | 登录后获取的 Token |
| Content-Type | `application/json` | 数据格式声明 |

#### 请求参数 (Body)

| 参数名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `name` | string | 是 | 关联对象的备注名称（如：孩子、父亲） |
| `telephone` | string | 是 | 目标的手机号 |

**示例 Body**:
```json
{
  "name": "孩子",
  "telephone": "13800138000"
}
```

#### 响应结果 (JSON)

**成功响应**

```json
{
  "code": 200,
  "message": "添加成功",
  "data": {
    "id": 3,
    "name": "孩子",
    "targetTelephone": "13800138000"
  }
}
```

**失败响应**

```json
{
  "code": 400,
  "message": "该账号已关联", // 或 "不能关联自己"
  "data": null
}
```

---

#### 2.6.3 移除关联账户

根据 ID 移除关联账户，并解除权限绑定。

*   **URL**: `/user/link/remove/:id`
*   **Method**: `DELETE`
*   **鉴权**: 需要 `Authorization` Header

#### 请求头 (Headers)

| 参数名 | 值 | 说明 |
| :--- | :--- | :--- |
| Authorization | `Bearer <access_token>` | 登录后获取的 Token |

#### URL 参数 (Params)

| 参数名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `id` | int | 是 | 关联记录的 ID (从列表接口获取) |

#### 响应结果 (JSON)

```json
{
  "code": 200,
  "message": "移除成功",
  "data": null
}
```

---


