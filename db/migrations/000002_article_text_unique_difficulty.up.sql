ALTER TABLE article_texts
ADD CONSTRAINT article_texts_article_id_difficulty_unique UNIQUE NULLS NOT DISTINCT (article_id, difficulty, deleted_at);
