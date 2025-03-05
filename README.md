# Greenlight

Greenlight is a JSON API for retrieving and managing information about movies. It provides functionality similar to the Open Movie Database API.

## Features

- Retrieve movie details
- Add, update, and delete movie entries
- Filtering, sorting, searching, pagination
- Authentication and authorization
- Request validation

## Table of Contents

- [Installation](#installation)
- [Setup](#setup)
- [Running the Application](#running-the-application)
- [API Endpoints](#api-endpoints)
- [Configuration](#configuration)
- [License](#license)

## Installation

To run Greenlight, you need the following dependencies:

- [Golang](https://go.dev/doc/install)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Migrate](https://github.com/golang-migrate/migrate)

For `Migrate` check the latest releases [here](https://github.com/golang-migrate/migrate/releases), then install like this:

```sh
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz
mv migrate.linux-amd64 $GOPATH/bin/migrate
```

## Setup

### 1. Clone the Repository

```sh
git clone https://github.com/ilyakutilin/greenlight.git
cd greenlight
```

### 2. Install Dependencies

```sh
go mod tidy
```

### 3. Set up PostgreSQL

```sh
sudo -u postgres psql
```

```sql
CREATE DATABASE greenlight
        WITH ENCODING 'UTF8'
        LC_COLLATE = 'en_US.UTF-8'
        LC_CTYPE = 'en_US.UTF-8'
        TEMPLATE template0;

CREATE ROLE greenlight WITH LOGIN PASSWORD 'pa55word';

ALTER DATABASE greenlight OWNER TO greenlight;

CREATE EXTENSION IF NOT EXISTS citext;
```

Use [this web-based tool](https://pgtune.leopard.in.ua) to generate suggested `postgresql.conf` values based on your available system hardware.

### 4. Set the Environment Variable

Create a new `GREENLIGHT_DB_DSN` environment variable by adding the following line to your `$HOME/.profile` or `$HOME/.bashrc` or `$HOME/.zshrc` files:

```sh
export GREENLIGHT_DB_DSN='postgres://greenlight:pa55word@localhost/greenlight'
```

Adjust the actual DSN if needed as per the [instructions](https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING-URIS).

### 5. Run Database Migrations

While in the project folder:

```sh
migrate -path=./migrations -database=$GREENLIGHT_DB_DSN up
```

## Running the Application

```sh
go run ./cmd/api
```

## API Endpoints

| Method | URL Pattern               | Action                                          |
| ------ | ------------------------- | ----------------------------------------------- |
| GET    | /v1/healthcheck           | Show application health and version information |
| GET    | /v1/movies                | Show the details of all movies                  |
| POST   | /v1/movies                | Create a new movie                              |
| GET    | /v1/movies/:id            | Show the details of a specific movie            |
| PATCH  | /v1/movies/:id            | Update the details of a specific movie          |
| DELETE | /v1/movies/:id            | Delete a specific movie                         |
| POST   | /v1/users                 | Register a new user                             |
| PUT    | /v1/users/activated       | Activate a specific user                        |
| PUT    | /v1/users/password        | Update the password for a specific user         |
| POST   | /v1/tokens/authentication | Generate a new authentication token             |
| POST   | /v1/tokens/password-reset | Generate a new password-reset token             |
| GET    | /debug/vars               | Display application metrics                     |

## Configuration

The following flags can be used when launching the application:

| flag               | values                               | default           |
| ------------------ | ------------------------------------ | ----------------- |
| -port              | integer                              | 4000              |
| -env               | development \| staging \| production | development       |
| -db-dsn            | DSN URI                              | \<see DSN above\> |
| -db-max-open-conns | integer                              | 25                |
| -db-max-idle-conns | integer                              | 25                |
| -db-max-idle-time  | %dm                                  | 15m               |

## License

This project is licensed under the MIT License. See the [LICENSE](https://opensource.org/license/mit) file for details.
