package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	v1 "github.com/project/api-gateway/api/handlers/email"
	"github.com/project/api-gateway/api/handlers/models"
	"google.golang.org/protobuf/encoding/protojson"

	pb "github.com/project/api-gateway/genproto/user_service"
	"github.com/project/api-gateway/pkg/etc"
	l "github.com/project/api-gateway/pkg/logger"
)

// CreateUser creates user
// @Summary creates user api
// @Description new user creation
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.User true "user"
// @Success 200 "success"
// @Router /v1/users [post]
func (h *handlerV1) CreateUser(c *gin.Context) {
	var (
		body        models.User
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().Create(ctx, &pb.User{
		Lastname: body.LastName,
		Name:     body.Name,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create user", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetUser gets user by id
// @Summary getting user by ID
// @Description user by ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 "success"
// @Router /v1/users/{id} [get]
func (h *handlerV1) GetUser(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	guid := c.Param("id")
	id, err := strconv.ParseInt(guid, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed parse string to int", l.Error(err))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().GetUser(
		ctx, &pb.GetUserRequest{
			Id: id,
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// Register registers user
// @Summary registre user api
// @Description register api
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.RegisterUserModel true "register"
// @Success 200 "success"
// @Router /v1/register [post]
func (h *handlerV1) Register(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
		body        models.RegisterUserModel
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind JSON", l.Error(err))
		return
	}

	body.Email = strings.TrimSpace(body.Email)
	body.Userame = strings.TrimSpace(body.Userame)
	body.Email = strings.ToLower(body.Email)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	existEmail, err := h.serviceManager.UserService().CheckField(ctx, &pb.CheckFieldRequest{
		Field: "email",
		Value: body.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error checking email uniqeuness", l.Error(err))
		return
	}

	if existEmail.Exists {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		h.log.Error("email already exists", l.Error(err))
		return
	}

	existUsername, err := h.serviceManager.UserService().CheckField(ctx, &pb.CheckFieldRequest{
		Field: "username",
		Value: body.Userame,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error checking username uniqeuness", l.Error(err))
		return
	}

	if existUsername.Exists {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error while checking username existance", l.Error(err))
		return
	}
	// code generation
	code := etc.GenerateCode(6)
	body.Code = code
	fmt.Println(body.Code)
	message := "Subject: Email verification\n Your verification code: " + body.Code
	err = v1.SendMail([]string{body.Email}, []byte(message))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Error:       err,
			Description: "Email is invalid, check it quickly)",
		})
	}
	c.JSON(http.StatusAccepted, models.Error{
		Error:       nil,
		Description: "Accepted successfully",
	})

	byteUser, err := json.Marshal(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error marshalling data to byte")
		return
	}

	err = h.redisStorage.SetWithTTL(body.Email, string(byteUser), 300)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error seting user data to redis")
		return
	}
	c.JSON(http.StatusAccepted, "success")
}

// func (h *handlerV1) Verify(c *gin.Context) {
// 	var userData models.RegisterUserModel
// 	code := c.Query("code")
// 	email := c.Query("email")

// 	codeToInt, err := strconv.ParseInt(code, 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("error parsing code to int", l.Error(err))
// 		return
// 	}

// 	data, err := h.redisStorage.Get(email)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("error while getting data from redis", l.Error(err))
// 		return
// 	}

// 	err = json.Unmarshal([]byte(data), &userData)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("error while unmarshalling data", l.Error(err))
// 		return
// 	}

// 	if codeToInt != userData.Code {
// 		if err != nil {
// 			c.JSON(http.StatusConflict, gin.H{
// 				"error": "Password not matched",
// 			})
// 			h.log.Error("incorrect password", l.Error(err))
// 			return
// 		}
// 	}
// }
