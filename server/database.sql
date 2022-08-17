CREATE TABLE `users` (
    `id` BINARY(16) NOT NULL,
    `firstname` VARCHAR(255) NOT NULL,
    `lastname` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `password` VARCHAR(255) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    PRIMARY KEY (`id`)
);

INSERT INTO `users` (`id`, `firstname`, `lastname`, `email`, `password`) VALUES (UUID_TO_BIN(UUID()), 'Bastiaan', 'van der Plaat', 'bastiaan.v.d.plaat@gmail.com', '$2y$10$dYo9.ABRMCdxo6fDdlWC7uL7MHbaD7sQFDPtVh97feHVnsQEgJh9m');

CREATE TABLE `sessions` (
    `id` BINARY(16) NOT NULL,
    `user_id` BINARY(16) NOT NULL,
    `token` VARCHAR(255) NOT NULL,
    `os` VARCHAR(32) NULL,
    `platform` VARCHAR(32) NULL,
    `version` VARCHAR(32) NULL,
    `expires_at` TIMESTAMP NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);

CREATE TABLE `artists` (
    `id` BINARY(16) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `deezer_id` BIGINT UNSIGNED NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    PRIMARY KEY (`id`)
);

CREATE TABLE `albums` (
    `id` BINARY(16) NOT NULL,
    `type` TINYINT NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `released_at` DATE NOT NULL,
    `explicit` BOOLEAN NOT NULL,
    `deezer_id` BIGINT UNSIGNED NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    PRIMARY KEY (`id`)
);

CREATE TABLE `album_artist` (
    `id` BINARY(16) NOT NULL,
    `album_id` BINARY(16) NOT NULL,
    `artist_id` BINARY(16) NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`album_id`) REFERENCES `albums`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`artist_id`) REFERENCES `artists`(`id`) ON DELETE CASCADE
);

CREATE TABLE `genres` (
    `id` BINARY(16) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `deezer_id` BIGINT UNSIGNED NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    PRIMARY KEY (`id`)
);

CREATE TABLE `album_genre` (
    `id` BINARY(16) NOT NULL,
    `album_id` BINARY(16) NOT NULL,
    `genre_id` BINARY(16) NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`album_id`) REFERENCES `albums`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`genre_id`) REFERENCES `genres`(`id`) ON DELETE CASCADE
);

CREATE TABLE `tracks` (
    `id` BINARY(16) NOT NULL,
    `album_id` BINARY(16) NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `disk` INT UNSIGNED NOT NULL,
    `position` INT UNSIGNED NOT NULL,
    `duration` FLOAT NOT NULL,
    `explicit` BOOLEAN NOT NULL,
    `deezer_id` BIGINT UNSIGNED NOT NULL,
    `youtube_id` VARCHAR(16) NOT NULL,
    `plays` BIGINT UNSIGNED NOT NULL DEFAULT 0,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`album_id`) REFERENCES `albums`(`id`) ON DELETE CASCADE
);

CREATE TABLE `track_artist` (
    `id` BINARY(16) NOT NULL,
    `track_id` BINARY(16) NOT NULL,
    `artist_id` BINARY(16) NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`track_id`) REFERENCES `tracks`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`artist_id`) REFERENCES `artists`(`id`) ON DELETE CASCADE
);

CREATE TABLE `track_play` (
    `id` BINARY(16) NOT NULL,
    `user_id` BINARY(16) NOT NULL,
    `track_id` BINARY(16) NOT NULL,
    `position` FLOAT NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`track_id`) REFERENCES `tracks`(`id`) ON DELETE CASCADE
);

-- artist_likes
-- album_likes
-- track_likes
