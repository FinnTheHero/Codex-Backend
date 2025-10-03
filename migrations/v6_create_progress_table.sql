CREATE TABLE progress (
    id uuid NOT NULL DEFAULT gen_random_uui(),
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    novel_id INTEGER NOT NULL REFERENCES novels(id) ON DELETE CASCADE,
    chapter_id INTEGER NOT NULL REFERENCES chapters(id) ON DELETE CASCADE,

    scroll_position INTEGER DEFAULT 0,
    progress_percentage DECIMAL(5,2) DEFAULT 0.00,

    created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
    completed boolean DEFAULT false,

    UNIQUE(user_id, chapter_id),
    CHECK (progress_percentage >= 0 AND progress_percentage <= 100)
)

CREATE INDEX idx_progress_user_id_novel_id ON progress (user_id, novel_id);
CREATE INDEX idx_progress_updated_at ON progress (updated_at);
