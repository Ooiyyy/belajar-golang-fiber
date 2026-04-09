package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "embed"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

var app = fiber.New(fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		ctx.Status(fiber.StatusInternalServerError)
		return ctx.SendString("Error : " + err.Error())
	},
})

func TestRoutingHello(t *testing.T) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("hello world")
	})

	request := httptest.NewRequest("GET", "/", nil)
	resposnse, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, resposnse.StatusCode)

	bytes, err := io.ReadAll(resposnse.Body)
	assert.Nil(t, err)
	assert.Equal(t, "hello world", string(bytes))

}

func TestCtx(t *testing.T) {
	app.Get("/hello", func(ctx *fiber.Ctx) error {
		name := ctx.Query("name", "Guest")
		return ctx.SendString("hello " + name)
	})

	request := httptest.NewRequest("GET", "/hello?name=roy", nil)
	resposnse, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, resposnse.StatusCode)

	bytes, err := io.ReadAll(resposnse.Body)
	assert.Nil(t, err)
	assert.Equal(t, "hello roy", string(bytes))

	app.Get("/hello", func(ctx *fiber.Ctx) error {
		name := ctx.Query("name", "Guest")
		return ctx.SendString("hello " + name)
	})

	request = httptest.NewRequest("GET", "/hello", nil)
	resposnse, err = app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, resposnse.StatusCode)

	bytes, err = io.ReadAll(resposnse.Body)
	assert.Nil(t, err)
	assert.Equal(t, "hello Guest", string(bytes))

}

func TestHttpRequest(t *testing.T) {
	app.Get("/request", func(ctx *fiber.Ctx) error {
		first := ctx.Get("firstname")
		last := ctx.Cookies("lastname")
		return ctx.SendString("hello " + first + " " + last)
	})

	request := httptest.NewRequest("GET", "/request", nil)
	request.Header.Set("firstname", "Roy")
	request.AddCookie(&http.Cookie{Name: "lastname", Value: "Nugroho"})
	resposnse, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, resposnse.StatusCode)

	bytes, err := io.ReadAll(resposnse.Body)
	assert.Nil(t, err)
	assert.Equal(t, "hello Roy Nugroho", string(bytes))

}

func TestRouteParameter(t *testing.T) {
	app.Get("/users/:userId/order/:orderId", func(ctx *fiber.Ctx) error {
		userId := ctx.Params("userId")
		orderId := ctx.Params("orderId")
		return ctx.SendString("Get Order " + orderId + " from User " + userId)
	})

	request := httptest.NewRequest("GET", "/users/1/order/4", nil)
	resposnse, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, resposnse.StatusCode)

	bytes, err := io.ReadAll(resposnse.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Get Order 4 from User 1", string(bytes))

}

func TestFormReq(t *testing.T) {
	app.Post("/hello", func(ctx *fiber.Ctx) error {
		name := ctx.FormValue("name")
		return ctx.SendString("hello " + name)
	})

	body := strings.NewReader("name=Roy")
	request := httptest.NewRequest("POST", "/hello", body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resposnse, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, resposnse.StatusCode)

	bytes, err := io.ReadAll(resposnse.Body)
	assert.Nil(t, err)
	assert.Equal(t, "hello Roy", string(bytes))

}

//go:embed source/contoh.txt
var contohFile []byte

func TestFormUpl(t *testing.T) {
	app.Post("/upload", func(ctx *fiber.Ctx) error {
		file, err := ctx.FormFile("file")
		if err != nil {
			return err
		}

		err = ctx.SaveFile(file, "./target/"+file.Filename)
		if err != nil {
			return err
		}
		return ctx.SendString("upload success")
	})

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	file, err := writer.CreateFormFile("file", "contoh.txt")
	assert.Nil(t, err)
	file.Write(contohFile)
	writer.Close()

	request := httptest.NewRequest("POST", "/upload", body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	resposnse, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, resposnse.StatusCode)

	bytes, err := io.ReadAll(resposnse.Body)
	assert.Nil(t, err)
	assert.Equal(t, "upload success", string(bytes))

}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestBodyReq(t *testing.T) {
	app.Post("/login", func(ctx *fiber.Ctx) error {
		body := ctx.Body()

		request := new(LoginReq)
		err := json.Unmarshal(body, request)
		if err != nil {
			return err
		}
		return ctx.SendString("hello " + request.Username)
	})

	req := strings.NewReader(`{"username":"roy", "password":"123"}`)
	request := httptest.NewRequest("POST", "/login", req)
	request.Header.Set("Content-Type", "application/json")
	resposnse, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, resposnse.StatusCode)

	bytes, err := io.ReadAll(resposnse.Body)
	assert.Nil(t, err)
	assert.Equal(t, "hello roy", string(bytes))

}

