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
	"github.com/rs/cors" // Added CORS package
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// Added CORS package
)

// Image define la estructura de un objeto de imagen que se almacenará en MongoDB.
type Image struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Link string             `bson:"link" json:"link"`
}

var client *mongo.Client // Variable global para el cliente de MongoDB.

func init() {
	// Carga las variables de entorno desde '.env'.
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Inicializa la conexión a MongoDB al iniciar el programa.
	mongoURI := os.Getenv("MONGO_URI") // Obtiene la URI de MongoDB de una variable de entorno.
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
	// Create a new CORS middleware handler allowing requests from your frontend
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5000"}, // Replace with your frontend's port
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},              // Add the methods you need
		AllowedHeaders: []string{"Content-Type"},
	})
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	handlerWithCORS := c.Handler(http.DefaultServeMux)
	// Configura el manejador de GraphQL.
	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})

	// Configura el servidor HTTP.
	http.Handle("/graphql", h)
	fmt.Println("Server is running on http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))

	// Set up your GraphQL endpoint with the handler including CORS
	http.Handle("/graphql", handlerWithCORS)

	fmt.Println("Server is running on http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))

}
