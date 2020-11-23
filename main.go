package main

import (
	"context"
	"errors"
	"log"
	"session-example/db"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var (
	v     *validator.Validate
	conn  *db.PrismaClient
	store *session.Store
)

// Middleware to ensure the requestor is logged in
func authMW(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return err
	}

	uid := sess.Get("userID")
	if uid == nil {
		return errors.New("not logged in")
	}
	c.Locals("userID", uid)

	return c.Next()
}

func initDB() *db.PrismaClient {

	client := db.NewClient()
	err := client.Connect()

	if err != nil {
		panic(err)
	}

	return client
}

func main() {
	v = validator.New()
	conn = initDB()

	sessStorage := sqlite3.New()
	store = session.New(session.Config{Storage: sessStorage})

	app := fiber.New()

	app.Post("/signup", func(c *fiber.Ctx) error {
		req := struct {
			Username string `validate:"required"`
			Password string `validate:"required"`
		}{}
		// Decode body
		if err := c.BodyParser(&req); err != nil {
			return err
		}
		// Validate body
		if err := v.Struct(&req); err != nil {
			return err
		}

		// Hash password
		// You may want to use argon2 here. For convenience I am using bcrypt.
		// While bcrypt has not been broken yet, argon2 is considered more future-proof.
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		// Create user
		_, err = conn.User.CreateOne(db.User.Username.Set(req.Username), db.User.Password.Set(string(hash))).Exec(context.TODO())
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusCreated)
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		req := struct {
			Username string `validate:"required"`
			Password string `validate:"required"`
		}{}
		// Decode body
		if err := c.BodyParser(&req); err != nil {
			return err
		}
		// Validate body
		if err := v.Struct(&req); err != nil {
			return err
		}

		// Find user
		user, err := conn.User.FindOne(db.User.Username.Equals(req.Username)).Exec(context.TODO())
		if err != nil {
			return err
		}

		// Validate password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.SendString("Invalid password")
		}

		// From here on, the password is verified to be correct
		// Create session
		sess, err := store.Get(c)
		if err != nil {
			return err
		}
		sess.Set("userID", user.ID)
		sess.Save()

		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/protected", authMW, func(c *fiber.Ctx) error {
		return c.SendString(c.Locals("userID").(string))
	})

	log.Fatal(app.Listen(":3000"))
}
