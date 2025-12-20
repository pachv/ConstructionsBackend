package handler

// import (
// 	"errors"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/pachv/constructions/constructions/internal/services"
// 	"github.com/pachv/constructions/constructions/transport/http/handler/responses"
// )

// func (h *Handler) AdminGetAllReviews(c *gin.Context) {
// 	items, err := h.reviewAdminService.GetAll(c.Request.Context())
// 	if err != nil {
// 		responses.BadRequestResponse(c, "cant get reviews: "+err.Error())
// 		return
// 	}
// 	responses.OkResponse(c, gin.H{"items": items})
// }

// func (h *Handler) AdminUpdateReview(c *gin.Context) {
// 	var req services.UpdateReviewRequest
// 	if err := c.BindJSON(&req); err != nil {
// 		responses.BadRequestResponse(c, "cant bind json: "+err.Error())
// 		return
// 	}

// 	if err := h.reviewAdminService.UpdateOne(c.Request.Context(), req); err != nil {
// 		if errors.Is(err, services.ErrReviewNotFound) {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
// 			return
// 		}
// 		responses.BadRequestResponse(c, "cant update review: "+err.Error())
// 		return
// 	}

// 	responses.OkResponse(c, gin.H{})
// }

// func (h *Handler) AdminDeleteReview(c *gin.Context) {
// 	id := c.Param("id")

// 	if err := h.reviewAdminService.DeleteOne(c.Request.Context(), id); err != nil {
// 		if errors.Is(err, service.ErrReviewNotFound) {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
// 			return
// 		}
// 		responses.BadRequestResponse(c, "cant delete review: "+err.Error())
// 		return
// 	}

// 	responses.OkResponse(c, gin.H{})
// }
