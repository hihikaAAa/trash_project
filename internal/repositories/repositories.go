// Package repositories
package repositories

import (
	"abr_paperless_office/internal/models"
	"abr_paperless_office/internal/repositories/postgres"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Links interface {
	AddLink(gc *gin.Context, gl models.GeneratedLink) error
	FindLinkByCrmUUID(gc *gin.Context, crmUUID uuid.UUID) (*models.GeneratedLink, error)
	CheckIfExistLink(gc *gin.Context, linkID string) (bool, error)
	FindCrmUUIDByCode(gc *gin.Context, linkID string) (uuid.UUID, error)
	FindDocInfo(gc *gin.Context, link string) (*models.DocInfo, error)
	UpdateDocInfo(gc *gin.Context, linkID string, di models.DocInfo, now time.Time) error
	CountTotalLinks(gc *gin.Context) (int, error)
}

type Docs interface {
	FindDocID(gc *gin.Context, docID string) (*string, error)
	FindDocStatus(gc *gin.Context, linkID string) (*string, error)
	FindSecureTokenHashByCode(gc *gin.Context, linkID, hash string) (bool, error)
	AddSecureTokenByCode(gc *gin.Context, lp models.DocPermission) error
	FindIsPEPAcceptedByDocID(gc *gin.Context, docID string) (*bool, error)
	FindIsPEPAcceptedByLinkID(gc *gin.Context, linkID string) (*bool, error)
}

type Repository struct {
	Links Links
	Docs  Docs
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Links: postgres.NewLinksRepository(db),
		Docs:  postgres.NewDocsRepository(db),
	}
}
