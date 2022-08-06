package main

import (
	"context"
	"fmt"

	"github.com/WinIT23/microservice-communication/comments/configs"
	"github.com/WinIT23/microservice-communication/comments/constants"
	"github.com/WinIT23/microservice-communication/comments/models"
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

	app.Get("/api/posts/:id/comments", func(c *fiber.Ctx) error {
		var cmts []models.Comment

		postId, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return err
		}

		cur, err := collection.Find(ctx, bson.M{"postid": postId})
		if err != nil {
			return err
		}
		defer cur.Close(ctx)

		for cur.Next(ctx) {
			var c models.Comment
			if err := cur.Decode(&c); err != nil {
				return err
			}
			cmts = append(cmts, c)
		}

		return c.JSON(cmts)
	})

	app.Post("/api/comments", func(c *fiber.Ctx) error {
		var cmt models.Comment
		var inf map[string]string
		if err := c.BodyParser(&inf); err != nil {
			return err
		}

		cmt.Id = primitive.NewObjectID()
		cmt.Text = inf["text"]
		PostId, err := primitive.ObjectIDFromHex(inf["post_id"])
		if err != nil {
			return err
		}
		cmt.PostId = PostId

		if _, err := collection.InsertOne(ctx, cmt); err != nil {
			fmt.Printf("%v", err)
			return err
		}
		return c.JSON(cmt)
	})

	app.Listen(":8001")
}
