package api

import (
	"github.com/ericlp/songbook/backend/internal/api/endpoints"
	"github.com/ericlp/songbook/backend/internal/api/endpoints/authentication"
	"github.com/ericlp/songbook/backend/internal/common"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var router *gin.Engine

func Init() {
	log.Println("Initializing GIN api")
	router = gin.Default()
	authentication.Init()

	envVars := common.GetEnvVars()
	store := cookie.NewStore([]byte(envVars.Secret))
	store.Options(
		sessions.Options{
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
		},
	)
	router.Use(sessions.Sessions("authentication", store))

	api := router.Group("/api")
	{
		api.Static("/images", envVars.ImageFolder)

		api.GET("/health", endpoints.HealthCheck)
		api.GET("/recipes/:uniqueName", endpoints.Recipe)
		api.GET("/recipes", endpoints.Recipes)
		api.GET("/books/:uniqueName", endpoints.SongBook)
		api.GET("/books", endpoints.SongBooks)
		api.GET("/tags", endpoints.Tags)
		api.GET("/authors", endpoints.Authors)

		authRequired := api.Group("")
		{
			authRequired.Use(authentication.CheckAuth())

			authRequired.GET("/me", authentication.Me)
			authRequired.POST("/recipes", endpoints.NewRecipe)
			authRequired.PUT("/recipes/:id", endpoints.EditRecipe)
			authRequired.DELETE("/recipes/:id", endpoints.RemoveRecipe)
			authRequired.PUT("/images", endpoints.ImageUpload)
			authRequired.POST("/books", endpoints.NewSongBook)
			authRequired.DELETE("/books/:id", endpoints.RemoveSongBook)
			authRequired.PUT("/books/:id", endpoints.EditSongBook)
			authRequired.POST("/tags", endpoints.NewTag)
			authRequired.DELETE("/tags/:id", endpoints.RemoveTag)
			authRequired.PUT("/tags/:id", endpoints.EditTag)
		}

		auth := api.Group("/auth")
		{
			auth.POST("/logout", authentication.Logout)

			gamma := auth.Group("/account")
			{
				gamma.GET("", authentication.GammaInitAuth)
				gamma.GET("/callback", authentication.GammaCallback)
			}
		}
	}
}

func Start() {
	err := router.Run()
	if err != nil {
		log.Fatalf("Failed to start webserver due to err: %s\n", err)
	}
}
