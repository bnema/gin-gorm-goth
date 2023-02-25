# gin-gorm-goth
A Boilerplate in Go with Gin, GORM, and Goth for creating a basic Oauth user authentication 

## Description

This is a boilerplate code for creating routes related to user authentication in a Go web application using the Gin framework and the GORM ORM library. 

The code defines a group of routes under the "/auth" prefix that can be accessed by users to authenticate themselves.

The first route simply returns a JSON message indicating that it's the authentication route.

The "/login" route is used to initiate the authentication process, where it sets the authentication providers (in this case, it uses the Discord provider) and starts the authentication process using the "gothic" package.

The "/callback" route is used to complete the authentication process, where it receives the user's authentication details, creates a new user and account in the database, creates a session token, and returns it in a HTTP-only cookie along with the user's session ID. The cookie is set to expire in 7 days. If any errors occur during these processes, appropriate error messages are returned as JSON responses with their corresponding HTTP status codes.

## WIP

- Check existing params
