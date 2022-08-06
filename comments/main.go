package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/WinIT23/microservice-communication/comments/configs"
	"github.com/WinIT23/microservice-communication/comments/constants"
	"github.com/WinIT23/microservice-communication/comments/models"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	ctx := context.Background()
	collection := configs.GetCollection(configs.GetMongoClient(), constants.MONGO_COLLECTION)

	cmt := models.Comment{Id: primitive.NewObjectID(), PostId: primitive.NewObjectID(), Text: "Some random comment."}
	if _, err := collection.InsertOne(ctx, cmt); err != nil {
		fmt.Printf("%v", err)
	}
	ps, err := json.Marshal(cmt)
	if err != nil {
		fmt.Printf("%v", err)
	}

	app.Use(cors.New())

	app.Get("/api/comments", func(c *fiber.Ctx) error {
		c.Response().Header.Add("Content-Type", "application/json")
		return c.SendString(string(ps))
	})

	app.Post("/api/comments", func(c *fiber.Ctx) error {
		var cmt models.Comment
		var inf map[string]string
		if err := c.BodyParser(&inf); err != nil {
			return err
		}

		cmt.Id = primitive.NewObjectID()
		cmt.Text = inf["text"]
		cmt.PostId, err = primitive.ObjectIDFromHex(inf["post_id"])
		if err != nil {
			return err
		}

		if _, err := collection.InsertOne(ctx, cmt); err != nil {
			fmt.Printf("%v", err)
			return err
		}
		return c.JSON(cmt)
	})

	app.Listen(":8001")
}
