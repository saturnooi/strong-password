package main

import (
	"database/sql"
	"log"
	"strong_password/app/handler"

	"github.com/acoshift/configfile"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	config := configfile.NewEnvReader()
	db, err := sql.Open("postgres", config.MustString("db_url"))
	if err != nil {
		log.Fatalf("can not open db: %v", err)
	}
	defer db.Close()

	port := config.StringDefault("port", "8080")

	r := gin.Default()
	r.POST("/api/strong_password_steps", handler.StrongPasswordSteps(db))
	r.Run(":" + port)
}
