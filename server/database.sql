-- BassieMusic database

-- Create BassieMusic MySQL user
-- CREATE USER 'bassiemusic'@'localhost' IDENTIFIED BY 'bassiemusic';
-- CREATE DATABASE `bassiemusic` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- GRANT ALL PRIVILEGES ON `bassiemusic`.* TO 'bassiemusic'@'localhost';
-- FLUSH PRIVILEGES;

-- MariaDB UUID_TO_BIN and BIN_TO_UUID pollyfills:
-- https://gist.github.com/bplaat/1d8d1bba135c726178ebdfc9df08e2ca

-- Tables
CREATE TABLE `users` (
    `id` BINARY(16) NOT NULL,
    `username` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `password` VARCHAR(255) NOT NULL,
    `avatar` BINARY(16) NULL,
    `role` TINYINT UNSIGNED NOT NULL,
    `theme` TINYINT UNSIGNED NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE (`username`),
    UNIQUE (`email`)
);

INSERT INTO `users` (`id`, `username`, `email`, `password`, `role`, `theme`) VALUES
    (UUID_TO_BIN(UUID()), 'bplaat', 'bastiaan.v.d.plaat@gmail.com', '$2a$10$21hEKLKeYntMkANwm.RCludVDbMU12PRqmc.k6febZUkJHNDoLEAq', 1, 0),
    (UUID_TO_BIN(UUID()), 'lplaat', 'leonard.van.der.plaat@gmail.com', '$2a$10$21hEKLKeYntMkANwm.RCludVDbMU12PRqmc.k6febZUkJHNDoLEAq', 1, 0);

CREATE TABLE `sessions` (
    `id` BINARY(16) NOT NULL,
    `user_id` BINARY(16) NOT NULL,
    `token` VARCHAR(255) NOT NULL,
    `ip` VARCHAR(48) NOT NULL,
    `ip_latitude` DECIMAL(10, 8) NULL,
    `ip_longitude` DECIMAL(11, 8) NULL,
    `ip_country` CHAR(2) NULL,
    `ip_city` VARCHAR(255) NULL,
    `client_os` VARCHAR(32) NULL,
    `client_name` VARCHAR(32) NULL,
    `client_version` VARCHAR(32) NULL,
    `expires_at` TIMESTAMP NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE,
    UNIQUE (`token`)
);

CREATE TABLE `artists` (
    `id` BINARY(16) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `deezer_id` BIGINT UNSIGNED NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
);

CREATE TABLE `artist_likes` (
    `id` BINARY(16) NOT NULL,
    `artist_id` BINARY(16) NOT NULL,
    `user_id` BINARY(16) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`artist_id`) REFERENCES `artists`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);

CREATE TABLE `albums` (
    `id` BINARY(16) NOT NULL,
    `type` TINYINT UNSIGNED NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `released_at` DATE NOT NULL,
    `explicit` BOOLEAN NOT NULL,
    `deezer_id` BIGINT UNSIGNED NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
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

CREATE TABLE `album_likes` (
    `id` BINARY(16) NOT NULL,
    `album_id` BINARY(16) NOT NULL,
    `user_id` BINARY(16) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`album_id`) REFERENCES `albums`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);

CREATE TABLE `genres` (
    `id` BINARY(16) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `deezer_id` BIGINT UNSIGNED NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
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

CREATE TABLE `track_likes` (
    `id` BINARY(16) NOT NULL,
    `track_id` BINARY(16) NOT NULL,
    `user_id` BINARY(16) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`track_id`) REFERENCES `tracks`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);

CREATE TABLE `track_plays` (
    `id` BINARY(16) NOT NULL,
    `track_id` BINARY(16) NOT NULL,
    `user_id` BINARY(16) NOT NULL,
    `position` FLOAT NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`track_id`) REFERENCES `tracks`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);

CREATE TABLE `download_tasks` (
    `id` BINARY(16) NOT NULL,
    `type` TINYINT UNSIGNED NOT NULL,
    `deezer_id` BIGINT UNSIGNED NULL,
    `singles` BOOLEAN NOT NULL DEFAULT 0,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
);
