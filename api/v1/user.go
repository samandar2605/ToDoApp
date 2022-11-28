package v1

import (
	"net/http"
	"strconv"

	"github.com/ToDoApp/api/models"
	"github.com/ToDoApp/storage/repo"
	"github.com/gin-gonic/gin"
)

// @Router /users/{id} [get]
// @Summary Get user by id
// @Description Get user by id
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.User
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	resp, err := h.storage.User().GetById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.User{
		Id:          resp.Id,
		FirstName:   resp.FirstName,
		LastName:    resp.LastName,
		PhoneNumber: resp.PhoneNumber,
		Password:    resp.Password,
		Email:       resp.Email,
		ImageUrl:    resp.ImageUrl,
		CreatedAt:   resp.CreatedAt,
	})
}

// @Security ApiKeyAuth
// @Router /users [post]
// @Summary Create a user
// @Description Create a user
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.CreateUser true "user"
// @Success 201 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) CreateUser(c *gin.Context) {
	var (
		req models.CreateUser
	)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	resp, err := h.storage.User().Create(&repo.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Password:    req.Password,
		ImageUrl:    req.ImageUrl,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.User{
		Id:          resp.Id,
		FirstName:   resp.FirstName,
		LastName:    resp.LastName,
		PhoneNumber: resp.PhoneNumber,
		Email:       resp.Email,
		Password:    resp.Password,
		ImageUrl:    resp.ImageUrl,
	})
}

// @Router /users [get]
// @Summary Get all Users
// @Description Get all Users
// @Tags users
// @Accept json
// @Produce json
// @Param filter query models.GetAllUsersParams false "Filter"
// @Success 200 {object} models.GetAllUsersResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetAllUsers(c *gin.Context) {
	req, err := usersParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	result, err := h.storage.User().GetAll(&repo.GetAllUsersParams{
		Page:       req.Page,
		Limit:      req.Limit,
		Search:     req.Search,
		SortByDate: req.SortByDate,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, usersResponse(result))
}

func usersParams(c *gin.Context) (*models.GetAllUsersParams, error) {
	var (
		limit      int = 10
		page       int = 1
		sortByDate string
		err        error
	)

	if c.Query("limit") != "" {
		limit, err = strconv.Atoi(c.Query("limit"))
		if err != nil {
			return nil, err
		}
	}

	if c.Query("page") != "" {
		page, err = strconv.Atoi(c.Query("page"))
		if err != nil {
			return nil, err
		}
	}

	if c.Query("sort_by_date") != "" &&
		(c.Query("sort_by_date") == "desc" || c.Query("sort_by_date") == "asc" || c.Query("sort_by_date") == "none") {
		sortByDate = c.Query("sort_by_date")
	}

	return &models.GetAllUsersParams{
		Page:       page,
		Limit:      limit,
		Search:     c.Query("search"),
		SortByDate: sortByDate,
	}, nil
}
func usersResponse(data *repo.GetAllUsersResult) *models.GetAllUsersResponse {
	response := models.GetAllUsersResponse{
		Users: make([]*models.User, 0),
		Count: data.Count,
	}

	for _, user := range data.Users {
		p := parseuserModel(user)
		response.Users = append(response.Users, &p)
	}

	return &response
}

func parseuserModel(user *repo.User) models.User {
	return models.User{
		Id:          user.Id,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Password:    user.Password,
		ImageUrl:    user.ImageUrl,
		CreatedAt:   user.CreatedAt,
	}
}

// @Security ApiKeyAuth
// @Summary Update a user
// @Description Update a userss
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param user body models.CreateUser true "user"
// @Success 200 {object} models.User
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id} [put]
func (h *handlerV1) UpdateUser(ctx *gin.Context) {
	var (
		req models.User
	)

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	req.Id = id
	user, err := h.storage.User().Update(&repo.User{
		Id:          req.Id,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Password:    req.Password,
		ImageUrl:    req.ImageUrl,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// @Security ApiKeyAuth
// @Summary Delete a User
// @Description Delete a user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id} [delete]
func (h *handlerV1) DeleteUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to convert",
		})
		return
	}

	err = h.storage.User().Delete(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "successful delete method",
	})
}