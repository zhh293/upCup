package test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dingdinglz/ai-swindle-detecter-backend/database"
	"github.com/dingdinglz/ai-swindle-detecter-backend/server"
	"github.com/dingdinglz/ai-swindle-detecter-backend/setting"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestApp() *fiber.App {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&database.UserTable{})
	database.MainDB = db

	setting.SettingVar = setting.SettingModel{
		JWT: setting.JWTSettingModel{
			Secret:    "test-secret-for-unit-test",
			ExpiresIn: 3600,
			TokenType: "Bearer",
		},
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"code": code,
				"msg":  err.Error(),
			})
		},
	})

	app.Post("/register", server.UserRegisterRoute)
	app.Post("/login", server.UserLoginRoute)

	return app
}

func createFormData(telephone string, password string) (io.Reader, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("telephone", telephone)
	_ = writer.WriteField("password", password)
	_ = writer.Close()
	return body, writer.FormDataContentType()
}

func TestUserRegisterRoute(t *testing.T) {
	app := setupTestApp()

	tests := []struct {
		name       string
		telephone  string
		password   string
		wantCode   string
		wantMsg    string
		wantFields []string
	}{
		{
			name:      "successful registration",
			telephone: "13800138000",
			password:  "test123456",
			wantCode:  "200",
			wantMsg:   "注册成功",
			wantFields: []string{"user_id", "telephone", "access_token", "token_type", "expires_in"},
		},
		{
			name:      "duplicate registration",
			telephone: "13800138000",
			password:  "test123456",
			wantCode:  "400",
			wantMsg:   "用户已存在",
		},
		{
			name:      "missing telephone",
			telephone: "",
			password:  "test123456",
			wantCode:  "400",
			wantMsg:   "参数不全",
		},
		{
			name:      "missing password",
			telephone: "13900139000",
			password:  "",
			wantCode:  "400",
			wantMsg:   "参数不全",
		},
		{
			name:      "empty both parameters",
			telephone: "",
			password:  "",
			wantCode:  "400",
			wantMsg:   "参数不全",
		},
		{
			name:      "new user registration",
			telephone: "13900139000",
			password:  "newuser123",
			wantCode:  "200",
			wantMsg:   "注册成功",
			wantFields: []string{"user_id", "telephone", "access_token", "token_type", "expires_in"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, contentType := createFormData(tt.telephone, tt.password)
			req := httptest.NewRequest("POST", "/register", body)
			req.Header.Set("Content-Type", contentType)

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			bodyBytes, _ := io.ReadAll(resp.Body)
			bodyStr := string(bodyBytes)

			assert.Contains(t, bodyStr, tt.wantMsg)
			assert.Contains(t, bodyStr, tt.wantCode)

			if tt.wantFields != nil {
				for _, field := range tt.wantFields {
					assert.Contains(t, bodyStr, field)
				}
			}
		})
	}
}

func TestUserLoginRoute(t *testing.T) {
	app := setupTestApp()

	// 先注册一个用户
	body, contentType := createFormData("13800138000", "test123456")
	req := httptest.NewRequest("POST", "/register", body)
	req.Header.Set("Content-Type", contentType)
	app.Test(req)

	tests := []struct {
		name       string
		telephone  string
		password   string
		wantCode   string
		wantMsg    string
		wantFields []string
	}{
		{
			name:      "successful login",
			telephone: "13800138000",
			password:  "test123456",
			wantCode:  "200",
			wantMsg:   "登录成功",
			wantFields: []string{"user_id", "telephone", "access_token", "token_type", "expires_in", "nickname", "avatar", "email"},
		},
		{
			name:       "non-existing user",
			telephone:  "13900139000",
			password:   "test123456",
			wantCode:   "401",
			wantMsg:    "用户不存在",
		},
		{
			name:       "wrong password",
			telephone:  "13800138000",
			password:   "wrongpassword",
			wantCode:   "401",
			wantMsg:    "密码错误",
		},
		{
			name:       "missing telephone",
			telephone:  "",
			password:   "test123456",
			wantCode:   "400",
			wantMsg:    "参数不全",
		},
		{
			name:       "missing password",
			telephone:  "13800138000",
			password:   "",
			wantCode:   "400",
			wantMsg:    "参数不全",
		},
		{
			name:       "empty both parameters",
			telephone:  "",
			password:   "",
			wantCode:   "400",
			wantMsg:    "参数不全",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, contentType := createFormData(tt.telephone, tt.password)
			req := httptest.NewRequest("POST", "/login", body)
			req.Header.Set("Content-Type", contentType)

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			bodyBytes, _ := io.ReadAll(resp.Body)
			bodyStr := string(bodyBytes)

			assert.Contains(t, bodyStr, tt.wantMsg)
			assert.Contains(t, bodyStr, tt.wantCode)

			if tt.wantFields != nil {
				for _, field := range tt.wantFields {
					assert.Contains(t, bodyStr, field)
				}
			}
		})
	}
}

func TestUserRegisterLoginFlow(t *testing.T) {
	app := setupTestApp()

	telephone := "13800138000"
	password := "test123456"

	// 注册
	body, contentType := createFormData(telephone, password)
	req := httptest.NewRequest("POST", "/register", body)
	req.Header.Set("Content-Type", contentType)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	resp.Body.Close()

	// 登录
	body, contentType = createFormData(telephone, password)
	req = httptest.NewRequest("POST", "/login", body)
	req.Header.Set("Content-Type", contentType)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)
	assert.Contains(t, bodyStr, "登录成功")
	assert.Contains(t, bodyStr, telephone)
	resp.Body.Close()
}

func TestInvalidRequestFormat(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("POST", "/register", strings.NewReader("invalid data"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	resp.Body.Close()
}
