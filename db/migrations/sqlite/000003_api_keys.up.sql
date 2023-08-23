CREATE TABLE `api_key` (
    `id` TEXT PRIMARY KEY NOT NULL,
    `client_id` TEXT NOT NULL,
    `description` TEXT NULL,
    `prefix` TEXT NOT NULL UNIQUE,
    `hash` TEXT NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `client` (`id`)
        ON DELETE CASCADE
);

CREATE TRIGGER [update_api_key_timestamp]
    AFTER UPDATE
    ON `api_key`
    FOR EACH ROW
BEGIN
    UPDATE `api_key` SET updated_at = CURRENT_TIMESTAMP WHERE id = old.id;
END;
