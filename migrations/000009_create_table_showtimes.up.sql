CREATE TABLE showtimes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    movie_id BIGINT UNSIGNED NOT NULL,
    studio_id BIGINT UNSIGNED NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_showtimes_movie FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_showtimes_studio FOREIGN KEY (studio_id) REFERENCES studios (id) ON DELETE RESTRICT ON UPDATE CASCADE,
    INDEX idx_showtimes_movie_id (movie_id),
    INDEX idx_showtimes_studio_id (studio_id),
    INDEX idx_showtimes_start_time (start_time)
) Engine = InnoDB;