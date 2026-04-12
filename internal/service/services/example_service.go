package services

import (
	"abr_paperless_office/internal/domainerrors"
	"abr_paperless_office/internal/models"
	"abr_paperless_office/internal/repositories"
	"abr_paperless_office/pkg/logger"
	"abr_paperless_office/pkg/utils"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type DocsService struct {
	LinksRepository repositories.Links
	DocsRepository  repositories.Docs
	CRM             CRMClient
}

func NewDocsService(linksRepo repositories.Links, docsRepo repositories.Docs, crm CRMClient) *DocsService {
	return &DocsService{
		LinksRepository: linksRepo,
		DocsRepository:  docsRepo,
		CRM:             crm,
	}
}

func (d *DocsService) GetAccessToDoc(gc *gin.Context, in models.DocPermissionInput) (*models.DocPermission, error) {
	exists, err := d.LinksRepository.CheckIfExistLink(gc, in.LinkID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domainerrors.ErrLinkNotFound
	}

	crmUUID, err := d.LinksRepository.FindCrmUUIDByCode(gc, in.LinkID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerrors.ErrLinkNotFound
		}
		return nil, err
	}

	_, err = d.LinksRepository.FindDocInfo(gc, in.LinkID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}

	logger.SetEventLinkID(gc, in.LinkID)
	crmCtx := logger.WithEventLinkID(gc.Request.Context(), in.LinkID)

	token, err := utils.CreateToken()
	if err != nil {
		return nil, err
	}

	exists, err = d.DocsRepository.FindSecureTokenHashByCode(gc, in.LinkID, *token.TokenHash)
	if err != nil {
		return nil, err
	}
	// Если такой токен существует - перегенерируем, чтобы был 100% уникальным
	if exists {
		for range 30 {
			token, err = utils.CreateToken()
			if err != nil {
				return nil, err
			}
			exists, err = d.DocsRepository.FindSecureTokenHashByCode(gc, in.LinkID, *token.TokenHash)
			if err != nil {
				return nil, err
			}
			if !exists {
				break
			}
		}
	}

	crmResp, err := d.CRM.GetDocumentPermission(crmCtx, models.GetDocumentPermissionRequest{
		UUID:        crmUUID.String(),
		Code:        in.OTP,
		UserAgent:   gc.GetHeader("User-Agent"),
		DeviceIP:    gc.ClientIP(),
		SecureToken: *token.Token,
	})

	if err != nil {
		var crmErr *domainerrors.CRMError
		if errors.As(err, &crmErr) {
			if mapped := mapCRMCodeToDomain(crmErr.Code); mapped != nil {
				return nil, mapped
			}
			return nil, crmErr
		}
		return nil, err
	}

	if crmResp.IsSuccess {
		docPerm := &models.DocPermission{
			GrantedAt:        time.Now(),
			LinkID:           in.LinkID,
			SecureTokenHash:  *token.TokenHash,
			SecureTokenPlain: *token.Token,
		}
		err = d.DocsRepository.AddSecureTokenByCode(gc, *docPerm)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, domainerrors.ErrLinkNotFound
			}
			return nil, err
		}
		di, findErr := d.LinksRepository.FindDocInfo(gc, docPerm.LinkID)
		if findErr != nil {
			if errors.Is(findErr, pgx.ErrNoRows) {
				return nil, domainerrors.ErrLinkNotFound
			}
			return nil, findErr
		}
		needsCRMRefresh := di == nil ||
			di.DocID == nil || strings.TrimSpace(*di.DocID) == "" ||
			di.DocStatus == nil || strings.TrimSpace(*di.DocStatus) == "" ||
			di.AttemptsLeft == nil ||
			di.IsPEPAccepted == nil

		if needsCRMRefresh {
			crmSecureToken := ""
			if token.Token != nil {
				crmSecureToken = *token.Token
			}

			crmInfo, infoErr := d.CRM.GetDocumentInfo(crmCtx,
				models.GetDocumentInfoRequest{
					UUID:        crmUUID.String(),
					UserAgent:   gc.GetHeader("User-Agent"),
					DeviceIP:    gc.ClientIP(),
					SecureToken: crmSecureToken,
				})
			if infoErr != nil {
				var crmErr *domainerrors.CRMError
				if errors.As(infoErr, &crmErr) {
					if mapped := mapCRMCodeToDomain(crmErr.Code); mapped !=
						nil {
						return nil, mapped
					}
					return nil, crmErr
				}
				return nil, infoErr
			}

			crmDocInfo := models.DocInfo{
				DocID:         crmInfo.Documents.DocID,
				DocName:       &crmInfo.Documents.DocName,
				DocStatus:     &crmInfo.Documents.DocStatus,
				AttemptsLeft:  &crmInfo.Documents.AttemptsLeft,
				ExpireAt:      crmInfo.Documents.ExpireAt,
				IsPEPAccepted: &crmInfo.Documents.IsPEPAccepted,
			}
			if err = d.LinksRepository.UpdateDocInfo(gc, docPerm.LinkID, crmDocInfo,
				time.Now()); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return nil, domainerrors.ErrLinkNotFound
				}
				return nil, err
			}

			di = &crmDocInfo
		}
		if di != nil && di.ExpireAt != nil {
			docPerm.ExpireAt = *di.ExpireAt
		}
		return docPerm, nil
	}

	if crmResp.AttemptsLeft != nil {
		docPerm := &models.DocPermission{
			AttemptsLeft: *crmResp.AttemptsLeft,
		}
		err = d.LinksRepository.UpdateDocInfo(gc, in.LinkID, models.DocInfo{AttemptsLeft: crmResp.AttemptsLeft}, time.Now())
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, domainerrors.ErrLinkNotFound
			}
			return nil, err
		}
		if *crmResp.AttemptsLeft == 0 {
			return nil, domainerrors.ErrNoAttemptsLeft
		}
		return docPerm, domainerrors.ErrWrongOTP
	}
	if mapped := mapCRMCodeToDomain(crmResp.ErrorCode); mapped != nil {
		return nil, mapped
	}
	if crmResp.ErrorMessage != "" {
		return nil, &domainerrors.CRMError{Code: crmResp.ErrorCode, Message: crmResp.ErrorMessage}
	}

	return nil, domainerrors.ErrInternalError
}

