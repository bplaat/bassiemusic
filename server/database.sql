CREATE TABLE `artists` (
    `id` BINARY(16) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    PRIMARY KEY (`id`)
);

CREATE TABLE `albums` (
    `id` BINARY(16) NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `released_at` DATE NOT NULL,
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
    FOREIGN KEY (`album_id`) REFERENCES `albums`(`id`),
    FOREIGN KEY (`artist_id`) REFERENCES `artists`(`id`)
);

CREATE TABLE `tracks` (
    `id` BINARY(16) NOT NULL,
    `album_id` BINARY(16) NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `disk` INT UNSIGNED NOT NULL,
    `position` INT UNSIGNED NOT NULL,
    `duration` INT UNSIGNED NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`album_id`) REFERENCES `albums`(`id`)
);

CREATE TABLE `track_artist` (
    `id` BINARY(16) NOT NULL,
    `track_id` BINARY(16) NOT NULL,
    `artist_id` BINARY(16) NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`track_id`) REFERENCES `tracks`(`id`),
    FOREIGN KEY (`artist_id`) REFERENCES `artists`(`id`)
);
