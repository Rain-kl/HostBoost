# Host Manager Demo

Simple Go HTTP service implementing the host CRUD endpoints described in `openapi.json`. The service uses the local file `hosts_mock.json` as its data store so it never touches the real operating system hosts file.

## Run

```bash
go run .
```

The server listens on `http://localhost:8080` by default.

## Sample Requests  

- List all hosts:
  ```bash
  curl http://localhost:8080/host/list
  ```
- Fetch a single host:
  ```bash
  curl "http://localhost:8080/host?domain=example.local"
  ```
- Add a host:
  ```bash
  curl -X POST http://localhost:8080/host \
       -H "Content-Type: application/json" \
       -d '{"domain":"demo.local","ip":"10.1.1.2","type":"dev"}'
  ```
- Delete a host:
  ```bash
  curl -X DELETE http://localhost:8080/host \
       -H "Content-Type: application/json" \
       -d '{"domain":"demo.local"}'
  ```

Responses follow the shapes defined in the OpenAPI document.
