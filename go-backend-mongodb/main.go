package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Image define la estructura de un objeto de imagen que se almacenará en MongoDB.
type Image struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Link string             `bson:"link" json:"link"`
}

var client *mongo.Client // Variable global para el cliente de MongoDB.

func init() {
	// Attempt to load .env file, but don't fail if it's not found
	_ = godotenv.Load()

	// Initialize MongoDB connection
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable is not set")
	}

	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
}

func ImagesResolver(p graphql.ResolveParams) (interface{}, error) {
	// Utiliza la variable global 'client' para acceder a MongoDB.
	collection := client.Database("oso").Collection("images")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var images []Image
	for cursor.Next(context.Background()) {
		var img Image
		err := cursor.Decode(&img)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, nil
}

func main() {
	// Define los tipos de GraphQL.
	imageType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Image",
		Fields: graphql.Fields{
			"ID": &graphql.Field{
				Type: graphql.String,
			},
			"Link": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	fields := graphql.Fields{
		"images": &graphql.Field{
			Type:    graphql.NewList(imageType),
			Resolve: ImagesResolver,
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// Configura el middleware de CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4000"}, // Cambia esto según sea necesario para tu frontend
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	})

	// Configura el manejador de GraphQL.
	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})

	// Manejador con CORS
	http.Handle("/graphql", c.Handler(h))

	fmt.Println("Server is running on http://0.0.0.0:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
