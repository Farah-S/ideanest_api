package main

// ? Require the packages
import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/example/golang-test/config"
	// "github.com/example/golang-test/pkg/api/routes"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// ? Create required variables that we'll re-assign later
var (
	server      *gin.Engine
	ctx         context.Context
	mongoclient *mongo.Client
	redisclient *redis.Client
)

// ? Init function that will run before the "main" function
func init() {

	// ? Load the .env variables
	config, err := config.LoadDBConfig("./config/")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	// ? Create a context
	ctx = context.TODO()

	// ? Connect to MongoDB
	mongoconn := options.Client().ApplyURI(config.DBUri)
	mongoclient, err := mongo.Connect(ctx, mongoconn)

	if err != nil {
		panic(err)
	}

	if err := mongoclient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB successfully connected...")

	// ? Connect to Redis
	redisclient = redis.NewClient(&redis.Options{
		Addr: config.RedisUri,
	})

	if _, err := redisclient.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	err = redisclient.Set(ctx, "test", "Welcome to Golang with Redis and MongoDB", 0).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("Redis client connected successfully...")

	// ? Create the Gin Engine instance
	server = gin.Default()
}

func main() {
	config, err := config.LoadAppConfig(".")

	if err != nil {
		log.Fatal("Could not load config", err)
	}

	defer mongoclient.Disconnect(ctx)

	value, err := redisclient.Get(ctx, "test").Result()

	if err == redis.Nil {
		fmt.Println("key: test does not exist")
	} else if err != nil {
		panic(err)
	}
	print(value)
	// server := gin.Default()
	server.LoadHTMLGlob("pkg/api/pages/*")
	server.GET("/api", func(c *gin.Context) {
	
	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"index.html",
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"title": "Home Page",
		},
	)
	// c.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
	})
	// apirouter := server.Group("/api")
	// apirouter.GET("/healthchecker", func(ctx *gin.Context) {
	// 	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
	// })
	// router.Run()
	server.Run(":" + config.Port)
	// routes.UserRoute(server)
}
