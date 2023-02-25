# gin-gorm-goth
A Boilerplate in Go with Gin, GORM, and Goth for creating a basic Oauth user authentication 

## Description:
This code provides a simple implementation of user authentication routes in a Go web application using the Gin web framework and the GORM ORM library. The "authRoutes.go" file defines the authentication routes and their actions. The "/auth" group of routes includes a simple message route and two routes for initiating and completing the authentication process. The "goth" package is used to handle the Discord authentication provider. Upon completion of the authentication process, the user's details are stored in the database, a session token is created and returned in an HTTP-only cookie, and a session ID is also returned. If any errors occur, appropriate error messages are returned with their corresponding HTTP status codes. This boilerplate provides a solid foundation for building user authentication systems in Go web applications.
