package main

import (
	"CrudGO/database"
	"CrudGO/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	defer database.DB.Close()

	r := gin.Default()

	r.LoadHTMLGlob("web/*")
	r.Static("/static", "./static")

	routes.SetupRoutes(r)
	routes.SetupUsuarioRoutes(r)
	routes.SetupPostRoutes(r)
	routes.SetupComentarioRoutes(r)
	routes.SetupLoginRoutes(r)
	routes.SetupRegistroRoutes(r)

	log.Println("http://localhost:8080")
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
