-- +goose Up
-- +goose StatementBegin
-- insurance_cofi.partner definition
CREATE TABLE IF NOT EXISTS `posts`
(
    `id` int unsigned NOT NULL AUTO_INCREMENT ,
    `title`       varchar(200)  NOT NULL DEFAULT '',
    `content`       TEXT NOT NULL DEFAULT '',
    `category`varchar(100)  NOT NULL DEFAULT '',
    `status` varchar(100) NOT NULL DEFAULT '',
    `created_date` timestamp    NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_date` timestamp    NULL DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `posts`;
-- +goose StatementEnd
