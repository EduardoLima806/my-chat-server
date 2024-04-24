package v1

import (
	"net/http"

	"github.com/eduardolima806/my-chat-server/internal/controller/http/v1/user_route"
	"github.com/eduardolima806/my-chat-server/internal/usecase/user_usecase"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, userUseCase user_usecase.UserBaseUserCase) {

	handler.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "The server is up and running. Chat Server")
	})

	unversionedGroup := handler.Group("/api/v1")
	{
		user_route.NewUserRoute(unversionedGroup, userUseCase)
	}
}
