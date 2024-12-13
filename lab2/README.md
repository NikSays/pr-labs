# Communicator

PR Lab 2

By Nejintev Nicolai, FAF-221

---

## Running locally

```shell
docker compose up --build
```

## Configuration

Configure environment variables in docker-compose.yml

## Usage

* HTTP server is available at `:8080`
* WS server is available at `:8081`
* TCP server is available at `:8082`

### Endpoints

* `POST :8080/movie (name, description, rating)` – Create a movie
* `GET :8080/movie (page, page_size)` – Returns paginated movies
* `PUT :8080/movie/:id (name, description, rating)` – Update a movie
* `DELETE :8080/movie/:id` – Delete a movie
* `GET :8081` – WebSocket chat
  * `{"username: "ABC","message":"ABC"}` – Send message
  * `{"username: "ABC","message":"/leave"}` – Disconnect
* `:8082` – Raw TCP server
  * `w ABC` – Writes data to file
  * `r` – Reads data from file
