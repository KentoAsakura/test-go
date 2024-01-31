package main

import (
	"net/http"
	"log"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"


	"test-go/User"

)




type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	err := t.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Printf("Error rendering template %s: %v", name, err)
	}
	return err
}



func main() {
	User.Init()
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Renderer = &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.GET("/main",func(c echo.Context)error{
		return c.Render(http.StatusOK,"index.html",nil)
	})


	// ログインフォームの表示用のルーティング
	e.GET("/login", func(c echo.Context) error {
		return c.Render(http.StatusOK, "login.html", nil)
	})

	e.Static("/css","css")
	e.Static("/js","js")
	e.Static("/images","images")
	e.Static("/QRCode","QRCode")

	// ログイン用のルーティング
	e.POST("/login", loginHandler)

	// ユーザー登録のルーティング
	e.POST("/register", registerHandler)

	e.Start(":8080")
}

func loginHandler(c echo.Context) error {
	username := c.FormValue("username")
	phoneNumber := c.FormValue("phoneNumber")

	var user User.User
	if err := User.DB.Where("user_name = ?", username).Preload("UserInfo").First(&user).Error; err != nil {
		return c.Render(http.StatusUnauthorized, "login.html", map[string]interface{}{
			"Error": "Invalid credentials",
		})
	}
	// パスワードのハッシュ化はセキュリティ上の理由から必要です
	// ここでは簡単な例として平文のパスワードをそのまま比較します
	if user.PhoneNumber == phoneNumber {
		return c.Render(http.StatusOK, "login_success.html", map[string]interface{}{
			"Username": username,
			"img":user.UserInfo.QRCODE_Number,
		})
	}

	return c.Render(http.StatusUnauthorized, "login.html", map[string]interface{}{
		"Error": "Invalid username or phone number",
	})
}

func registerHandler(c echo.Context) error {
    username := c.FormValue("username")
    password := c.FormValue("phoneNumber")

    if err := User.CreateUser(username, password); err != nil {
        return c.String(http.StatusInternalServerError, "Failed to create user")
    }

    return c.String(http.StatusOK, "User registered successfully")
}


