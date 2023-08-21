CREATE TABLE `session` (
    `token` TEXT PRIMARY KEY NOT NULL,
    `user_id` TEXT NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `expires_at` DATETIME NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES `user` (`id`)
        ON DELETE CASCADE
);