func (d *DocsService) DownloadDoc(gc *gin.Context, in models.DocDownloadInput, secureToken *utils.TokenDetails) ([]byte, error) {
	linkID, err := d.DocsRepository.FindDocID(gc, in.DocID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerrors.ErrNoDoc
		}
		return nil, err
	}

	logger.SetEventLinkID(gc, *linkID)
	crmCtx := logger.WithEventLinkID(gc.Request.Context(), *linkID)

	if secureToken == nil || secureToken.Token == nil {
		return nil, domainerrors.ErrWrongSecureToken
	}

	exists, err := d.DocsRepository.FindSecureTokenHashByCode(gc, *linkID, *secureToken.TokenHash)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domainerrors.ErrWrongSecureToken
	}

	status, err := d.DocsRepository.FindDocStatus(gc, *linkID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerrors.ErrWrongDocStatus
		}
		return nil, err
	}
	if status == nil || (*status != string(models.StatusSentForSignature) && *status != string(models.StatusSigned)) {
		return nil, domainerrors.ErrWrongDocStatus
	}

	crmResp, err := d.CRM.GetDocumentBinary(crmCtx, models.GetDocumentBinaryRequest{
		DocID:       in.DocID,
		SecureToken: *secureToken.Token,
	})
	if err != nil {
		var crmErr *domainerrors.CRMError
		if errors.As(err, &crmErr) {
			if mapped := mapCRMCodeToDomain(crmErr.Code); mapped != nil {
				return nil, mapped
			}
			return nil, crmErr
		}
		return nil, err
	}

	return crmResp, nil
}

