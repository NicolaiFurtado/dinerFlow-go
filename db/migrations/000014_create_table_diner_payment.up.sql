CREATE TABLE `diner-flow`.`diner_payments`
(
    `id`           INT            NOT NULL AUTO_INCREMENT,
    `user_id`      INT            NOT NULL,
    `client_name`  VARCHAR(255)   NOT NULL,
    `total_price`  DECIMAL(10, 2) NOT NULL,
    `receipt_data` JSON           NOT NULL,
    `status`       VARCHAR(255)   NOT NULL DEFAULT 'Not Paid',
    `type_payment` VARCHAR(255) NULL,
    `created_at`   TIMESTAMP      NOT NULL,
    `created_by`   INT            NOT NULL,
    `closed_at`    TIMESTAMP NULL,
    `closed_by`    INT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB;
