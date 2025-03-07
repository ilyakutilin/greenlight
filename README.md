# Greenlight

Greenlight is a JSON API for retrieving and managing information about movies. It provides functionality similar to the Open Movie Database API.

## Features

- Retrieve movie details
- Add, update, and delete movie entries
- Filtering, sorting, searching, pagination
- Request validation
- Authentication and authorization
- Rate limiting
- Logging and error handling
- Monitoring (metrics)

## Table of Contents

- [Installation](#installation)
- [Setup](#setup)
- [Running the Application](#running-the-application)
- [API Endpoints](#api-endpoints)
- [Configuration](#configuration)
- [Audit](#audit)
- [License](#license)

## Installation

To run Greenlight, you need the following dependencies:

- [Golang](https://go.dev/doc/install)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Migrate](https://github.com/golang-migrate/migrate)

For `Migrate` check the latest releases [here](https://github.com/golang-migrate/migrate/releases), then install like this:

```sh
curl -L https://github.com/golang-migrate/migrate/releases/download/<version>/<filename> | tar xvz
sudo mv migrate.linux-amd64 $GOPATH/bin/migrate
```

replacing the \<version\> and the \<filename\> with the actual version and filename.

- [Staticcheck](https://staticcheck.dev/) (used during audits)

```sh
go install honnef.co/go/tools/cmd/staticcheck@latest
which staticcheck
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

Create the `.envrc` file in the project root and open it:

```sh
touch .envrc
nano .envrc
```

Add the `GREENLIGHT_DB_DSN` environment variable and save the file:

```
export GREENLIGHT_DB_DSN='postgres://greenlight:pa55word@localhost/greenlight'
```

Adjust the actual DSN if needed as per the [instructions](https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING-URIS).

### 5. Run Database Migrations

While in the project folder:

```sh
migrate -path=./migrations -database=$GREENLIGHT_DB_DSN up
```

or

```sh
make db/migrations/up
```

## Running the Application

```sh
go run ./cmd/api
```

or

```sh
make run/api
```

## API Endpoints

| Method   | URL Pattern                 | Action                                          |
| -------- | --------------------------- | ----------------------------------------------- |
| `GET`    | `/v1/healthcheck`           | Show application health and version information |
| `GET`    | `/v1/movies`                | Show the details of all movies                  |
| `POST`   | `/v1/movies`                | Create a new movie                              |
| `GET`    | `/v1/movies/:id`            | Show the details of a specific movie            |
| `PATCH`  | `/v1/movies/:id`            | Update the details of a specific movie          |
| `DELETE` | `/v1/movies/:id`            | Delete a specific movie                         |
| `POST`   | `/v1/users`                 | Register a new user                             |
| `PUT`    | `/v1/users/activated`       | Activate a specific user                        |
| `PUT`    | `/v1/users/password`        | Update the password for a specific user         |
| `POST`   | `/v1/tokens/authentication` | Generate a new authentication token             |
| `POST`   | `/v1/tokens/password-reset` | Generate a new password-reset token             |
| `GET`    | `/debug/vars`               | Display application metrics                     |

## Configuration

The following flags can be used when launching the application:

| flag                  | values                               | default                |
| --------------------- | ------------------------------------ | ---------------------- |
| -port                 | integer                              | `4000`                 |
| -env                  | development \| staging \| production | `development`          |
| -db-dsn               | DSN URI                              | empty                  |
| -db-max-open-conns    | integer                              | `25`                   |
| -db-max-idle-conns    | integer                              | `25`                   |
| -db-max-idle-time     | %dm                                  | `15m`                  |
| -limiter-rps          | integer                              | `2`                    |
| -limiter-burst        | integer                              | `4`                    |
| -limiter-enabled      | true \| false                        | `true`                 |
| -smtp-host            | string                               | dev smtp host          |
| -smtp-port            | integer                              | `25`                   |
| -smtp-username        | string                               | dev smtp username      |
| -smtp-password        | string                               | dev smtp password      |
| -smtp-sender          | string                               | dev dummy sender email |
| -cors-trusted-origins | space-separated list of URLs         | empty                  |

## Audit

Make sure that the [Staticcheck](https://staticcheck.dev/) is installed before launching audits.

Launch audit:

```sh
make audit
```

Audit includes:

- [`go mod tidy`](https://go.dev/ref/mod#go-mod-tidy)

- [`go mod verify`](https://go.dev/ref/mod#go-mod-verify)

- [`go mod vendor`](https://go.dev/ref/mod#go-mod-vendor)

- [`go fmt ./...`](https://pkg.go.dev/cmd/go#hdr-Gofmt__reformat__package_sources)

- [`go vet ./...`](https://pkg.go.dev/cmd/vet)

- `go test -race -vet=off ./...` command to run all tests in the project directory. The `-race` flag enables Go’s race detector, which can help pick up certain classes of race conditions while tests are running.

- third-party [`staticcheck`](https://staticcheck.dev/docs/running-staticcheck/cli/) tool to carry out some additional static analysis checks.

==========

Additionally, running `make vendor` regularly (and mandatorily after each install).

## License

This project is licensed under the MIT License. See the [LICENSE](https://opensource.org/license/mit) file for details.
