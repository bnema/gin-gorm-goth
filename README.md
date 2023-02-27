# gin-gorm-goth
A Boilerplate in Go with Gin, GORM, and Goth for creating a basic Oauth user authentication (with Discord as provider)

## Description

This is a boilerplate codebase created for learning purposes and heavily inspired by NextAuth. It provides basic RESTful API endpoints for user authentication and authorization, using Gin as web/http server library, GORM as the ORM, and Goth/Gothic as the Oauth2 client library. The API supports CRUD operations for managing user accounts,authentication sessions and posts.

## Routes

The `/auth/login` route is used to initiate the authentication process, where it sets the authentication providers (in this case, it uses the Discord provider) and starts the authentication process using the "gothic" package.

The `/auth/callback` route is used to complete the authentication process, where it receives the user's authentication details, creates a new user and account in the database, creates a session token, and returns it in a HTTP-only cookie along with the user's session ID.

The `/auth/logout` route is used to destroy the user's session and remove the session cookie from the user's browser.

The `/blog` route return all the Posts in the database as JSON responses.

The `/blog/:title` route returns a single Post in the database as a JSON response.

The `/blog/admin` route is used to test the authentication process and the authMiddleware. It is only accessible to authenticated users AND if isAdmin is set to true. It returns a JSON response indicating that it's the admin route.

The `/blog/admin/createpost` route is used to create a new Post in the database. 

The `/blog/admin/updatepost` route is used to update a Post in the database.

The `/blog/admin/deletepost` route is used to delete a Post in the database.

Bonus: The `/benchmark` route is used to test the performance of the application. (A snippet of the http client tool is provided in `/benchmarkClient/benchmarkClient.go`)


## Env variables

```
GIN_MODE=debug
DOMAIN=localhost
PORT=3000
DATABASE_URL=postgresql://yourname:yourpasswordE@yourPGserver:PORT/yourdb
SESSION_SECRET=secret
JWT_SECRET=secret
DISCORD_CLIENT_ID=client_id
DISCORD_CLIENT_SECRET=client_secret
AUTH_CALLBACK_URL=http://localhost:3000/auth/callback

```

## Installation

Clone the repository

- Set up the environment variables by creating a .env file in the root directory and filling in the required variables (see the section on Env variables above).
- Run `go mod download` to download the required packages.
- Run `go run main.go` to start the application.