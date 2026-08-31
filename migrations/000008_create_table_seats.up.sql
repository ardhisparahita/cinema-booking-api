CREATE TABLE seats (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    studio_id BIGINT UNSIGNED NOT NULL,
    row_label VARCHAR(2) NOT NULL,
    col_number SMALLINT UNSIGNED NOT NULL,
    seat_type ENUM('regular', 'vip') NOT NULL DEFAULT 'regular',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_seats_studio FOREIGN KEY (studio_id) REFERENCES studios (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT uq_seats_position UNIQUE (
        studio_id,
        row_label,
        col_number
    ),
    INDEX idx_seats_studio_id (studio_id)
) Engine = InnoDB;