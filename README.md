# xm-golang-exercise
Test task for Go devs from XM

## Prerequisites
Newer Go release (1.23 or up) and Docker

# Services

## Auth Service
Auth Server with two endpoints for registering and retrieving a JWT. Stores user data in Mongo.
### Endpoints
*POST /register*
Req body example
```json
{
    "username": "test_user",
    "password": "test_password"
}
```
*POST /login*
Req body example
```json
{
    "username": "test_user",
    "password": "test_password"
}
```

## Companies Service
Has CRUD endpoints for managing Companies stored in a Postgres instance. *runMigrations.sh* script exists for running migrations. Data mutating endpoints check for valid JWT (using middleware) that should be fetched from *Auth Service*. Each operation also fires an event that results in sending a message to Kafka.
*POST /company*
Req body example
```json
{
  "name": "Skynet",
  "description": "A leading technology company specializing in AI solutions.",
  "amount_of_employees": 100,
  "registered": true,
  "type": "Corporations"
}
```
*PATCH /company/:id*
Req body example
```json
{
  "name": "Skynet",
  "description": "A leading technology company specializing in AI solutions.",
  "amount_of_employees": 100,
  "registered": true,
  "type": "Corporations"
}
```
*DELETE /company/:id*
*GET /company/:id*

## Event Consumer Service
A very simple service that subscribes to Kafka Topic and logs received messages

All services, except Consumer, have scripts for running tests, building, linting, mock generation. All services have Alpine build Dockerfiles and docker-compose files for spinning up service dependencies like databases and Kafka. You can start each of the services separately and use the provided endpoints. Just for testing it's easier to start docker-compose-local.yml.

## Running all services
*docker-compose-local.yml* will start all the required databases, run migrations on Postgres, start Kafka and the three services. If starting for the first time give it a minute or two for the migration container to finish before hitting the endpoints.
