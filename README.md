# Codex Backend

Backend for Codex - novel reading platform.

## Details

Codex-Backend is built in `GoLang`, using `Gin` for server and ~AWS-dynamoDB~ ~firestore (moving to Heroku Postgres)~ PostgreSQL for database.

~It is deployed on `Heroku` (thats why the code is in api directory).~

Deployed on personal server in Docker.

## Run
run server:

```bash
GIN_MODE=debug go run api/cmd/web/main.go
```

```bash
GIN_MODE=debug go run api/cmd/worker/main.go
```

Both are needed

## Endpoints

5 Groups of endpoints: Client, Manage, User, Validate and Health.

- Client is responsible for basic GET requests.
- Manage is responsible for Upload/Modification/Delete operations.
- User is responsible for user authentication, authorization and Registration (Delete is not yet implemented).
- Validate is responsible for validating user tokens.
- Health is responsible for checking the health of the server.

### Client: `/api` followed by request path
- GET `/all` - Get all novels
- GET `/:novel` - Get a novel by id
- GET `/:novel/:chapter` - Get chapter from novel using both ids
- GET `/:novel/all` - Get all chapters from novel using id
- GET `/:novel/chapter` - Get cursor paginated chapters from novel using id

Pagination querries:
- `?limit=100` - Limit the number of results returned, Max = 200, Min = 1
- `?cursor=""` - Encoded offset, will be handled automatically
- `?sort="asc"` - Sort order, asc or desc

### Manage: `/api/manage` followed by request path
- POST `/epub` Create Novel/Chapters from epub file.

- POST `/create/novel` Create Novel
- POST `/create/chapter` Create Chapter

- PUT `/update/novel` - Update novel
- PUT `/update/chapter` - Update chapter

- DELETE `/delete/novel` - Delete novel
- DELETE `/delete/chapter` - Delete chapter

### User: `/api/user` followed by request path
- POST `/login` - Login user
- POST `/logout` - Logout user
- POST `/register` - Register user

### Validate: `/api/validate` followed by request path
- GET `/` - Validate user token

## Health: `/health` followed by request path
- GET `/` - Check health of server
For now this does nothing, but will be used to check the health of Docker image.
