# Host Manager Demo

Simple Go HTTP service implementing the host CRUD endpoints described in `openapi.json`. The service uses the local file `hosts.json` as its data store so it never touches the real operating system hosts file.

## 配置说明

程序使用 YAML 配置文件进行配置管理。首次运行时，如果配置文件不存在，会自动创建默认的 `config.yaml` 文件。

### 配置文件示例

```yaml
# 服务器配置
server:
  port: ":15920"

# 数据存储配置
data:
  host_file: "hosts.json"

# CORS 跨域配置
cors:
  allow_origins:
    - "*"
  allow_methods:
    - "GET"
    - "POST"
    - "DELETE"
    - "OPTIONS"
  allow_headers:
    - "Content-Type"
    - "Authorization"
    - "X-Requested-With"
  expose_headers:
    - "Content-Length"
  allow_credentials: false
  max_age: "12h"
```

## 运行

### 使用默认配置文件（config.yaml）

```bash
go run .
```

### 指定配置文件

```bash
go run . --config /path/to/custom-config.yaml
```

或者

```bash
go run . --config=myconfig.yaml
```

服务器默认监听在 `http://localhost:15920`（可通过配置文件修改）。

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
