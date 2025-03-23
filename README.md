# chirpy

chirpy is a RESTful api that manages users and their `chirps` with built in authentication and authorization,
built with the golang's `http` package

## requirements
to run this project you'll need:
- `go` with version at least `1.24`
- `.env` file with neccessary values set in it
    - `DB_URL` a url to your chirpy database it should look something like `postgres://user:password@localhost:5432/chirpy`
    - `JWT_SECRET` a secret string that is going to be used to generate and validate JWT tokens
    - `POLKA_KEY` a webhook api key
- `postgresql` database running on your local machine, or somwhere remote but remember to set the `DB_URL` appropriately

with all setup you can just `go run .` in root directory of the project, the app should print on what `port` is the server starting

to run tests:
- run `go test ./...` in the root directory of the project

## api

### /api/chirps

- `GET /api/chirps` displays all the chirps sorted by creation date in ascending order
- `GET /api/chirps/{id}` displays chirp by id
- `GET /api/chirps?author_id={id}` displays chirps by author id
- `GET /api/chirps?sort=desc` displays all the chirps sorted by creation date in descending order
- `GET /api/chirps?author_id={id}?sort={sorting}` displays chirps by author id sorted in given order
- `POST /api/chirps` creates a new chirp for an authorized user
- `DELETE /api/chirps/{id}` deletes a chirp by id for an authorized user

requests and responses used by `/api/chirp`

```
    type createChirpRQ struct {
        Body string `json:"body"`
    }

    type chirpRes struct {
        Id        string    `json:"id"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
        Body      string    `json:"body"`
        UserId    string    `json:"user_id"`
    }

```

### /api/users

- `POST /api/users` creates a new user with provided email and password, the password is hashed before storing
- `PUT /api/users` updates the users email and password

request and responses used by `/api/users`

```
    type userRQ struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    type userRes struct {
        Id           string    `json:"id"`
        CreatedAt    time.Time `json:"created_at"`
        UpdatedAt    time.Time `json:"updated_at"`
        Email        string    `json:"email"`
        IsChirpyRed  bool      `json:"is_chirpy_red"`
        Token        string    `json:"token,omitempty"`
        RefreshToken string    `json:"refresh_token,omitempty"`
    }
```

### auth api

- `POST /api/login` logs a user in, returning a token and a refresh token in response
- `POST /api/refresh` refreshes the JWT token provided a refresh token
- `POST /api/revoke` revokes a refresh token

### /admin/

- `GET /admin/metrics` returns a HTML with server hits value
- `POST /admin/reset` resets the database
- `GET /admin/users` returns all the users
  
### /api/healthz

- `GET /api/healthz` returns the status of the service
