package v1

import (
	"net/http"
	"strconv"

	"github.com/ToDoApp/api/models"
	"github.com/ToDoApp/storage/repo"
	"github.com/gin-gonic/gin"
)

// @Router /notes/{id} [get]
// @Summary Get note by id
// @Description Get note by id
// @Tags Notes
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.Note
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetNote(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp, err := h.storage.Note().GetById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, models.Note{
		Id:          resp.Id,
		UserId:      resp.UserId,
		Title:       resp.Title,
		Description: resp.Description,
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
		DeletedAt:   resp.DeletedAt,
	})
}

// @Security ApiKeyAuth
// @Router /notes [post]
// @Summary Create a Note
// @Description Create a Note
// @Tags Notes
// @Accept json
// @Produce json
// @Param Note body models.CreateNote true "Note"
// @Success 201 {object} models.Note
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) CreateNote(c *gin.Context) {
	var (
		req models.CreateNote
	)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	payload, err := h.GetAuthPayload(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp, err := h.storage.Note().Create(&repo.Note{
		UserId:      payload.UserID,
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Note{
		Id:          resp.Id,
		UserId:      resp.UserId,
		Title:       resp.Title,
		Description: resp.Description,
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
		DeletedAt:   resp.DeletedAt,
	})
}

// @Router /notes [get]
// @Summary Get all Notes
// @Description Get all Notes
// @Tags Notes
// @Accept json
// @Produce json
// @Param filter query models.GetAllNotesParams false "Filter"
// @Success 200 {object} models.GetAllNotesResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetAllNotes(c *gin.Context) {
	req, err := NotesParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	result, err := h.storage.Note().GetAll(&repo.GetAllNotesParams{
		Page:       req.Page,
		Limit:      req.Limit,
		Search:     req.Search,
		SortByDate: req.SortByDate,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, notesResponse(result))
}

func NotesParams(c *gin.Context) (*models.GetAllNotesParams, error) {
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

	return &models.GetAllNotesParams{
		Page:       page,
		Limit:      limit,
		Search:     c.Query("search"),
		SortByDate: sortByDate,
	}, nil
}
func notesResponse(data *repo.GetAllNotesResult) *models.GetAllNotesResponse {
	response := models.GetAllNotesResponse{
		Notes: make([]*models.Note, 0),
		Count: data.Count,
	}

	for _, note := range data.Notes {
		p := parseNoteModel(note)
		response.Notes = append(response.Notes, &p)
	}

	return &response
}

func parseNoteModel(note *repo.Note) models.Note {
	return models.Note{
		Id:          note.Id,
		UserId:      note.UserId,
		Title:       note.Title,
		Description: note.Description,
		CreatedAt:   note.CreatedAt,
		UpdatedAt:   note.UpdatedAt,
		DeletedAt:   note.DeletedAt,
	}
}

// @Security ApiKeyAuth
// @Summary Update a note
// @Description Update a notes
// @Tags Notes
// @Accept json
// @Produce json
// @Param note body models.CreateNote true "note"
// @Param id path int true "ID"
// @Success 200 {object} models.Note
// @Failure 500 {object} models.ErrorResponse
// @Router /notes/{id} [put]
func (h *handlerV1) UpdateNote(ctx *gin.Context) {
	var (
		req models.Note
	)

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
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
	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	req.UserId = payload.UserID
	note, err := h.storage.Note().Update(&repo.Note{
		Id:          req.Id,
		UserId:      req.UserId,
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   req.CreatedAt,
		UpdatedAt:   req.UpdatedAt,
		DeletedAt:   req.DeletedAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, note)
}

// @Security ApiKeyAuth
// @Summary Delete a Note
// @Description Delete a Note
// @Tags Notes
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ErrorResponse
// @Router /notes/{id} [delete]
func (h *handlerV1) DeleteNote(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to convert",
		})
		return
	}

	err = h.storage.Note().Delete(id)
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
