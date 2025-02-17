

CREATE TABLE users (
  email string(255) NOT NULL,
  name string(20) NOT NULL,
  user_id string(36) NOT NULL,
) PRIMARY KEY(
    user_id
);