package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

var db *sqlx.DB

func main() {

	var err error
	db, err = sqlx.Open("mysql", "root:keep1234@tcp(127.0.0.1)/keeplearning")
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	app.Post("/signup", Signup)
	app.Post("/login", Login)
	app.Get("/hello", Hello)

	app.Listen(":8000")
}

func Signup(c *fiber.Ctx) error {
	request := SignupRequest{}
	err := c.BodyParser(&request)
	if err != nil {
		return nil
	}

	if request.Username == "" || request.Password == "" {
		return fiber.ErrUnprocessableEntity
	}

	query := "insert user (username, password) values (?, ?)"
	result, err := db.Exec(query, request.Username, request.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	user := User{
		Id: int(id),
		Username: request.Username,
		Password: request.Password,
	}

	
	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	return nil
}

func Hello(c *fiber.Ctx) error {
	return nil
}

type User struct {
	Id       int    `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
}

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Fiber() {
	app := fiber.New(fiber.Config{
		Prefork: true,
	})

	//Middleware
	app.Use("/hello", func(c *fiber.Ctx) error {
		c.Locals("name", "abc")
		fmt.Println("before")
		err := c.Next()
		fmt.Println("after")
		return err
	})

	app.Use(requestid.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "*",
		AllowHeaders: "*",
	}))

	// TimeZone: "Asia/Bangkok",
	app.Use(logger.New(logger.ConfigDefault))

	app.Get("/hello", func(c *fiber.Ctx) error {
		name := c.Locals("name")
		return c.SendString(fmt.Sprintf("GET Hello %v", name))
	})

	app.Post("/hello", func(c *fiber.Ctx) error {
		return c.SendString("POST Hello")
	})

	//Parameters and Optional Parameters
	// app.Post("/hello/:name/:surname?", func(c *fiber.Ctx) error {
	// 	name := c.Params("name")
	// 	surname := c.Params("surname")
	// 	return c.SendString("POST Hello, " + name + surname)
	// })

	//ParamsInt
	app.Post("/hello/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return fiber.ErrBadRequest
		}
		return c.SendString(fmt.Sprintf("ID: %v", id))
	})

	//Query
	app.Get("/query", func(c *fiber.Ctx) error {
		name := c.Query("name")
		surname := c.Query("surname")
		return c.SendString("name: " + name + " surname: " + surname)
	})

	//Query Parser
	app.Get("/query2", func(c *fiber.Ctx) error {
		person := Person{}
		c.QueryParser(&person)
		return c.JSON(person)
	})

	//Wildcards
	app.Get("/wildcards/*", func(c *fiber.Ctx) error {
		wildcards := c.Params("*")
		return c.SendString(wildcards)
	})

	//Static file
	app.Static("/", "./wwwroot", fiber.Static{
		Index:         "index.html",
		CacheDuration: time.Second * 10,
	})

	//New Error
	app.Get("/error", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusNotFound, "content not found")
	})

	//Group
	v1 := app.Group("/v1", func(c *fiber.Ctx) error {
		c.Set("Version", "v1")
		return c.Next()
	})
	v1.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello v1")
	})

	v2 := app.Group("/v2", func(c *fiber.Ctx) error {
		c.Set("Version", "v2")
		return c.Next()
	})
	v2.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello v2")
	})

	//Mount  // curl "localhost:8000/user/login"
	userApp := fiber.New()
	userApp.Get("/login", func(c *fiber.Ctx) error {
		return c.SendString("Login")
	})

	app.Mount("/user", userApp)

	//Server
	app.Server().MaxConnsPerIP = 1
	app.Get("/server", func(c *fiber.Ctx) error {
		time.Sleep(time.Second * 30)
		return c.SendString("server")
	})

	//Environment  // curl "localhost:8000/env?name=asdf" | jq
	app.Get("/env", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"BaseURL":     c.BaseURL(),
			"Hostname":    c.Hostname(),
			"IP":          c.IP(),
			"IPs":         c.IPs(),
			"OriginalURL": c.OriginalURL(),
			"Path":        c.Path(),
			"Protocol":    c.Protocol(),
			"Subdomain":   c.Subdomains(),
		})
	})

	//Body
	app.Post("/body", func(c *fiber.Ctx) error {
		fmt.Printf("isJSON: %v\n", c.Is("json"))
		fmt.Println(string(c.Body()))

		person := Person{}
		err := c.BodyParser(&person)
		if err != nil {
			return err
		}
		fmt.Println(person)

		return nil
	})
	// curl localhost:8000/body -d '{"name": "J"}' -H content-type:application/json
	// curl localhost:8000/body -d 'id=1&name=J' -H content-type:application/x-www-form-urlencoded

	//Body
	app.Post("/body2", func(c *fiber.Ctx) error {
		fmt.Printf("isJSON: %v\n", c.Is("json"))
		// fmt.Println(string(c.Body()))

		data := map[string]interface{}{}
		err := c.BodyParser(&data)
		if err != nil {
			return err
		}
		fmt.Println(data)

		return nil
	})
	// curl localhost:8000/body2 -d '{"id":1,"name":"Hi"}' -H content-type:application/json

	app.Listen(":8000")
}

type Person struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
