package server

import (
	"log"
	"time"

	"github.com/99designs/gqlgen/handler"
	auth0 "github.com/auth0-community/go-auth0"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	easyapi "github.com/gradiented/easyisa/pkg/graphql"
	"gopkg.in/square/go-jose.v2"
)

func graphqlHandler() gin.HandlerFunc {
	c := easyapi.Config{Resolvers: &easyapi.Resolver{}}
	c.Directives.HasRole = easyapi.Auth_Directive

	h := handler.GraphQL(easyapi.NewExecutableSchema(c))
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func authMiddleWare() gin.HandlerFunc {

	client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: "https://easy-isa.auth0.com/.well-known/jwks.json"}, nil)
	audience := "http://localhost:3000/graphql"
	configuration := auth0.NewConfiguration(client, []string{audience}, "https://easy-isa.auth0.com/", jose.RS256)
	validator := auth0.NewValidator(configuration, nil)

	return gin.HandlerFunc(func(c *gin.Context) {

		tok, err := validator.ValidateRequest(c.Request)

		if err != nil {
			// c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			// c.Abort()
			log.Println("Invalid token:", err)
		}

		claims := map[string]interface{}{}
		err = validator.Claims(c.Request, tok, &claims)
		if err != nil {
			// c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			// c.Abort()
			log.Println("Invalid claims:", err)
		}

		c.Next()
	})
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := handler.Playground("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func Start() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{"POST"},
		AllowHeaders:     []string{"Origin Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(static.Serve("/", static.LocalFile("./web/build", true)))
	r.Use(authMiddleWare())
	r.POST("/query", graphqlHandler())
	r.GET("/easyapi", playgroundHandler())

	r.Run("localhost:8080")
}
