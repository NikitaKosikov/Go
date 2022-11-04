package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"test/internal/domain"
	"test/internal/service/dto"
	"test/pkg/api"
	"test/pkg/api/auth"
	"test/pkg/api/params"

	"github.com/gin-gonic/gin"
)

const (
	idNameURL           = "id"
	SortByParametersURL = "sort_by"
	usersGroup          = "/users"
	adminGroup          = "/admins"
)

func (h *handler) initUsersRoutes(api *gin.RouterGroup) {

	users := api.Group(usersGroup)
	{

		users.POST("/", h.Create)
		users.GET(auth.RefreshURL, h.RefreshToken)

		admin := users.Group("/").Use(h.tokenManager.VerifyJWTMiddleware(auth.AdminRole))
		{
			admin.GET("/", h.FindAll)
		}

		authencticated := users.Group("/").Use(h.tokenManager.VerifyJWTMiddleware(auth.UserRole, auth.AdminRole))
		{
			authencticated.GET("/:id", h.FindOne)
			authencticated.PUT("/:id", h.Update)
			authencticated.DELETE("/:id", h.Delete)

		}

	}
}

// @Summary Create
// @Tags users
// @Description Create user
// @ID create-user
// @Accept json
// @Produce json
// @Param userDTO body user.CreateUserDTO true "user info"
// @Seccess 200 {integer} integer 1
// @Router /users [post]

func (h *handler) Create(ctx *gin.Context) {

	var userDTO dto.CreateUserDTO
	err := ctx.BindJSON(&userDTO)
	if err != nil {
		newResponse(ctx, http.StatusBadRequest, "failed to bind user and json")
		return
	}

	tokenDTO, err := h.services.Users.Create(ctx, userDTO)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			newResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Header("Access-Token", tokenDTO.AccessToken)
	ctx.Header("Refresh-Token", tokenDTO.RefreshToken)
	ctx.Status(http.StatusCreated)
}

// @Summary Find user
// @Tags user/:id
// @Description  Find user details
// @ID find-user
// @Accept json
// @Produce json
// @Param id body user.CreateUserDTO true "user info"
// @Seccess 200 {integer} integer 1
// @Router /user/:id [get]

func (h *handler) FindOne(ctx *gin.Context) {

	id := ctx.Param(idNameURL)
	user, err := h.services.Users.FindOne(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			newResponse(ctx, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, params.ErrInvalidIdParam) {
			newResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}
		
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	userBytes, err := json.Marshal(user)
	if err != nil {
		newResponse(ctx, http.StatusBadRequest, "failed to marshal user to json")
		return
	}
	ctx.Writer.Write(userBytes)
	ctx.Status(http.StatusOK)
}

// @Summary Find users
// @Tags users
// @Description Find users details
// @ID find-users
// @Accept json
// @Produce json
// @Seccess 200 {integer} integer 1
// @Router /users [get]

func (h *handler) FindAll(ctx *gin.Context) {
	sortBy := ctx.Param(SortByParametersURL)
	filter := ctx.Param(api.FilterByParametersURL)
	limit := ctx.Param(api.LimitByParametersURL)
	offset := ctx.Param(api.OffsetByParametersURL)
	users, err := h.services.Users.FindAll(ctx, limit, offset, filter, sortBy)
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	usersBytes, err := json.Marshal(users)
	if err != nil {
		newResponse(ctx, http.StatusBadRequest, "failed to marshal user to json")
		return
	}
	ctx.Writer.Write(usersBytes)
	ctx.Status(http.StatusOK)
}

// @Summary Create
// @Tags users
// @Description Create user
// @ID create-user
// @Accept json
// @Produce json
// @Param userDTO body user.CreateUserDTO true "user info"
// @Seccess 200 {integer} integer 1
// @Router /users [post]

func (h *handler) Update(ctx *gin.Context) {
	var userDTO dto.UpdateUserDTO
	id := ctx.Param(idNameURL)
	userDTO.Id = id
	err := ctx.BindJSON(&userDTO)
	if err != nil {
		newResponse(ctx, http.StatusBadRequest, "failed to bind user and json")
		return
	}

	err = h.services.Users.Update(ctx, userDTO)
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

// @Summary Create
// @Tags users
// @Description Create user
// @ID create-user
// @Accept json
// @Produce json
// @Param userDTO body user.CreateUserDTO true "user info"
// @Seccess 200 {integer} integer 1
// @Router /users [post]

func (h *handler) Delete(ctx *gin.Context) {
	id := ctx.Param(idNameURL)
	err := h.services.Users.Delete(ctx, id)
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

func (h *handler) RefreshToken(ctx *gin.Context) {
	userId := ctx.Param("id")
	tokenDTO, err := h.services.Users.RefreshUserToken(ctx, userId)
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Header("Access-Token", tokenDTO.AccessToken)
	ctx.Header("Refresh-Token", tokenDTO.RefreshToken)
	ctx.Status(http.StatusOK)
}
