# TODO-List managment service

A containerized web server with database for todo list managment.

# How to use
Setup: 

- Create a `.env` file in the repo directory using env.example. Make sure you've set all the variables. Even tho the app itself has it's own default values in case of missing envs, docker-compose however needs all required values.
Set `APP_STATIC_TOKEN`, server will start with a token-auth middleware. Default token is `test_token`".

- Run `make run`. It will build docker containers and expose an API at `localhost:8080`. 

Usage:
Swagger UI will be accessible at `http://localhost:8081/swagger/index.html#`

Click `authorize` button and specify your token that you've set in .env (`APP_STATIC_TOKEN`). If you haven't specified a token, type anything and click 'Authorize' to sart using Swagger.


# Running tests

- Run `make db`
- In separate terminal window run `make test`
- Run `make stop` after tests finish
