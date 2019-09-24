package main

import (
	"context"
	"time"

	"github.com/99designs/gqlgen/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	graphql "github.com/gradiented/easyisa/graphql"
)

func graphqlHandler() gin.HandlerFunc {
	c := graphql.Config{Resolvers: &graphql.Resolver{}}
	c.Directives.HasRole = graphql.Auth_Directive

	h := handler.GraphQL(graphql.NewExecutableSchema(c))
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func authMiddleWare(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	ctx := context.WithValue(c.Request.Context(), "token", token)
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := handler.Playground("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(authMiddleWare)
	r.POST("/query", graphqlHandler())
	r.GET("/", playgroundHandler())
	r.Run()
}