func (d *DocsService) SignDoc(gc *gin.Context, in models.DocSignInput, secureToken *utils.TokenDetails) error {
	linkID, err := d.DocsRepository.FindDocID(gc, in.DocID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainerrors.ErrNoDoc
		}
		return err
	}

	logger.SetEventLinkID(gc, *linkID)
	crmCtx := logger.WithEventLinkID(gc.Request.Context(), *linkID)

	if secureToken == nil || secureToken.Token == nil {
		return domainerrors.ErrWrongSecureToken
	}

	exists, err := d.DocsRepository.FindSecureTokenHashByCode(gc, *linkID, *secureToken.TokenHash)
	if err != nil {
		return err
	}
	if !exists {
		return domainerrors.ErrWrongSecureToken
	}

	isSigned, err := d.DocsRepository.FindIsPEPAcceptedByDocID(gc, in.DocID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainerrors.ErrNoPEP
		}
		return err
	}

	if isSigned == nil || !*isSigned {
		return domainerrors.ErrNoPEP
	}

	status, err := d.DocsRepository.FindDocStatus(gc, *linkID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainerrors.ErrWrongDocStatus
		}
		return err
	}
	if status == nil || *status != string(models.StatusSentForSignature) {
		return domainerrors.ErrWrongDocStatus
	}

	err = d.CRM.SignDoc(crmCtx, models.SignDocRequest{
		DocID:       in.DocID,
		UserAgent:   gc.GetHeader("User-Agent"),
		DeviceIP:    gc.ClientIP(),
		SecureToken: *secureToken.Token,
	})
	if err != nil {
		var crmErr *domainerrors.CRMError
		if errors.As(err, &crmErr) {
			if mapped := mapCRMCodeToDomain(crmErr.Code); mapped != nil {
				return mapped
			}
			return crmErr
		}
		return err
	}

	newStatus := string(models.StatusSigned)
	err = d.LinksRepository.UpdateDocInfo(gc, *linkID, models.DocInfo{DocStatus: &newStatus}, time.Now())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainerrors.ErrNoDoc
		}
		return err
	}

	return nil
}

func (d *DocsService) AcceptConditions(gc *gin.Context, in models.AcceptConditionsInput, secureToken *utils.TokenDetails) error {
	logger.SetEventLinkID(gc, in.LinkID)
	crmCtx := logger.WithEventLinkID(gc.Request.Context(), in.LinkID)

	exists, err := d.LinksRepository.CheckIfExistLink(gc, in.LinkID)
	if err != nil {
		return err
	}
	if !exists {
		return domainerrors.ErrLinkNotFound
	}

	if secureToken == nil || secureToken.Token == nil {
		return domainerrors.ErrWrongSecureToken
	}

	exists, err = d.DocsRepository.FindSecureTokenHashByCode(gc, in.LinkID, *secureToken.TokenHash)
	if err != nil {
		return err
	}
	if !exists {
		return domainerrors.ErrWrongSecureToken
	}

	isSigned, err := d.DocsRepository.FindIsPEPAcceptedByLinkID(gc, in.LinkID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainerrors.ErrLinkNotFound
		}
		return err
	}
	if isSigned != nil && *isSigned {
		return domainerrors.ErrPepExists
	}

	status, err := d.DocsRepository.FindDocStatus(gc, in.LinkID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainerrors.ErrWrongDocStatus
		}
		return err
	}
	if status == nil || *status != string(models.StatusSentForSignature) {
		return domainerrors.ErrWrongDocStatus
	}

	crmUUID, err := d.LinksRepository.FindCrmUUIDByCode(gc, in.LinkID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainerrors.ErrLinkNotFound
		}
		return err
	}

	err = d.CRM.AcceptConditions(crmCtx, models.GrantPEPAccessRequest{
		UUID:        crmUUID.String(),
		Type:        in.Type,
		UserAgent:   gc.GetHeader("User-Agent"),
		DeviceIP:    gc.ClientIP(),
		SecureToken: *secureToken.Token,
	})
	if err != nil {
		var crmErr *domainerrors.CRMError
		if errors.As(err, &crmErr) {
			if mapped := mapCRMCodeToDomain(crmErr.Code); mapped != nil {
				return mapped
			}
			return crmErr
		}
		return err
	}

	accepted := true
	err = d.LinksRepository.UpdateDocInfo(
		gc, in.LinkID, models.DocInfo{IsPEPAccepted: &accepted}, time.Now(),
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainerrors.ErrLinkNotFound
		}
		return err
	}

	return nil
}
