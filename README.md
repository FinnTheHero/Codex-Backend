# Codex Backend

Backend for Codex - novel reader app.

# Details

Codex-Backend is built in `GoLang`, using `Gin` for server and `AWS-dynamoDB` for database.

It is deployed on `Heroku`.

## Development

### Setup

-   clone repo
    ```bash
    git clone -b Development https://github.com/FinnTheHero/Codex-Backend.git
    ```
-   create `.env` in root directory of the project
    ```bash
    cd "Codex-Backend" && touch .env
    ```
-   add AWS credentials to `.env` file
    ```bash
    AWS_ACCESS_KEY_ID = "" # Add your keys here
    AWS_SECRET_ACCESS_KEY = "" # Add your keys here
    AWS_REGION = "eu-central-1" # Change region if needed
    ```
