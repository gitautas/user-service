DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
  `id`           int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_id`      char(36) CHARACTER SET ascii,
  `first_name`   varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `last_name`    varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `nickname`     varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `password`     varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `email`        varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `country`      varchar(3) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at`   timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`   timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
