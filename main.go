package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID        int    `json:"id" bson:"_id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

var collection *mongo.Collection

func main() {
	fmt.Println("Hello world")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	MONGO_URI := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB Atlas")

	collection = client.Database("go-tanstack-todo").Collection("todos")

	app := fiber.New()

	app.Get("/api/todos", getTodos)
	// app.Post("/api/todos", createTodos)
	// app.Patch("/api/todos/:id", updateTodos)
	// app.Delete("/api/todos/:id", deleteTodos)

	port := os.Getenv("PORT")

	if port == "" {
		port = "5000"
	}

	log.Fatal(app.Listen("0.0.0.0:", port))
}

func getTodos(c *fiber.Ctx) error {
	var todos []Todo
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}

	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return err
		}
		todos = append(todos, todo)

	}

	return c.JSON(todos)
}
