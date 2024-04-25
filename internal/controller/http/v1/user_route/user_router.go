package user_route

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/eduardolima806/my-chat-server/internal/domain"
	"github.com/eduardolima806/my-chat-server/internal/usecase/user_usecase"
	"github.com/gin-gonic/gin"
)

type userRouter struct {
	useCase user_usecase.UserBaseUserCase
}

type createUserBody struct {
	UserName    string `json:"userName" binding:"required"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type loginBody struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewUserRoute(handler *gin.RouterGroup, userUseCase user_usecase.UserBaseUserCase) {
	h := handler.Group("/users")
	r := &userRouter{useCase: userUseCase}

	{
		h.POST("/create-user", r.createUser)
		h.POST("/login", r.loginUser)
	}
}

func (route *userRouter) createUser(ctx *gin.Context) {
	var body createUserBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		// TODO: Include logger interface
		// fmt.Errorf("http - v1 - create a user route")
		fmt.Println("http - v1 - create a user route")
		strErr := strings.ReplaceAll(err.Error(), "\n", "\\n")
		bindErr := domain.CreateError(domain.ErrBadRequest.Error(), strErr)
		bindErrStruct := domain.ErrorCodeResponse(bindErr)
		bindErrStruct.ErrorMessage = fmt.Sprintf("Error to bind user data: %s", bindErr.Error())
		ctx.JSON(domain.GetHttpStatusCode(bindErr), bindErrStruct)
		return
	}

	userOutput, err := route.useCase.CreateUserUseCase.Execute(*body.toUserInput())

	if err != nil {
		ctx.JSON(domain.GetHttpStatusCode(err), domain.ErrorCodeResponse(err))
	} else {
		ctx.JSON(http.StatusOK, userOutput)
	}
}

func (route *userRouter) loginUser(ctx *gin.Context) {
	var body loginBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		// TODO: Include logger interface
		// fmt.Errorf("http - v1 - create a user route")
		fmt.Println("http - v1 - login user route")
		strErr := strings.ReplaceAll(err.Error(), "\n", "\\n")
		bindErr := domain.CreateError(domain.ErrBadRequest.Error(), strErr)
		bindErrStruct := domain.ErrorCodeResponse(bindErr)
		bindErrStruct.ErrorMessage = fmt.Sprintf("Error to bind login data: %s", bindErr.Error())
		ctx.JSON(domain.GetHttpStatusCode(bindErr), bindErrStruct)
		return
	}

	userOutput, err := route.useCase.LoginUserUseCase.Execute(*body.tLoginInput())

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, domain.ErrorCodeResponse(err))
	} else {
		if userOutput.IsSucceed {
			ctx.JSON(http.StatusOK, "login succeed")
		} else {
			err := domain.CreateError(domain.ErrBadRequest.Error(), userOutput.ErrorType.Description)
			ctx.JSON(http.StatusBadRequest, domain.ErrorCodeResponse(err))
		}
	}
}

func (body *createUserBody) toUserInput() *user_usecase.UserInput {
	return &user_usecase.UserInput{
		UserName:    body.UserName,
		DisplayName: body.DisplayName,
		Email:       body.Email,
		Password:    body.Password,
	}
}

func (body *loginBody) tLoginInput() *user_usecase.LoginInput {
	return &user_usecase.LoginInput{
		Login:    body.Login,
		Password: body.Password,
	}
}
