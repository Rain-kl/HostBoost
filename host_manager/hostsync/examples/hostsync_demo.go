package main

import (
	"flag"
	"fmt"
	"hostMgr/hostsync"
	"log"
	"os"
)

func main() {
	// 命令行参数
	syncCmd := flag.Bool("sync", false, "同步 hosts.json 到系统 hosts 文件")
	listBackups := flag.Bool("list-backups", false, "列出所有备份文件")
	restore := flag.String("restore", "", "从指定备份恢复 (文件名)")
	restoreLatest := flag.Bool("restore-latest", false, "从最新备份恢复")
	deleteBackups := flag.Bool("delete-backups", false, "删除所有备份文件")
	noBackup := flag.Bool("no-backup", false, "同步时不创建备份")
	hostsFile := flag.String("hosts", "hosts.json", "指定 hosts.json 文件路径")

	flag.Parse()

	// 创建 Syncer 实例
	syncer := hostsync.NewSyncer(*hostsFile)

	// 如果指定不备份
	if *noBackup {
		syncer.SetBackupEnabled(false)
	}

	// 根据命令执行不同操作
	switch {
	case *syncCmd:
		syncHosts(syncer)
	case *listBackups:
		listBackupFiles(syncer)
	case *restore != "":
		restoreFromBackup(syncer, *restore)
	case *restoreLatest:
		restoreFromLatest(syncer)
	case *deleteBackups:
		deleteAllBackups(syncer)
	default:
		flag.Usage()
		fmt.Println("\n示例:")
		fmt.Println("  sudo go run examples/hostsync_demo.go -sync")
		fmt.Println("  sudo go run examples/hostsync_demo.go -list-backups")
		fmt.Println("  sudo go run examples/hostsync_demo.go -restore hosts_backup_20250102_150405")
		fmt.Println("  sudo go run examples/hostsync_demo.go -restore-latest")
		os.Exit(1)
	}
}

func syncHosts(syncer *hostsync.Syncer) {
	fmt.Println("🔍 检查权限...")
	if err := syncer.ValidatePermissions(); err != nil {
		if err == hostsync.ErrPermissionDenied {
			log.Fatalf("❌ 权限不足: %v\n💡 提示: 请使用 sudo 运行程序", err)
		}
		log.Fatalf("❌ 权限检查失败: %v", err)
	}

	fmt.Printf("📖 读取 %s...\n", syncer.GetHostsJSONPath())
	entries, err := syncer.SyncFromJSON()
	if err != nil {
		log.Fatalf("❌ 读取失败: %v", err)
	}

	fmt.Printf("📝 找到 %d 条 host 记录\n", len(entries))
	for _, entry := range entries {
		fmt.Printf("  - %s -> %s (类型: %s)\n", entry.Domain, entry.IP, entry.Type)
	}

	fmt.Printf("\n🔄 同步到系统 hosts 文件: %s\n", syncer.GetSystemHostsPath())
	if err := syncer.Sync(); err != nil {
		log.Fatalf("❌ 同步失败: %v", err)
	}

	fmt.Println("✅ Hosts 同步成功!")
	fmt.Println("🔄 DNS 缓存已自动刷新")
}

func listBackupFiles(syncer *hostsync.Syncer) {
	fmt.Println("📋 备份文件列表:")
	backups, err := syncer.ListBackups()
	if err != nil {
		log.Fatalf("❌ 获取备份列表失败: %v", err)
	}

	if len(backups) == 0 {
		fmt.Println("  (无备份文件)")
		return
	}

	for i, backup := range backups {
		fmt.Printf("  %d. %s\n", i+1, backup)
	}
}

func restoreFromBackup(syncer *hostsync.Syncer, backupName string) {
	fmt.Printf("🔄 从备份恢复: %s\n", backupName)
	if err := syncer.RestoreFromBackup(backupName); err != nil {
		if err == hostsync.ErrPermissionDenied {
			log.Fatalf("❌ 权限不足: %v\n💡 提示: 请使用 sudo 运行程序", err)
		}
		log.Fatalf("❌ 恢复失败: %v", err)
	}
	fmt.Println("✅ 恢复成功!")
}

func restoreFromLatest(syncer *hostsync.Syncer) {
	fmt.Println("🔄 从最新备份恢复...")
	if err := syncer.RestoreLatestBackup(); err != nil {
		if err == hostsync.ErrPermissionDenied {
			log.Fatalf("❌ 权限不足: %v\n💡 提示: 请使用 sudo 运行程序", err)
		}
		log.Fatalf("❌ 恢复失败: %v", err)
	}
	fmt.Println("✅ 恢复成功!")
}

func deleteAllBackups(syncer *hostsync.Syncer) {
	fmt.Println("🗑️  删除所有备份文件...")
	if err := syncer.DeleteAllBackups(); err != nil {
		log.Fatalf("❌ 删除失败: %v", err)
	}
	fmt.Println("✅ 删除成功!")
}
