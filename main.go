package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		IdleTimeout:  time.Second * 5,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		Prefork:      true,
	})

	app.Use(func(ctx *fiber.Ctx) error {
		fmt.Println("Middleware before processing")
		err := ctx.Next()
		fmt.Println("Middleware after processing")
		return err
	})

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello world!")
	})
	if fiber.IsChild() {
		fmt.Println("I'm child process")
	} else {
		fmt.Println("I'm Parent process")
	}

	err := app.Listen("localhost:8000")
	if err != nil {
		panic(err)
	}
}
