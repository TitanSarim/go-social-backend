package store

import (
	"context"

	"github.com/lib/pq"
)

func (s *PostStore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetadata ,error) {

	// shouldFilterByOwnPost := len(fq.Tags) >= 1 || fq.Search != ""
	// p.user_id = $1

	query := `
				SELECT p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags, u.username, COUNT(c.id) AS comments_count 
				FROM posts p
				LEFT JOIN comments c ON c.post_id = p.id
				LEFT JOIN users u ON p.user_id = u.id
				JOIN followers f ON f.follower_id = p.user_id OR p.user_id = $1
				WHERE 
					f.user_id = $1  AND 
					(p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%') AND
					(p.tags @> $5 OR $5 = '{}')
				GROUP BY p.id, u.username
				ORDER BY p.created_at ` + fq.Sort +`
				LIMIT $2 OFFSET $3
			`
	rows, err := s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset, fq.Search, pq.Array(fq.Tags))
	if err !=nil {
        return nil, err
    }

	defer rows.Close()

	var feed []PostWithMetadata
	for rows.Next() {
		var p PostWithMetadata
        err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt, &p.Version, pq.Array(&p.Tags), &p.User.Username, &p.CommentCount)
        if err != nil {
            return nil, err
        }
        feed = append(feed, p)
	}
	return feed, nil
}