package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/allangrds/fullcycle-mba-go-expert/aulas/7-apis/internal/dto"
	"github.com/allangrds/fullcycle-mba-go-expert/aulas/7-apis/internal/entity"
	"github.com/allangrds/fullcycle-mba-go-expert/aulas/7-apis/internal/infra/database"
	"github.com/go-chi/jwtauth"
)

type Error struct {
	Message string `json:"message"`
}

type UserHandler struct {
	UserDB       database.UserInterface
	JWT          *jwtauth.JWTAuth
	JWTExpiresIn int
}

func NewUserHandler(db database.UserInterface, jwt *jwtauth.JWTAuth, jwtExpiresIn int) *UserHandler {
	return &UserHandler{
		UserDB:       db,
		JWT:          jwt,
		JWTExpiresIn: jwtExpiresIn,
	}
}

// GetJWT godoc
// @Summary Get a user JWT
// @Description Get a user JWT
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.GetJWTInput true "user credentials"
// @Success 200 {object} dto.GetJWTOutput
// @Failure 400 {object} Error
// @Failure 401 {object} Error
// @Router /users/generate-token [post]
func (userHandler *UserHandler) GetJWT(writer http.ResponseWriter, request *http.Request) {
	var userJwtDto dto.GetJWTInput
	err := json.NewDecoder(request.Body).Decode(&userJwtDto)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)

		error := Error{Message: err.Error()}
		json.NewEncoder(writer).Encode(error)

		return
	}

	// Buscar o user pelo email
	user, err := userHandler.UserDB.FindByEmail(userJwtDto.Email)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)

		error := Error{Message: err.Error()}
		json.NewEncoder(writer).Encode(error)

		return
	}

	// Validar a senha
	if !user.IsPasswordValid(userJwtDto.Password) {
		writer.WriteHeader(http.StatusUnauthorized)

		error := Error{Message: err.Error()}
		json.NewEncoder(writer).Encode(error)

		return
	}

	_, tokenString, _ := userHandler.JWT.Encode(map[string]interface{}{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(userHandler.JWTExpiresIn)).Unix(),
	})

	accessToken := dto.GetJWTOutput{AccessToken: tokenString}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(accessToken)
}

// CreateUser godoc
// @Summary Create user
// @Description Create user
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.CreateUserInput true "user request"
// @Success 201
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /users [post]
func (userHandler *UserHandler) CreateUser(writer http.ResponseWriter, request *http.Request) {
	// Lê o body da request
	// Converte o body para um struct dto
	var userDto dto.CreateUserInput
	err := json.NewDecoder(request.Body).Decode(&userDto)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)

		error := Error{Message: err.Error()}
		json.NewEncoder(writer).Encode(error)

		return
	}

	// Fazer diretamente acesso da entidade no coração não é comum. Em vez disso, usaremos no futuro um use case(clean arch)
	user, err := entity.NewUser(userDto.Name, userDto.Email, userDto.Password)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)

		error := Error{Message: err.Error()}
		json.NewEncoder(writer).Encode(error)

		return
	}

	// Persistir o user
	err = userHandler.UserDB.Create(user)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)

		error := Error{Message: err.Error()}
		json.NewEncoder(writer).Encode(error)

		return
	}

	// Retornar o user criado
	writer.WriteHeader(http.StatusCreated)
}
