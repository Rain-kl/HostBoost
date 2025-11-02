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

## Host Sync (系统 Hosts 文件同步)

`hostsync` 包提供了将 `hosts.json` 文件同步到系统 hosts 文件的功能。

### 功能特性

- 📁 读取 `hosts.json` 文件并解析 host 条目
- 🔄 **覆盖模式同步**：删除管理区域所有旧内容，完全使用 `hosts.json` 的内容重新写入
- 💾 自动备份系统 hosts 文件（保留最近 10 个备份）
- 🔒 使用标记注释管理同步区域，不影响系统 hosts 文件的其他内容
- 🔙 支持从备份恢复系统 hosts 文件
- ⚠️ 权限检查，确保有足够的权限修改系统文件
- 🌍 跨平台支持（Windows、macOS、Linux）
- 🚀 **自动刷新 DNS 缓存**：同步后自动刷新系统 DNS 缓存，使更改立即生效

### 使用示例

```go
package main

import (
	"fmt"
	"log"
	"hostMgr/hostsync"
)

func main() {
	// 创建 Syncer 实例
	syncer := hostsync.NewSyncer("hosts.json")
	
	// 检查权限（可选）
	if err := syncer.ValidatePermissions(); err != nil {
		log.Fatalf("权限不足: %v\n提示: 请使用 sudo 运行程序", err)
	}
	
	// 同步 hosts.json 到系统 hosts 文件
	if err := syncer.Sync(); err != nil {
		log.Fatalf("同步失败: %v", err)
	}
	
	fmt.Println("✅ Hosts 同步成功!")
}
```

### DNS 缓存刷新

同步成功后会自动刷新系统 DNS 缓存，使更改立即生效。支持的平台：

- **macOS**: 使用 `dscacheutil -flushcache` 和 `killall -HUP mDNSResponder`
- **Linux**: 自动尝试 `systemd-resolve`、`resolvectl`、`nscd` 等
- **Windows**: 使用 `ipconfig /flushdns`

如需禁用自动刷新：

```go
syncer := hostsync.NewSyncer("hosts.json")
syncer.SetAutoFlushDNSCache(false) // 禁用自动 DNS 缓存刷新
syncer.Sync()
```

手动刷新 DNS 缓存：

```go
syncer := hostsync.NewSyncer("hosts.json")
if err := syncer.FlushDNSCache(); err != nil {
	log.Printf("刷新 DNS 缓存失败: %v", err)
}
```

### 禁用自动备份

```go
syncer := hostsync.NewSyncer("hosts.json")
syncer.SetBackupEnabled(false) // 禁用自动备份
syncer.Sync()
```

### 备份管理

```go
syncer := hostsync.NewSyncer("hosts.json")

// 列出所有备份
backups, err := syncer.ListBackups()
if err != nil {
	log.Fatal(err)
}
for _, backup := range backups {
	fmt.Println("备份:", backup)
}

// 从指定备份恢复
err = syncer.RestoreFromBackup("hosts_backup_20250102_150405")

// 从最新备份恢复
err = syncer.RestoreLatestBackup()

// 删除所有备份
err = syncer.DeleteAllBackups()
```

### 系统 Hosts 文件位置

- **Windows**: `C:\Windows\System32\drivers\etc\hosts`
- **macOS/Linux**: `/etc/hosts`

### 同步模式说明

**覆盖模式**：每次同步时，会完全删除管理区域内的所有旧内容，然后使用 `hosts.json` 中的内容重新写入。这确保了管理区域与 `hosts.json` 文件保持完全一致。

### 管理区域标记

同步时会在系统 hosts 文件中使用标记注释来标识管理区域：

```
# === HostBoost Managed Section Start ===
172.66.166.61   github.com      # type:cloudflare
# === HostBoost Managed Section End ===
```

- **管理区域内**：每次同步时会被完全覆盖为 `hosts.json` 的内容
- **管理区域外**：不会被修改，确保与其他工具或手动配置兼容

### 权限要求

修改系统 hosts 文件需要管理员权限：

- **macOS/Linux**: 使用 `sudo` 运行程序
  ```bash
  sudo go run .
  ```

- **Windows**: 以管理员身份运行命令提示符或 PowerShell

### 注意事项

⚠️ **重要**: 
- 修改系统 hosts 文件可能影响网络连接和域名解析
- 建议在修改前先备份（默认已启用自动备份）
- 确保 `hosts.json` 文件格式正确
- 需要管理员权限才能修改系统 hosts 文件
