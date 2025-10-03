ALTER TABLE chapters
ADD CONSTRAINT fk_chapters_novel_id
FOREIGN KEY (novel_id)
REFERENCES novels(id)
ON DELETE CASCADE;
