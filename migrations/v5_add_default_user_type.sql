ALTER TABLE users
ALTER COLUMN type SET DEFAULT "User";

ALTER TABLE users
ADD CONSTRAINT chk_user_type
CHECK (type IN ('User', 'Admin'));
