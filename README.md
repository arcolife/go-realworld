- [Getting started](#getting-started)
  - [Pre-requisites](#pre-requisites)
  - [Locally](#locally)
  - [Testing](#testing)
    - [Direct](#direct)
    - [Docker](#docker)
- [TODO](#todo)

---

![logo](logo.png)

- [Go net/http](//golang.org) codebase containing real world examples (CRUD, auth, advanced patterns, etc) that adheres to the [RealWorld](https://github.com/gothinkster/realworld) spec and API.

- [Demo](https://demo.realworld.io/)
- [RealWorld framework](https://github.com/gothinkster/realworld)

---

This codebase was created to demonstrate a fully fledged fullstack application built with **Go net/http library** including CRUD operations, authentication, routing, pagination, and more.

We've gone to great lengths to adhere to the **Golang** community styleguides & best practices.

For more information on how to this works with other frontends/backends, head over to the [RealWorld](https://github.com/gothinkster/realworld) repo.


## How it works

The project structure was inspired by two posts on [Ben Johnson's](https://twitter.com/benbjohnson) blog which can be found [here](https://www.gobeyond.dev/packages-as-layers/) and [here](https://www.gobeyond.dev/standard-package-layout/).

# Getting started

This project uses Go version 1.17 and postgresql 14

## Pre-requisites

You also need to have [migrate](https://github.com/golang-migrate/migrate) tool installed to run all migrations against the database

refs:
- https://github.com/golang-migrate/migrate/blob/master/GETTING_STARTED.md
- https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
- https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md
- https://github.com/golang-migrate/migrate/issues/179#issuecomment-475821264

```sh
source .env 
make run-migration
make down-migration
make down-migration
# migrate -database postgres://admin:Y3Z5wiuefz4tidZW@localhost:5432/conduit?sslmode=disable -path postgres/migrations down 
# Are you sure you want to apply all down migrations? [y/N]
# y
# Applying all down migrations
# error: Dirty database version 3. Fix and force version.
# make: *** [Makefile:13: down-migration] Error 1

make force-migration version=1
# migrate -database postgres://admin:Y3Z5wiuefz4tidZW@localhost:5432/conduit?sslmode=disable -path postgres/migrations force 1

make run-migration           
# migrate -database postgres://admin:Y3Z5wiuefz4tidZW@localhost:5432/conduit?sslmode=disable -path postgres/migrations up
# 2/u create_users_table (12.569443ms)
# 3/u create_articles_table (25.719797ms)
# 4/u create_tags_table (35.836959ms)
# 5/u create_follow_table (53.77812ms)
# 6/u add_comments_table (62.155861ms)
# 7/u create_comments_table (80.76418ms)
# 8/u add_favorite_table (94.583558ms)
```

## Locally

- make sure [Go](https://golang.org/dl) is installed on your machine.
- make sure to have the postgresql database installed locally or remote. Refer to [docker/README.md](./docker/README.md)
- set the .env file or env var as shown in the .env.example file
- Install migrate tool
- run the migrations in postgres/migrations or run `make run-migrate`
- fetch all dependencies using `go mod download`
- run `make run` to start the server locally

- http://localhost:9000/api/v1/health

  ```json
  {"data":{"hello":"beautiful"},"message":"healthy","status":"available"}
  ```

## Testing


```sh
# start DB. Refer to docker/
docker run --rm --name postgresql --env-file docker/db.env -p 5432:5432 bitnami/postgresql:14

# follow all steps above; then start server
cd ..
make run
```

### Direct

- **realworld's newman**:

```sh
chmod +x ./api/run-api-tests.sh
export APIURL=127.0.0.1:9000/api/v1

# run it first time
./api/run-api-tests.sh
# if it fails, it'll prompt to install newman (postman's cli tool)
# then run it again
./api/run-api-tests.sh
```

- **custom tests**:

```sh
curl -XPOST -H "Content-Type: application/json" 127.0.0.1:9000/api/v1/users -d @data.json
```

### Docker

```sh
docker run --network="host" postman/newman_alpine33 \
  run https://github.com/gothinkster/realworld/raw/main/api/Conduit.postman_collection.json  \
      --global-var "APIURL=http://127.0.0.1:9000/api/v1" \
      --global-var "USERNAME=user2021" \
      --global-var "EMAIL=user2021@example.com" \
      --global-var "PASSWORD=password"
```

# TODO
- Revisit error handling
