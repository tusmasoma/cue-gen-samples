

CREATE TABLE m_regions (
  country_code string(10) NOT NULL,
  created_at timestamp NOT NULL,
  description string(255) NOT NULL,
  region_id string(36) NOT NULL,
  region_name string(100) NOT NULL,
  updated_at timestamp NOT NULL,
) PRIMARY KEY(
    region_id
);
CREATE UNIQUE NULL_FILTERED INDEX m_regions_by_country_code ON m_regions (country_code);