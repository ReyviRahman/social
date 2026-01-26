CREATE TABLE IF NOT EXISTS user_invitations (
  token bytea PRIMARY KEY,
  user_id bigint NOT NULL,
  expiry TIMESTAMP(0) WITH TIME ZONE NOT NULL,
  
  CONSTRAINT fk_user 
    FOREIGN KEY (user_id) 
    REFERENCES users (id) 
    ON DELETE CASCADE
);