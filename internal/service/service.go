// Package service
package service

import (
	"abr_paperless_office/internal/models"
	adapters "abr_paperless_office/internal/repositories"
	"abr_paperless_office/internal/service/services"
	"abr_paperless_office/pkg/utils"

	"github.com/gin-gonic/gin"
)

type Links interface {
	GenerateLink(gc *gin.Context, lg models.LinkGenerateInput, domain string, codeLength int) (*models.GeneratedLink, error)
	GetLinkInfo(gc *gin.Context, li models.LinkInfoInput, secureToken *utils.TokenDetails) (*models.DocInfo, error)
}

type Docs interface {
	GetAccessToDoc(gc *gin.Context, in models.DocPermissionInput) (*models.DocPermission, error)
	DownloadDoc(gc *gin.Context, in models.DocDownloadInput, secureToken *utils.TokenDetails) ([]byte, error)
	SignDoc(gc *gin.Context, in models.DocSignInput, secureToken *utils.TokenDetails) error
	AcceptConditions(gc *gin.Context, in models.AcceptConditionsInput, secureToken *utils.TokenDetails) error
}

type Service struct {
	Links Links
	Docs  Docs
}

func NewService(repository *adapters.Repository, crmClient services.CRMClient) *Service {
	return &Service{
		Links: services.NewLinksService(repository.Links, repository.Docs, crmClient),
		Docs:  services.NewDocsService(repository.Links, repository.Docs, crmClient),
	}
}
