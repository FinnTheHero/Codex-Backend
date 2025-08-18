# Codex Backend

Backend for Codex - novel reading platform.

## Details

Codex-Backend is built in `GoLang`, using `Gin` for server and ~AWS-dynamoDB~ firestore (moving to Heroku Postgres) for database.

It is deployed on `Heroku` (thats why the code is in api directory).

air config is outdated and not recommended. use [Run](Run guide instead)

## Run
run server:

    ```bash
    go run api/cmd/web/main.go
    ```

    ```bash
    go run api/cmd/worker/main.go
    ```

Both are needed

## Endpoints

3 Groups of endpoints: Client, Manage and User.

- Client is responsible for basic GET requests.
- Manage is responsible for Upload/Modification operations.
- User is responsible for user authentication, authorization and Registration (Delete is not yet implemented).

### Client: base path followed by request path
- `/all` - Get all novels
- `/:novel` - Get a novel by id
- `/:novel/:chapter` - Get chapter from novel using both ids
- `/:novel/all` - Get all chapters from novel using id
- `/:novel/chapter` - Get cursor paginated chapters from novel using id

    Options: limit (max 100), cursor (chapter index (integer)) and sort ("asc" || "desc").

    Defaults: limit=100, cursor=0, sort="desc"

### Manage: `/manage` followed by request path
- `/upload` - Upload novel
- `/:novel` - Update novel
- `/:novel/:chapter` - Update chapter

### User: `/user` followed by request path
- `/validate` - Validate user token
- `/login` - Login user
- `/register` - Register user
- `/logout` - Logout user
