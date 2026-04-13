package v1

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hihikaAAa/trash_project/internal/domain/task"
	domainerrors "github.com/hihikaAAa/trash_project/internal/domainerrors"
	"github.com/hihikaAAa/trash_project/internal/service/services"
	httpres "github.com/hihikaAAa/trash_project/pkg/http_res"
)

type createOrderRequest struct {
	Address       string     `json:"address" binding:"required"`
	Description   *string    `json:"description"`
	PreferredTime *time.Time `json:"preferred_time"`
}

type assignOrderRequest struct {
	WorkerID string `json:"worker_id" binding:"required"`
}

func (h *Handler) initOrdersRoutes(router *gin.RouterGroup) {
	router.POST("", h.createOrder)
	router.GET("/:id", h.getOrder)
	router.GET("/my", h.listOwnOrders)

	router.GET("/available", h.listAvailableOrders)
	router.GET("/assigned", h.listAssignedOrders)
	router.POST("/:id/accept", h.acceptOrder)
	router.POST("/:id/complete", h.completeOrder)

	router.GET("/all", h.listAllOrders)
	router.POST("/:id/assign", h.assignOrder)
	router.POST("/:id/reassign", h.reassignOrder)
	router.POST("/:id/cancel", h.cancelOrder)
}

func (h *Handler) createOrder(ctx *gin.Context) {
	actor, err := actorFromContext(ctx)
	if err != nil {
		handleErr(ctx, err)
		return
	}

	var req createOrderRequest
	if err = ctx.ShouldBindJSON(&req); err != nil {
		handleErr(ctx, httpres.NewHTTPError(err, http.StatusBadRequest, httpres.CodeBadRequest, "bad_request"))
		return
	}

	input := services.CreateOrderInput{
		Address:       req.Address,
		Description:   req.Description,
		PreferredTime: req.PreferredTime,
	}

	order, err := h.services.Orders.Create(ctx.Request.Context(), actor, input)
	if err != nil {
		handleErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, order)
}