type RegisterReq struct {
	Username string `json:"username" xml:"username" form:"username"`
	Password string `json:"password" xml:"password" form:"password"`
	Name     string `json:"name" xml:"name" form:"name"`
}

func TestBodyParserJSON(t *testing.T) {
	app.Post("/register", func(ctx *fiber.Ctx) error {
		request := new(RegisterReq)
		err := ctx.BodyParser(request)
		if err != nil {
			return nil
		}
		return ctx.SendString("Register Success " + request.Username)
	})

	req := strings.NewReader(`{"username":"roy", "password":"123", "name":"Roikhan"}`)
	request := httptest.NewRequest("POST", "/register", req)
	request.Header.Set("Content-Type", "application/json")
	resposnse, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, resposnse.StatusCode)

	bytes, err := io.ReadAll(resposnse.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Register Success roy", string(bytes))

}

func TestBodyParserForm(t *testing.T) {
	app.Post("/register", func(ctx *fiber.Ctx) error {
		request := new(RegisterReq)
		err := ctx.BodyParser(request)
		if err != nil {
			return nil
		}
		return ctx.SendString("Register Success " + request.Username)
	})

	req := strings.NewReader(`username=roy&password=123&name=Roikhan}`)
	request := httptest.NewRequest("POST", "/register", req)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resposnse, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, resposnse.StatusCode)

	bytes, err := io.ReadAll(resposnse.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Register Success roy", string(bytes))

}

func TestBodyParseXml(t *testing.T) {
	app.Post("/register", func(ctx *fiber.Ctx) error {
		request := new(RegisterReq)
		err := ctx.BodyParser(request)
		if err != nil {
			return nil
		}
		return ctx.SendString("Register Success " + request.Username)
	})

	req := strings.NewReader(
		`<RegisterReq>
			<username>roy</username>
			<password>123</password>
			<name>Roikhan</name>
		</RegisterReq>`)
	request := httptest.NewRequest("POST", "/register", req)
	request.Header.Set("Content-Type", "application/xml")
	resposnse, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, resposnse.StatusCode)

	bytes, err := io.ReadAll(resposnse.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Register Success roy", string(bytes))

}

func TestResponseJSON(t *testing.T) {
	app.Get("/user", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"username": "roy",
			"name":     "roikhan",
		})
	})

	request := httptest.NewRequest("GET", "/user", nil)
	request.Header.Set("Content-Type", "application/json")
	resposnse, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, resposnse.StatusCode)

	bytes, err := io.ReadAll(resposnse.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"name":"roikhan","username":"roy"}`, string(bytes))

}

func TestDownloadFile(t *testing.T) {
	app.Get("/download", func(ctx *fiber.Ctx) error {
		return ctx.Download("./source/contoh.txt", "contoh.txt")
	})

	request := httptest.NewRequest("GET", "/download", nil)
	resposnse, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, `attachment; filename="contoh.txt"`, resposnse.Header.Get("Content-Disposition"))
	assert.Equal(t, 200, resposnse.StatusCode)

	bytes, err := io.ReadAll(resposnse.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Download Berhasil!", string(bytes))

}

func TestRouterGroup(t *testing.T) {
	helloWorld := func(ctx *fiber.Ctx) error {
		return ctx.SendString("hello world")
	}

	api := app.Group("/api")
	api.Get("/hello", helloWorld)
	api.Get("/world", helloWorld)

	web := app.Group("/web")
	web.Get("/hello", helloWorld)
	web.Get("/world", helloWorld)

	request := httptest.NewRequest("GET", "/api/hello", nil)
	resposnse, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, resposnse.StatusCode)

	bytes, err := io.ReadAll(resposnse.Body)
	assert.Nil(t, err)
	assert.Equal(t, "hello world", string(bytes))

}

func TestStatic(t *testing.T) {
	app.Static("/public", "./source")

	request := httptest.NewRequest("GET", "/public/contoh.txt", nil)
	resposnse, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, resposnse.StatusCode)

	bytes, err := io.ReadAll(resposnse.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Download Berhasil!", string(bytes))

}
