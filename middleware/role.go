package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/irfanseptian/fims-backend/models"
	"github.com/irfanseptian/fims-backend/utils"
)

// RequireRole returns a Gin middleware that checks if the authenticated user
// has one of the required roles.
func RequireRole(roles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		role := models.Role(userRole.(string))

		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}

		utils.ErrorResponse(c, http.StatusForbidden, "Anda tidak memiliki akses untuk resource ini")
		c.Abort()
	}
}
