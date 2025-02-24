# Codex Backend

Backend for Codex - novel reader app.

# Details

Codex-Backend is built in `GoLang`, using `Gin` for server and `AWS-dynamoDB` for database.

It is deployed on `Heroku`.

## Endpoints

### Group: `auth` - `/auth`

- GET - `/login` - Login user

    Middleware will check if user has a token, if its valid it will return user data.
    Otherwise it will be redirected to login handler where provided credentials will be used to login.
    Finally it will return token as a cookie.

    ```json
    // Request body example
    {
        "email": "email",
        "password": "password"
    }
    ```

- POST - `/register` - Register user
    ```json
    // Request body example
    {
        "username": "user",
        "password": "password",
        "email": "example@mail"
    }
    ```

### Group: `client` - `/`

- `GET` - `/all` - Get all novels
- `GET` - `/:novel` - Get novel by name
- `GET` - `/:novel/all` - Get all chapters in novel
- `GET` - `/:novel/:chapter` - Get specific chapter in novel

### Group: `admin` - `/admin`

- `POST` - `/novel` - Create novel

    ```json
    // Request body example
    {
        "title": "Novel",
        "author": "Author",
        "description": "Description",
        "creation_date": "2024-08-03T15:35:24.621148124+04:00",
        "upload_date": "2024-08-03T15:35:24.621148124+04:00",
        "update_date": "2024-08-03T15:35:24.621148124+04:00"
    }
    ```

- `POST` - `/:novel/chapter` - Create chapter to novel

    ```json
    // Request body example
    {
        "title": "Chapter",
        "author": "Author",
        "description": "Description",
        "creation_date": "2024-08-03T15:35:24.621148124+04:00",
        "upload_date": "2024-08-03T15:35:24.621148124+04:00",
        "update_date": "2024-08-03T15:35:24.621148124+04:00",
        "content": "Content"
    }
    ```

- `PUT` - `/:novel` - Update novel

    ```json
    // TODO
    ```

- `PUT` - `/:novel/:chapter` - Update chapter

    ```json
    // TODO
    ```

- `DELETE` - `/:novel` - Delete novel

    ```json
    // TODO
    ```

- `DELETE` - `/:novel/:chapter` - Delete chapter
    ```json
    // TODO
    ```

## Development

### Setup

- clone repo
    ```bash
    git clone -b Development https://github.com/FinnTheHero/Codex-Backend.git
    ```
- create `.env` in root directory of the project
    ```bash
    cd "Codex-Backend" && touch .env
    ```
- add AWS credentials to `.env` file
    ```bash
    AWS_ACCESS_KEY_ID = "" # Add your keys here
    AWS_SECRET_ACCESS_KEY = "" # Add your keys here
    AWS_REGION = "eu-central-1" # Change region if needed
    ```

### Run

- run server:
    ```bash
    go run api/cmd/codex/main.go
    ```
