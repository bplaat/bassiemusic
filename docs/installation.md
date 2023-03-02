[&laquo; Back to the README.md](../README.md)

# Installation Documentation

## Windows

### Database
- Download and install the latest stable [MariaDB](https://mariadb.org/download/)
- In the setup also install HeidiSQL
- Login as root to your database with HeidiSQL
- Create BassieMusic user and database
    ```sql
    CREATE USER 'bassiemusic'@'localhost' IDENTIFIED BY 'bassiemusic';
    CREATE DATABASE `bassiemusic` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
    GRANT ALL PRIVILEGES ON `bassiemusic`.* TO 'bassiemusic'@'localhost';
    FLUSH PRIVILEGES;
    ```
- Change database with `use bassiemusic;`
- Create MariaDB UUID_TO_BIN and BIN_TO_UUID [pollyfills](https://gist.github.com/bplaat/1d8d1bba135c726178ebdfc9df08e2ca)
- Create the tables in [`server/database.sql`](../server/database.sql)

### Server
- Download and install [Go compiler](https://go.dev/dl/) 1.18 or higher
- Go into the server folder with `cd server`
- Build the server `go build`
- Copy the `.env.example` file to `.env`
- Start the server `./bassiemusic serve`

### Web frontend
- Download and install [Node.js](https://nodejs.org/en/) 18 or higher
- Go into the web folder with `cd web`
- Install NPM dependencies with `npm install`
- Copy the `.env.example` file to `.env`
- Start the dev server with `npm run dev`

### Youtube downloader
- Download and install [yt-dlp](https://github.com/yt-dlp/yt-dlp#installation)
- Move `yt-dlp.exe` executable to folder that is in your path
- Download and install [FFmpeg](https://www.gyan.dev/ffmpeg/builds/)
- Add the `ffmpeg/bin/` folder to your path

### You're done, now login as admin
- Now run the server and the web frontend and go to http://localhost:5173/
- Login as admin with the username: **admin**, password: **admin**
- Go to the Admin Downloader page, search and download an album you like
- It will show up in the global Albums page and you can listen to it

## macOS

### Database
- Install [Homebrew](https://brew.sh/)
- Install the latest stable MySQL with `brew install mysql`
- Login as root to your database with `sudo mysql`
- Create BassieMusic user and database
    ```sql
    CREATE USER 'bassiemusic'@'localhost' IDENTIFIED BY 'bassiemusic';
    CREATE DATABASE `bassiemusic` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
    GRANT ALL PRIVILEGES ON `bassiemusic`.* TO 'bassiemusic'@'localhost';
    FLUSH PRIVILEGES;
    ```
- Change database with `use bassiemusic;`
- Create the tables in [`server/database.sql`](../server/database.sql)

### Server
- Install the latest Go compiler with `brew install go`
- Go into the server folder with `cd server`
- Build the server `go build`
- Copy the `.env.example` file to `.env`
- Start the server `./bassiemusic serve`

### Web frontend
- Install the latest Node.js with `brew install node`
- Go into the web folder with `cd web`
- Install NPM dependencies with `npm install`
- Copy the `.env.example` file to `.env`
- Start the dev server with `npm run dev`

### Youtube downloader
- Install yt-dlp and ffmpeg with `brew install yt-dlp ffmpeg`

### You're done, now login as admin
- Now run the server and the web frontend and go to http://localhost:5173/
- Login as admin with the username: **admin**, password: **admin**
- Go to the Admin Downloader page, search and download an album you like
- It will show up in the global Albums page and you can listen to it

## Linux (Ubuntu based distro's)

### Database
- Install the latest stable MySQL with `sudo apt install mysql-server`
- Login as root to your database with `sudo mysql`
- Create BassieMusic user and database
    ```sql
    CREATE USER 'bassiemusic'@'localhost' IDENTIFIED BY 'bassiemusic';
    CREATE DATABASE `bassiemusic` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
    GRANT ALL PRIVILEGES ON `bassiemusic`.* TO 'bassiemusic'@'localhost';
    FLUSH PRIVILEGES;
    ```
- Change database with `use bassiemusic;`
- Create the tables in [`server/database.sql`](../server/database.sql)

### Server
- Download and install [Go compiler](https://go.dev/dl/) 1.18 or higher
- Add the `go/bin` folder to your path
- Go into the server folder with `cd server`
- Build the server `go build`
- Copy the `.env.example` file to `.env`
- Start the server `./bassiemusic serve`

### Web frontend
- Install the [Node Version Manager](https://github.com/nvm-sh/nvm#install--update-script)
- Restart terminal and install Node with `nvm install node`
- Go into the web folder with `cd web`
- Install NPM dependencies with `npm install`
- Copy the `.env.example` file to `.env`
- Start the dev server with `npm run dev`

### Youtube downloader
- Download and install [yt-dlp](https://github.com/yt-dlp/yt-dlp#installation)
- Move `yt-dlp` executable to folder that is in your path
- Install ffmpeg with `sudo apt install ffmpeg`

### You're done, now login as admin
- Now run the server and the web frontend and go to http://localhost:5173/
- Login as admin with the username: **admin**, password: **admin**
- Go to the Admin Downloader page, search and download an album you like
- It will show up in the global Albums page and you can listen to it
