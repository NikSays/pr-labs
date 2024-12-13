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

* `POST :8080/monitor (name, price_mdl, price_eur, warranty)` – Create a monitor
* `GET :8080/monitor (page, page_size)` – Returns paginated monitor
* `PUT :8080/monitor/:id (name, price_mdl, price_eur, warranty)` – Update a monitor
* `DELETE :8080/monitor/:id` – Delete a monitor
