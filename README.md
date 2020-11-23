# session-example
This is a minimal example of an authentication system using [Fiber](https://github.com/gofiber/fiber), Prisma's [Go Client](https://github.com/prisma/prisma-client-go) and Fiber's [Session middleware](https://github.com/gofiber/fiber/tree/master/middleware/session).

## Directory structure
* `main.go` - code
* `schema.prisma` - prisma schema
* `migrations` - prisma auto generated migrations
* `db` - prisma auto generated db package
* `fiber.sqlite3` - session storage
* `data.db` - prisma user storage

## Routes
By default the application runs on port 3000.


### POST `/signup` (signs up a user)
*Body:* 
```json
{
    "username": string,
    "password": string
}
```

### POST `/login` (logs in a user)
*Body:* 
```json
{
    "username": string,
    "password": string
}
```

*Returns* HTTP 201 

### GET `/protected` (returns a users ID)

*Returns* HTTP 200, user ID

