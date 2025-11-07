package opt

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

var (
	ErrNoOptDataFound    = errors.New("no opt data found")
	ErrEmptyOptList      = errors.New("opt list is empty")
	ErrInvalidType       = errors.New("invalid type")
	ErrOnlyOneOptRemains = errors.New("优选 IP 存量不足, 无法更换IP, 考虑重新进行优选操作")
)

// Repository 优选数据仓库
type Repository struct {
	mu       sync.RWMutex
	filePath string              // JSON 文件路径
	store    map[string]*OptData // key 为 type
}

// NewRepository 创建新的优选数据仓库
func NewRepository(filePath string) (*Repository, error) {
	repo := &Repository{
		filePath: filePath,
		store:    make(map[string]*OptData),
	}

	// 尝试从文件加载数据
	if err := repo.load(); err != nil {
		// 如果文件不存在,这是正常的,创建空存储
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load opt data: %w", err)
		}
	}

	return repo, nil
}

// load 从 JSON 文件加载数据
func (r *Repository) load() error {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return err
	}

	var optStore OptStore
	if err := json.Unmarshal(data, &optStore); err != nil {
		return fmt.Errorf("failed to unmarshal opt data: %w", err)
	}

	r.store = optStore.Opts
	if r.store == nil {
		r.store = make(map[string]*OptData)
	}

	return nil
}

// save 保存数据到 JSON 文件
func (r *Repository) save() error {
	optStore := OptStore{
		Opts: r.store,
	}

	data, err := json.MarshalIndent(optStore, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal opt data: %w", err)
	}

	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write opt data: %w", err)
	}

	return nil
}

// SaveOptData 保存优选数据(type 相同则替换,不同则新增)
func (r *Repository) SaveOptData(optType string, data []OptInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if optType == "" {
		return ErrInvalidType
	}

	r.store[optType] = &OptData{
		Type:    optType,
		Data:    data,
		Current: 0,
	}

	return r.save()
}

// GetCurrentOpt 获取指定类型的当前优选
func (r *Repository) GetCurrentOpt(optType string) (string, OptInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if optType == "" {
		return "", OptInfo{}, ErrInvalidType
	}

	optData, exists := r.store[optType]
	if !exists || len(optData.Data) == 0 {
		return "", OptInfo{}, ErrNoOptDataFound
	}

	if optData.Current >= len(optData.Data) {
		return "", OptInfo{}, ErrNoOptDataFound
	}

	return optType, optData.Data[optData.Current], nil
}

// ChangeToNext 切换到下一个优选,并删除当前的
func (r *Repository) ChangeToNext(optType string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if optType == "" {
		return ErrInvalidType
	}

	optData, exists := r.store[optType]
	if !exists || len(optData.Data) == 0 {
		return ErrEmptyOptList
	}

	// 删除当前的
	if optData.Current < len(optData.Data) {
		optData.Data = append(optData.Data[:optData.Current], optData.Data[optData.Current+1:]...)
	}

	// 如果删除后列表为空
	if len(optData.Data) == 0 {
		delete(r.store, optType)
		return r.save()
	}

	// 如果当前索引超出范围,重置为 0
	if optData.Current >= len(optData.Data) {
		optData.Current = 0
	}

	return r.save()
}

// GetAllTypes 获取所有优选类型
func (r *Repository) GetAllTypes() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]string, 0, len(r.store))
	for t := range r.store {
		types = append(types, t)
	}
	return types
}

// GetOptListSize 获取指定类型的优选列表大小
func (r *Repository) GetOptListSize(optType string) int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if optData, exists := r.store[optType]; exists {
		return len(optData.Data)
	}
	return 0
}
