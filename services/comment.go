package services

import (
	"context"
	"github.com/gocql/gocql"
	"time"

	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"go-fundraising/db"
	"go-fundraising/models"
)

type CommentService struct{}

func (s *CommentService) InsertComment(ctx context.Context, comment models.Comment) error {
	stmt, names := qb.Insert(models.CommentTable.Name).
		Columns(models.CommentTable.Columns...).
		ToCql()

	return gocqlx.Query(db.ScyllaSession.Query(stmt).Consistency(gocql.One), names).
		BindStruct(comment).
		ExecRelease()
}

func (s *CommentService) GetCommentsByPostID(ctx context.Context, postID string, perPage int, lastCreatedAt time.Time) ([]models.Comment, error) {

	qbSelect := qb.Select(models.CommentTable.Name).
		Where(qb.Eq("post_id")).
		OrderBy("created_at", qb.DESC).
		Limit(uint(perPage))

	if !lastCreatedAt.IsZero() {
		qbSelect = qbSelect.Where(qb.Lt("created_at"))
	}

	stmt, names := qbSelect.ToCql()

	var comments []models.Comment
	q := gocqlx.Query(db.ScyllaSession.Query(stmt).Consistency(gocql.One), names).
		BindMap(qb.M{
			"post_id": postID,
		})

	if !lastCreatedAt.IsZero() {
		q.Bind(lastCreatedAt)
	}

	if err := q.SelectRelease(&comments); err != nil {
		return nil, err
	}

	return comments, nil
}
