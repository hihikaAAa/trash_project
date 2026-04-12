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
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initDocsRoutes(router *gin.RouterGroup) {
	router.POST("/otp", h.GetDocPermission)
	router.GET("/docDownload", h.DownloadDoc)
	router.POST("/docSign", h.SignDoc)
	router.POST("/acceptConditions", h.AcceptConditions)
}

// GetDocPermission GoDoc
// @Summary Получение доступа к документу
// @Description Сервис принимает одноразовый код, запрашивает доступ в CRM и при успехе возвращает только статус 200.
// @Tags Документы
// @Accept json
// @Produce json
// @Param   data  body   models.DocPermissionInput true "otp data json format"
// @Success 200 "OK"
// @Failure 400 {object} httpres.HTTPError
// @Failure 403 {object} httpres.HTTPError "code wrong_otp (with attemptsLeft field) or code denied (without attemptsLeft field)"
// @Failure 404 {object} httpres.HTTPError
// @Failure 500 {object} httpres.HTTPError
// @Router /docs/otp [post]
func (h *Handler) GetDocPermission(ctx *gin.Context) {
	var in models.DocPermissionInput

	if err := ctx.ShouldBindJSON(&in); err != nil {
		if errors.Is(err, io.EOF) {
			httpres.HandleDomainError(ctx, domainerrors.ErrEmptyRequestBody, nil, nil)
			return
		}
		httpres.NewError(ctx, http.StatusBadRequest, httpres.CodeBadRequest, err, nil)
		return
	}

	dp, err := h.services.Docs.GetAccessToDoc(ctx, in)
	if err != nil {
		if dp != nil {
			httpres.HandleDomainError(ctx, err, metrics.DocPermissionTotal, &dp.AttemptsLeft)
			return
		}
		httpres.HandleDomainError(ctx, err, metrics.DocPermissionTotal, nil)
		return
	}

	maxAge := 0
	if !dp.ExpireAt.IsZero() {
		maxAge = max(int(time.Until(dp.ExpireAt).Seconds()), 0)
	}

	metrics.DocPermissionTotal.WithLabelValues("granted").Inc()
	ctx.SetCookie("secure_token", dp.SecureTokenPlain, maxAge, "/", "", false, true)
	ctx.Status(http.StatusOK)
}

// DownloadDoc GoDoc
// @Summary Получение документа
// @Description Возвращает документ в бинарном формате.
// @Tags Документы
// @Produce application/pdf
// @Param   docId   query   string  true  "Doc ID"
// @Success 200 {file} file
// @Failure 400 {object} httpres.HTTPError
// @Failure 403 {object} httpres.HTTPError
// @Failure 404 {object} httpres.HTTPError
// @Failure 500 {object} httpres.HTTPError
// @Router /docs/docDownload [get]
func (h *Handler) DownloadDoc(ctx *gin.Context) {
	docID := ctx.Query("docId")
	if docID == "" {
		httpres.NewError(ctx, http.StatusBadRequest, httpres.CodeNotExistantDoc, domainerrors.ErrNoDoc, nil)
		return
	}

	tokenDetails, err := utils.ReadRequiredSecureToken(ctx)
	if err != nil {
		httpres.HandleDomainError(ctx, err, nil, nil)
		return
	}

	doc, err := h.services.Docs.DownloadDoc(ctx, models.DocDownloadInput{DocID: docID}, tokenDetails)
	if err != nil {
		httpres.HandleDomainError(ctx, err, metrics.DocDownloadTotal, nil)
		return
	}

	contentType := http.DetectContentType(doc)
	ctx.Data(http.StatusOK, contentType, doc)
}

// SignDoc GoDoc
// @Summary Подпись документа
// @Description Сервис принимает на вход docId документа и передает в CRM информацию о подписи.
// @Tags Документы
// @Accept json
// @Produce json
// @Param   data  body   models.DocSignInput true "doc sign data json format"
// @Success 200 "OK"
// @Failure 400 {object} httpres.HTTPError
// @Failure 403 {object} httpres.HTTPError
// @Failure 404 {object} httpres.HTTPError
// @Failure 500 {object} httpres.HTTPError
// @Router /docs/docSign [post]
func (h *Handler) SignDoc(ctx *gin.Context) {
	var in models.DocSignInput

	if err := ctx.ShouldBindJSON(&in); err != nil {
		if errors.Is(err, io.EOF) {
			httpres.HandleDomainError(ctx, domainerrors.ErrEmptyRequestBody, nil, nil)
			return
		}
		httpres.NewError(ctx, http.StatusBadRequest, httpres.CodeBadRequest, err, nil)
		return
	}

	tokenDetails, err := utils.ReadRequiredSecureToken(ctx)
	if err != nil {
		httpres.HandleDomainError(ctx, err, metrics.DocSignTotal, nil)
		return
	}

	err = h.services.Docs.SignDoc(ctx, in, tokenDetails)
	if err != nil {
		httpres.HandleDomainError(ctx, err, metrics.DocSignTotal, nil)
		return
	}

	metrics.DocSignTotal.WithLabelValues("success").Inc()
	ctx.Status(http.StatusOK)
}

// AcceptConditions GoDoc
// @Summary Передача согласия на использование ПЭП
// @Description Сервис принимает на вход UUID ссылки и тип согласия, после чего передает информацию в CRM
// @Tags Документы
// @Accept json
// @Produce json
// @Param   data  body   models.AcceptConditionsInput true "accept conditions data json format"
// @Success 200 "OK"
// @Failure 400 {object} httpres.HTTPError
// @Failure 403 {object} httpres.HTTPError
// @Failure 404 {object} httpres.HTTPError
// @Failure 409 {object} httpres.HTTPError
// @Failure 500 {object} httpres.HTTPError
// @Router /docs/acceptConditions [post]
func (h *Handler) AcceptConditions(ctx *gin.Context) {
	var in models.AcceptConditionsInput

	if err := ctx.ShouldBindJSON(&in); err != nil {
		if errors.Is(err, io.EOF) {
			httpres.HandleDomainError(ctx, domainerrors.ErrEmptyRequestBody, nil, nil)
			return
		}
		httpres.NewError(ctx, http.StatusBadRequest, httpres.CodeBadRequest, err, nil)
		return
	}

	tokenDetails, err := utils.ReadRequiredSecureToken(ctx)
	if err != nil {
		httpres.HandleDomainError(ctx, err, metrics.AcceptConditionsTotal, nil)
		return
	}

	err = h.services.Docs.AcceptConditions(ctx, in, tokenDetails)
	if err != nil {
		httpres.HandleDomainError(ctx, err, metrics.AcceptConditionsTotal, nil)
		return
	}
	metrics.AcceptConditionsTotal.WithLabelValues("success").Inc()
	ctx.Status(http.StatusOK)
}
