package handler

import (
	"database/sql"
	"net/http"
	strongpassword "strong_password/app/usecase/strong_password"

	"github.com/gin-gonic/gin"
)

func StrongPasswordSteps(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req strongpassword.StrongPasswordReq
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		steps, err := strongpassword.StrongPasswordSteps(c, req, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"num_of_steps": steps.NumOfSteps})
	}
}