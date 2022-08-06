package main

import (
	"context"
	"fmt"

	"github.com/WinIT23/microservice-communication/posts/configs"
	"github.com/WinIT23/microservice-communication/posts/constants"
	"github.com/WinIT23/microservice-communication/posts/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	ctx := context.Background()
	collection := configs.GetCollection(configs.GetMongoClient(), constants.MONGO_COLLECTION)

	app.Use(cors.New())

	app.Get("/api/posts", func(c *fiber.Ctx) error {
		posts := []models.Post{}
		cur, err := collection.Find(ctx, bson.M{})
		if err != nil {
			return err
		}
		defer cur.Close(ctx)

		for cur.Next(ctx) {
			var post models.Post
			if err := cur.Decode(&post); err != nil {
				fmt.Printf("%v", err)
			}
			posts = append(posts, post)
		}
		c.Response().Header.Add("Content-Type", "application/json")
		return c.JSON(posts)
	})

	app.Post("/api/posts", func(c *fiber.Ctx) error {
		var post models.Post
		if err := c.BodyParser(&post); err != nil {
			return err
		}

		post.Id = primitive.NewObjectID()
		if _, err := collection.InsertOne(ctx, post); err != nil {
			return err
		}
		return c.JSON(post)
	})
	app.Listen(":8000")
}
