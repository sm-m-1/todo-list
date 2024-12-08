To do app in go lang with several features
- authentication based on username and password
- session based authentication with postgress DB. Supports distributed architecture.
- basic crud operations




steps

0. run the go server with `go run main.go`

1. create a user

curl --location 'http://localhost:8080/register' \
--header 'Content-Type: application/json' \
--data '{"username": "user3", "password": "p"}'


2. login as the user
curl --location 'http://localhost:8080/login' \
--header 'Content-Type: application/json' \
--data '{"username": "user3", "password": "p"}'

A response header will be returned which is needed for the next step.

Set-Cookie: session=3NgznVMsOgGIYe3ASJauA6qGetmVp0_zsOFvqVbYkz0; Path=/; Expires=Mon, 09 Dec 2024 03:58:08 GMT; Max-Age=86400; HttpOnly; SameSite=Lax

3. get all the todo list
curl --location 'http://localhost:8080/todos' \
--header 'Cookie: session=9Ph9aqvhNdJRPvb0QlOZBRq7M_F_nDLJkraiUCp2chk; session=3NgznVMsOgGIYe3ASJauA6qGetmVp0_zsOFvqVbYkz0' \
--header 'Content-Type: application/json'

