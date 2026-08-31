CREATE TABLE movies (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    synopsis TEXT NULL,
    duration_min SMALLINT UNSIGNED NOT NULL,
    rating ENUM('SU', '13+', '17+', '21+') NOT NULL,
    poster_url VARCHAR(255) NULL,
    release_date DATE NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_movies_release_date (release_date),
    INDEX idx_movies_is_active (is_active)
) ENGINE = InnoDB;