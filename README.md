# TODO-List managment service

A containerized web server with database for todo list managment.

# How to use
Setup: 

- Create a `.env` file in the repo directory using env.example. Make sure you've set all the variables. Even tho the app itself has it's own default values in case of missing envs, docker-compose however needs all required values.
If you set an `APP_STATIC_TOKEN`, server will start with a token-auth middleware. Make sure to add `token` header with the value you set up in your .env file for each request.


- Run `make run`. It will build docker containers and expose an API at `localhost:8080`. 

Usage:
to be described in swagger.



# Running tests

- Run `make db`
- In separate terminal window run `make test`
- Run `make stop` after tests finish
