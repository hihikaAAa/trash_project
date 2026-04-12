package postgres

import (
	"abr_paperless_office/internal/models"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DocsRepository interface {
	FindDocID(gc *gin.Context, docID string) (*string, error)
	FindDocStatus(gc *gin.Context, linkID string) (*string, error)
	FindSecureTokenHashByCode(gc *gin.Context, linkID, hash string) (bool, error)
	AddSecureTokenByCode(gc *gin.Context, lp models.DocPermission) error
	FindIsPEPAcceptedByDocID(gc *gin.Context, docID string) (*bool, error)
	FindIsPEPAcceptedByLinkID(gc *gin.Context, linkID string) (*bool, error)
}

type docsRepository struct {
	db   DB
	pool *pgxpool.Pool
}

func NewDocsRepository(db *pgxpool.Pool) DocsRepository {
	return &docsRepository{
		db:   db,
		pool: db,
	}
}

func NewDocsRepositoryWithDB(db DB, pool *pgxpool.Pool) DocsRepository {
	return &docsRepository{
		db:   db,
		pool: pool,
	}
}

func (d *docsRepository) FindDocID(gc *gin.Context, docID string) (*string, error) {
	ctx := gc.Request.Context()

	var exists bool
	err := d.db.QueryRow(ctx, `
	SELECT EXISTS(SELECT 1 FROM paperless_office WHERE doc_id=$1)`, docID).Scan(&exists)

	if err != nil {
		return nil, fmt.Errorf("query docID: %w", err)
	}

	if !exists {
		return nil, pgx.ErrNoRows
	}

	var linkID *string
	err = d.db.QueryRow(ctx, `
	SELECT link_id
	FROM paperless_office
	WHERE doc_id=$1`, docID).Scan(&linkID)

	if err != nil {
		return nil, fmt.Errorf("query link_id by doc_id: %w", err)
	}

	return linkID, nil
}

func (d *docsRepository) FindDocStatus(gc *gin.Context, linkID string) (*string, error) {
	ctx := gc.Request.Context()

	var status *string
	err := d.db.QueryRow(ctx, `
	SELECT doc_status
	FROM paperless_office
	WHERE link_id=$1`, linkID).Scan(&status)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, pgx.ErrNoRows
	}

	if err != nil {
		return nil, fmt.Errorf("query secure token hash: %w", err)
	}

	return status, nil
}

func (d *docsRepository) FindSecureTokenHashByCode(gc *gin.Context, linkID, hash string) (bool, error) {
	ctx := gc.Request.Context()

	var exists bool
	err := d.db.QueryRow(ctx, `
    SELECT EXISTS(SELECT 1 FROM secure_tokens WHERE link_id=$1 AND secure_token_hash=$2)`, linkID, hash).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("query secure_token_hash: %w", err)
	}

	return exists, nil
}

func (d *docsRepository) AddSecureTokenByCode(gc *gin.Context, lp models.DocPermission) error {
	ctx := gc.Request.Context()

	_, err := d.db.Exec(ctx, `
	INSERT INTO secure_tokens (link_id, secure_token_hash)
	VALUES ($1, $2)`, lp.LinkID, lp.SecureTokenHash)

	if err != nil {
		return fmt.Errorf("insert secure_token for link: %w", err)
	}

	return nil
}

func (d *docsRepository) FindIsPEPAcceptedByDocID(gc *gin.Context, docID string) (*bool, error) {
	ctx := gc.Request.Context()

	var isSigned *bool
	err := d.db.QueryRow(ctx, `
    SELECT is_pep_accepted
	FROM paperless_office
	WHERE doc_id=$1`, docID).Scan(&isSigned)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, pgx.ErrNoRows
	}

	if err != nil {
		return nil, fmt.Errorf("query isPEPAccepted: %w", err)
	}

	return isSigned, nil
}

func (d *docsRepository) FindIsPEPAcceptedByLinkID(gc *gin.Context, linkID string) (*bool, error) {
	ctx := gc.Request.Context()

	var isSigned *bool
	err := d.db.QueryRow(ctx, `
    SELECT is_pep_accepted
	FROM paperless_office
	WHERE link_id=$1`, linkID).Scan(&isSigned)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, pgx.ErrNoRows
	}

	if err != nil {
		return nil, fmt.Errorf("query isPEPAccepted: %w", err)
	}

	return isSigned, nil
}
