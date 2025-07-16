package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/models"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/services"
)

type PostHandler struct {
	BaseHandler
	PostService *services.PostService
}

func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{
		PostService: postService,
	}
}

func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	var req models.PostRequest
	if err := c.BodyParser(&req); err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request", err)
	}

	post, err := h.PostService.CreatePost(&req)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create post", err)
	}

	return h.SuccessResponse(c, "Post created successfully", post)
}

func (h *PostHandler) GetPostByID(c *fiber.Ctx) error {
	id, err := h.GetIDParam(c)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	post, err := h.PostService.GetPostByID(id)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusNotFound, "Post not found", err)
	}

	return h.SuccessResponse(c, "Post retrieved successfully", post)
}

func (h *PostHandler) GetAllPosts(c *fiber.Ctx) error {
	queryPaginate := h.GetPaginationParams(c)

	posts, err := h.PostService.GetAllPosts(queryPaginate.Limit, queryPaginate.Offset)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve posts", err)
	}

	// Count total posts for pagination
	total, err := h.PostService.Count(&models.Post{})
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to count posts", err)
	}

	meta := h.CalculatePaginationMeta(int(c.QueryInt("page", 1)), queryPaginate.Limit, total)
	return h.PaginatedSuccessResponse(c, "Posts retrieved successfully", posts, meta)
}

func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
	id, err := h.GetIDParam(c)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	var req models.PostRequest
	if err := c.BodyParser(&req); err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request", err)
	}

	post, err := h.PostService.UpdatePost(id, &req)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update post", err)
	}

	return h.SuccessResponse(c, "Post updated successfully", post)
}

func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	id, err := h.GetIDParam(c)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	err = h.PostService.DeletePost(id)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete post", err)
	}

	return h.SuccessResponse(c, "Post deleted successfully", nil)
}
