-- Enable the pg_trgm extension if it is not already enabled.
-- This extension provides support for trigram-based indexing, improving text search performance.
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Create a GIN index on the 'content' column of the 'comments' table using trigram operations.
-- This helps speed up full-text searches on comment content.
CREATE INDEX IF NOT EXISTS idx_comments_content ON comments USING gin (content gin_trgm_ops);

-- Create a GIN index on the 'title' column of the 'posts' table using trigram operations.
-- This allows for efficient searches and similarity matching on post titles.
CREATE INDEX IF NOT EXISTS idx_posts_title ON posts USING gin (title gin_trgm_ops);

-- Create a GIN index on the 'tags' column of the 'posts' table.
-- This is useful for searching and filtering posts by tags efficiently.
CREATE INDEX IF NOT EXISTS idx_posts_tags ON posts USING gin (tags);

-- Create a B-tree index on the 'username' column of the 'users' table.
-- This speeds up queries that search for users by username.
CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);

-- Create a B-tree index on the 'user_id' column of the 'posts' table.
-- This optimizes queries that filter or join posts based on the user who created them.
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts (user_id);

-- Create a B-tree index on the 'post_id' column of the 'comments' table.
-- This speeds up queries that retrieve comments for a specific post.
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments (post_id);
