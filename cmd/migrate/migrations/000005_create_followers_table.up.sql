CREATE TABLE IF NOT EXISTS followers (
  user_id bigint NOT NULL,
  following_id bigint NOT NULL,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),

  CHECK (user_id <> following_id),
  PRIMARY KEY (user_id, following_id),
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
  FOREIGN KEY (following_id) REFERENCES users (id) ON DELETE CASCADE
);