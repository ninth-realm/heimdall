CREATE TABLE `user` (
    `id` TEXT PRIMARY KEY NOT NULL,
    `first_name` TEXT NOT NULL,
    `last_name` TEXT NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TRIGGER [update_user_timestamp]
    AFTER UPDATE
    ON `user`
    FOR EACH ROW
BEGIN
    UPDATE `user` SET updated_at = CURRENT_TIMESTAMP WHERE id = old.id;
END;

CREATE TABLE `email` (
    `id` TEXT PRIMARY KEY NOT NULL,
    `user_id` TEXT NOT NULL,
    `email` TEXT NOT NULL UNIQUE,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`user_id`) REFERENCES `user` (`id`)
        ON DELETE CASCADE
);

CREATE TRIGGER [update_email_timestamp]
    AFTER UPDATE
    ON `email`
    FOR EACH ROW
BEGIN
    UPDATE `email` SET updated_at = CURRENT_TIMESTAMP WHERE id = old.id;
END;

CREATE TABLE `password` (
    `id` TEXT PRIMARY KEY NOT NULL,
    `user_id` TEXT NOT NULL UNIQUE,
    `hash` TEXT NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`user_id`) REFERENCES `user` (`id`)
        ON DELETE CASCADE
);

CREATE TRIGGER [update_password_timestamp]
    AFTER UPDATE
    ON `password`
    FOR EACH ROW
BEGIN
    UPDATE `password` SET updated_at = CURRENT_TIMESTAMP WHERE id = old.id;
END;
