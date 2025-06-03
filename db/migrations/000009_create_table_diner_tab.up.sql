CREATE TABLE IF NOT EXISTS `diner_tab`
(
    `id`
    INT
    NOT
    NULL
    AUTO_INCREMENT,
    `primary_user_id`
    INT
    NOT
    NULL,
    `table_id`
    INT
    NOT
    NULL,
    `client_name`
    VARCHAR
(
    255
) NOT NULL,
    `orders` JSON NOT NULL,
    `status` VARCHAR
(
    255
) NOT NULL DEFAULT 'open',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` INT NOT NULL,
    `update_info` JSON NULL,
    PRIMARY KEY
(
    `id`
)
    ) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4 COLLATE =utf8mb4_unicode_ci;