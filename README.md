# Nexlab core backend

## Project Structure

### Backend
- `services/controller`: Controller service [http://localhost:8080](http://localhost:8080)
- `services/auth`: Auth webhook service

## Prerequisites

- Hasura GraphQL Engine 2.0
- Golang 1.15+
- Docker + docker-compose

## Database design and migration

Use Hasura 2.0 CLI: https://docs.hasura.io/1.0/graphql/manual/hasura-cli/install-hasura-cli.html#install

```sh
# rename the CLI to hasura2 to avoid conflict if using Hasura v1 in parallel
sudo wget https://github.com/hasura/graphql-engine/releases/download/v2.0.9/cli-hasura-linux-amd64 -O /usr/local/bin/hasura2
sudo chmod +x /usr/local/bin/hasura2
```

- Design

```
hasura2 console --admin-secret hasura
```

- Migrate: 

```
hasura2 migrate apply --all-databases --admin-secret hasura
hasura2 metadata apply --admin-secret hasura
```

## How to Run

### Full stack

Copy `dotenv` file to `.env` and edit configuration if necessary, then start docker

```sh
SERVICE=<service-name> make dev
# or rebuild 
SERVICE=<service-name> make dev-build
# sometimes we need to clean all database and run from scratch 
make dev-clean
# view logs 
docker-compose logs -f portal
```

At the first startup, it's easier to work with basic mock data. Just run the following script after `controller` are online.

```sh
make bootstrap
```

## Deployment 

### Test environment

Run docker-compose

### Production environment
