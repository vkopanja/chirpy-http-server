# Chirpy REST API—boot.dev HTTP servers course

Go project created for
the [boot.dev](https://boot.dev) - [HTTP Servers](https://www.boot.dev/lessons/50f37da8-72c0-4860-a7d1-17e4bda5c243)
course.

This is a simple REST API created during the guided HTTP Servers course from boot.dev

## Quickstart

You need to have Go and PostgreSQL installed on your machine.

### Installing Go

#### macOS (with Homebrew [brew.sh](https://brew.sh/))

```
brew install go
```

#### Linux

#### Ubuntu/Debian

```
sudo apt update && sudo apt install golang-go
```

#### Fedora

```
sudo dnf install golang
```

#### Arch

```
sudo pacman -S go
```

### Verify Go is working

```
go version
```

You will also need latest [goose](https://github.com/pressly/goose) for migrations
and [sqlc](https://github.com/sqlc-dev/sqlc) for
generating SQL code.

Easiest way to install `goose` is to use Go itself: `go install github.com/pressly/goose/v3/cmd/goose@latest`

Same thing for `sqlc`, just install it with: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`

## Usage

You will need to create a `.env` file with the following variables:

```dotenv
DB_URL=postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable
PLATFORM=dev
SECRET=zC+Lder3Sj859yD1K6F1eYqo9dKRf+/0HtsmxRIX
POLKA_KEY=f271c81ff7084ee5b99a5091b42d486e
```

You can update all the variables in `.env` to your needs/setup.

For usage, refer to the Swagger documentation
at [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)