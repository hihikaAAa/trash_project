package v1

import (
	"abr_paperless_office/internal/domainerrors"
	"abr_paperless_office/internal/metrics"
	"abr_paperless_office/internal/models"
	httpres "abr_paperless_office/pkg/http_res"
	"abr_paperless_office/pkg/utils"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initLinksRoutes(router *gin.RouterGroup) {
	router.POST("/clientLink", h.GenerateLink)
	router.GET("/linkInfo", h.GetLinkInfo)
}

// GenerateLink GoDoc
// @Summary Создание ссылки
// @Description По переданному crm_uuid создается уникальная короткая ссылка
// @Tags Ссылки
// @Produce json
// @Param   data  body   models.LinkGenerateInput true "link data json format"
// @Success 200 {object} models.LinkGenerateResponse
// @Success 201 {object} models.LinkGenerateResponse
// @Failure 400 {object} httpres.HTTPError
// @Failure 500 {object} httpres.HTTPError
// @Router /links/clientLink [post]
func (h *Handler) GenerateLink(ctx *gin.Context) {
	var in models.LinkGenerateInput

	if err := ctx.ShouldBindJSON(&in); err != nil {
		if errors.Is(err, io.EOF) {
			httpres.HandleDomainError(ctx, domainerrors.ErrEmptyRequestBody, nil, nil)
			return
		}
		httpres.NewError(ctx, http.StatusBadRequest, httpres.CodeBadRequest, err, nil)
		return
	}

	gl, err := h.services.Links.GenerateLink(ctx, in, h.cnf.Links.Domain, h.cnf.Links.CodeLength)
	if err != nil {
		httpres.HandleDomainError(ctx, err, metrics.LinkGenerateTotal, nil)
		return
	}

	if gl.WasCreated {
		metrics.LinkGenerateTotal.WithLabelValues("created").Inc()
		ctx.JSON(http.StatusCreated, models.FilteredGeneratedLinkResponse(gl))
		return
	}
	metrics.LinkGenerateTotal.WithLabelValues("existing").Inc()
	ctx.JSON(http.StatusOK, models.FilteredGeneratedLinkResponse(gl))
}

// GetLinkInfo GoDoc
// @Summary Получение информации о ссылке
// @Description Сервис вызывается клиентской частью при переходе пользователя по ссылке. В ответе 200 поле documents.docId возвращается только если у secure_token есть доступ к документу.
// @Tags Ссылки
// @Produce json
// @Param   linkId   query   string  true  "Link ID"
// @Success 200 {object} models.LinkInfoResponse
// @Failure 400 {object} httpres.HTTPError
// @Failure 403 {object} httpres.HTTPError
// @Failure 404 {object} httpres.HTTPError
// @Failure 500 {object} httpres.HTTPError
// @Router /links/linkInfo [get]
func (h *Handler) GetLinkInfo(ctx *gin.Context) {
	linkID := ctx.Query("linkId")
	if linkID == "" {
		httpres.NewError(ctx, http.StatusBadRequest, httpres.CodeNotExistantLink, domainerrors.ErrLinkNotFound, nil)
		return
	}

	tokenDetails := utils.ReadOptionalSecureToken(ctx)

	out, err := h.services.Links.GetLinkInfo(ctx, models.LinkInfoInput{LinkID: linkID}, tokenDetails)
	if err != nil {
		httpres.HandleDomainError(ctx, err, metrics.LinkInfoTotal, nil)
		return
	}
	metrics.LinkInfoTotal.WithLabelValues("ok").Inc()
	ctx.JSON(http.StatusOK, models.FilteredLinkInfoResponse(out))
}