func (h *Handler) getOrder(ctx *gin.Context) {
	actor, err := actorFromContext(ctx)
	if err != nil {
		handleErr(ctx, err)
		return
	}

	orderID, err := parsePathID(ctx.Param("id"))
	if err != nil {
		handleErr(ctx, err)
		return
	}

	order, err := h.services.Orders.GetByID(ctx.Request.Context(), actor, orderID)
	if err != nil {
		handleErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func (h *Handler) listOwnOrders(ctx *gin.Context) {
	actor, err := actorFromContext(ctx)
	if err != nil {
		handleErr(ctx, err)
		return
	}

	orders, err := h.services.Orders.ListOwn(ctx.Request.Context(), actor)
	if err != nil {
		handleErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

func (h *Handler) listAvailableOrders(ctx *gin.Context) {
	actor, err := actorFromContext(ctx)
	if err != nil {
		handleErr(ctx, err)
		return
	}

	orders, err := h.services.Orders.ListAvailable(ctx.Request.Context(), actor)
	if err != nil {
		handleErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

func (h *Handler) listAssignedOrders(ctx *gin.Context) {
	actor, err := actorFromContext(ctx)
	if err != nil {
		handleErr(ctx, err)
		return
	}

	orders, err := h.services.Orders.ListAssigned(ctx.Request.Context(), actor)
	if err != nil {
		handleErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

func (h *Handler) acceptOrder(ctx *gin.Context) {
	actor, err := actorFromContext(ctx)
	if err != nil {
		handleErr(ctx, err)
		return
	}

	orderID, err := parsePathID(ctx.Param("id"))
	if err != nil {
		handleErr(ctx, err)
		return
	}

	order, err := h.services.Orders.Accept(ctx.Request.Context(), actor, orderID)
	if err != nil {
		handleErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func (h *Handler) completeOrder(ctx *gin.Context) {
	actor, err := actorFromContext(ctx)
	if err != nil {
		handleErr(ctx, err)
		return
	}

	orderID, err := parsePathID(ctx.Param("id"))
	if err != nil {
		handleErr(ctx, err)
		return
	}

	order, err := h.services.Orders.Complete(ctx.Request.Context(), actor, orderID)
	if err != nil {
		handleErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func (h *Handler) listAllOrders(ctx *gin.Context) {
	actor, err := actorFromContext(ctx)
	if err != nil {
		handleErr(ctx, err)
		return
	}

	orders, err := h.services.Orders.ListAll(ctx.Request.Context(), actor)
	if err != nil {
		handleErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

func (h *Handler) assignOrder(ctx *gin.Context) {
	h.assignLike(ctx)
}

func (h *Handler) reassignOrder(ctx *gin.Context) {
	h.assignLike(ctx)
}

func (h *Handler) assignLike(ctx *gin.Context) {
	actor, err := actorFromContext(ctx)
	if err != nil {
		handleErr(ctx, err)
		return
	}

	orderID, err := parsePathID(ctx.Param("id"))
	if err != nil {
		handleErr(ctx, err)
		return
	}

	var req assignOrderRequest
	if err = ctx.ShouldBindJSON(&req); err != nil {
		handleErr(ctx, httpres.NewHTTPError(err, http.StatusBadRequest, httpres.CodeBadRequest, "bad_request"))
		return
	}

	workerID, err := uuid.Parse(req.WorkerID)
	if err != nil {
		handleErr(ctx, httpres.NewHTTPError(err, http.StatusBadRequest, httpres.CodeBadRequest, "bad_request"))
		return
	}

	order, err := h.services.Orders.Assign(ctx.Request.Context(), actor, orderID, workerID)
	if err != nil {
		handleErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func (h *Handler) cancelOrder(ctx *gin.Context) {
	actor, err := actorFromContext(ctx)
	if err != nil {
		handleErr(ctx, err)
		return
	}

	orderID, err := parsePathID(ctx.Param("id"))
	if err != nil {
		handleErr(ctx, err)
		return
	}

	order, err := h.services.Orders.Cancel(ctx.Request.Context(), actor, orderID)
	if err != nil {
		handleErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func actorFromContext(ctx *gin.Context) (services.Actor, error) {
	uid, err := extractUserID(ctx)
	if err != nil {
		return services.Actor{}, err
	}
	role, err := extractRole(ctx)
	if err != nil {
		return services.Actor{}, err
	}
	return services.Actor{ID: uid, Role: role}, nil
}

func extractUserID(ctx *gin.Context) (uuid.UUID, error) {
	keys := []string{"user_id", "userID", "userId", "account_id"}
	for _, key := range keys {
		if v, ok := ctx.Get(key); ok {
			switch val := v.(type) {
			case uuid.UUID:
				if val == uuid.Nil {
					return uuid.Nil, domainerrors.ErrForbidden
				}
				return val, nil
			case string:
				id, err := uuid.Parse(strings.TrimSpace(val))
				if err != nil {
					return uuid.Nil, domainerrors.ErrForbidden
				}
				return id, nil
			}
		}
	}
	return uuid.Nil, domainerrors.ErrForbidden
}

func extractRole(ctx *gin.Context) (task.Role, error) {
	keys := []string{"role", "user_role", "user_group"}
	for _, key := range keys {
		if v, ok := ctx.Get(key); ok {
			s, ok := v.(string)
			if !ok {
				continue
			}
			r, err := services.ParseRole(strings.ToLower(strings.TrimSpace(s)))
			if err != nil {
				return "", domainerrors.ErrForbidden
			}
			return r, nil
		}
	}
	return "", domainerrors.ErrForbidden
}

func parsePathID(value string) (uuid.UUID, error) {
	id, err := uuid.Parse(strings.TrimSpace(value))
	if err != nil {
		return uuid.Nil, httpres.NewHTTPError(err, http.StatusBadRequest, httpres.CodeBadRequest, "bad_request")
	}
	return id, nil
}

func handleErr(ctx *gin.Context, err error) {
	httpres.HandleDomainErrors(ctx, err, nil)
}
