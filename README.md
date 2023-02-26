# gin-gorm-goth
A Boilerplate in Go with Gin, GORM, and Goth for creating a basic Oauth user authentication (with Discord as provider)

## Description

This is a boilerplate code for creating routes related to user authentication in a Go web application using the Gin framework and the GORM ORM library. 

The code defines a group of routes under the "/auth" prefix that can be accessed by users to authenticate themselves.

The first route simply returns a JSON message indicating that it's the authentication route.

The "/login" route is used to initiate the authentication process, where it sets the authentication providers (in this case, it uses the Discord provider) and starts the authentication process using the "gothic" package.

The "/callback" route is used to complete the authentication process, where it receives the user's authentication details, creates a new user and account in the database, creates a session token, and returns it in a HTTP-only cookie along with the user's session ID.

The "/logout" route is used to destroy the user's session and remove the session cookie from the user's browser.

The "/blog" route return all the Posts in the database as JSON responses.

The "/blog/:id" route returns a single Post in the database as a JSON response.

The "/blog/admin" route is used to test the authentication process and the authMiddleware. It is only accessible to authenticated users AND if isAdmin is set to true. It returns a JSON response indicating that it's the admin route.

The "/blog/admin/createpost" route is used to create a new Post in the database. 

The "/blog/admin/updatepost" route is used to update a Post in the database.

The "/blog/admin/deletepost" route is used to delete a Post in the database.


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
