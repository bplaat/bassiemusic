-- BassieMusic database

-- Users
CREATE TABLE `users` (
    `id` BINARY(16) NOT NULL,
    `username` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `password` VARCHAR(255) NOT NULL,
    `avatar` BINARY(16) NULL,
    `allow_explicit` TINYINT(1) UNSIGNED NOT NULL,
    `role` TINYINT UNSIGNED NOT NULL,
    `language` CHAR(2) NOT NULL,
    `theme` TINYINT UNSIGNED NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE (`username`),
    UNIQUE (`email`)
);

INSERT INTO `users` (`id`, `username`, `email`, `password`, `allow_explicit`, `role`, `language`, `theme`) VALUES
    (UUID_TO_BIN(UUID()), 'admin', 'admin@plaatsoft.nl', '$2a$10$GwDKz/4HjEklaq3FtdMYo.p3ildTU36iX1.29rdDRIIi9qgIlT7n2', 1, 1, 'en', 0);

-- Sessions
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

-- Artists
CREATE TABLE `artists` (
    `id` BINARY(16) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `synced` TINYINT(1) UNSIGNED NOT NULL,
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

-- Albums
CREATE TABLE `albums` (
    `id` BINARY(16) NOT NULL,
    `type` TINYINT UNSIGNED NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `released_at` DATE NOT NULL,
    `explicit` TINYINT(1) UNSIGNED NOT NULL,
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

-- Genres
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

-- Tracks
CREATE TABLE `tracks` (
    `id` BINARY(16) NOT NULL,
    `album_id` BINARY(16) NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `disk` INT UNSIGNED NOT NULL,
    `position` INT UNSIGNED NOT NULL,
    `duration` FLOAT NOT NULL,
    `explicit` TINYINT(1) UNSIGNED NOT NULL,
    `deezer_id` BIGINT UNSIGNED NOT NULL,
    `youtube_id` VARCHAR(16) NULL,
    `plays` BIGINT UNSIGNED NOT NULL,
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

-- Playlists
CREATE TABLE `playlists` (
    `id` BINARY(16) NOT NULL,
    `user_id` BINARY(16) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `image` BINARY(16) NULL,
    `public` TINYINT(1) UNSIGNED NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
);

CREATE TABLE `playlist_track` (
    `id` BINARY(16) NOT NULL,
    `playlist_id` BINARY(16) NOT NULL,
    `position` INT UNSIGNED NOT NULL,
    `track_id` BINARY(16) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`playlist_id`) REFERENCES `playlists`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`track_id`) REFERENCES `tracks`(`id`) ON DELETE CASCADE
);

CREATE TABLE `playlist_likes` (
    `id` BINARY(16) NOT NULL,
    `playlist_id` BINARY(16) NOT NULL,
    `user_id` BINARY(16) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`playlist_id`) REFERENCES `playlists`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);

-- Download tasks
CREATE TABLE `download_tasks` (
    `id` BINARY(16) NOT NULL,
    `type` TINYINT UNSIGNED NOT NULL,
    `deezer_id` BIGINT UNSIGNED NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
);
