

CREATE TABLE i_user (
  created_at timestamp NOT NULL,
  email string(255) NOT NULL,
  name string(20) NOT NULL,
  updated_at timestamp NOT NULL,
  user_id string(36) NOT NULL,
) PRIMARY KEY(
    user_id
);

CREATE TABLE i_user_profile (
  bio string(MAX) NOT NULL,
  created_at timestamp NOT NULL,
  profile_id string(36) NOT NULL,
  updated_at timestamp NOT NULL,
  user_id string(36) NOT NULL,
  website string(255) NOT NULL,
) PRIMARY KEY(
    profile_id
);
ALTER TABLE i_user_profile
ADD CONSTRAINT fk_i_user_profile_user_id
FOREIGN KEY (user_id) REFERENCES i_user(user_id);