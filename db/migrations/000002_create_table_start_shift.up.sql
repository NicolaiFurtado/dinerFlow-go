CREATE TABLE IF NOT EXISTS start_shift (
                                           id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
                                           user_id BIGINT UNSIGNED NOT NULL,
                                           start_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                           end_time TIMESTAMP,
                                           FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );