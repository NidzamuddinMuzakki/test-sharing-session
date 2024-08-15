-- +goose Up
-- +goose StatementBegin
-- sharing-vision.log_posts definition
CREATE TABLE IF NOT EXISTS `log_posts`
(
    `id` int unsigned NOT NULL AUTO_INCREMENT ,
    `article_id` int NOT NULL DEFAULT 0,
    `data_before` TEXT NULL,
    `data_after` TEXT NULL,
    `category_status`varchar(100)  NOT NULL DEFAULT '',
    `created_date` timestamp    NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_date` timestamp    NULL DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `log_posts`;
-- +goose StatementEnd
