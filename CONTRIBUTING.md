# Contributing

## Backend

### Project structure

```
services
├── auth
├── controller
├── controller-api
├── ...
```

### Define new module

Each module is defined independently in separated database. So basically there isn't restriction between modules. We can copy the module's database and metadata to another backend.

First most, define new database connection environment to `controller`:
- Update `docker-compose.yaml`
- Use Hasura CLI to add data source and migrations 

### Squash migrations

Ideally we should create only one migration folder per PR. Let's squash migrations before merging

### Recommended development tools

- Visual Studio Code 
- Extensions: ESlint, Prettier

## Local development

View `README.md`

## Development environment 

### Managed services

- Hasura cloud 
- Cloud SQL