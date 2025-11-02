package opt

import (
	"hostMgr/internal/extSvc"
	"log"
)

// Service 优选服务
type Service struct {
	repo *Repository
}

// NewService 创建新的优选服务
func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// HostService 定义 host service 接口，用于避免循环依赖
type HostService interface {
	UpdateHostsByType(hostType string) (int, error)
}

// ReportOpt 上报优选数据
func (s *Service) ReportOpt(req ReportRequest) error {
	if len(req.Data) == 0 {
		return ErrEmptyOptList
	}

	// 保存优选数据
	if err := s.repo.SaveOptData(req.Type, req.Data); err != nil {
		return err
	}

	// 调用 host service 的 UpdateHostsByType 方法更新相关主机的 IP
	if hostSvc, ok := extSvc.HostService.(HostService); ok && hostSvc != nil {
		count, err := hostSvc.UpdateHostsByType(req.Type)
		if err != nil {
			// 记录错误但不影响优选数据的保存
			log.Printf("Warning: failed to update hosts after reporting opt (type=%s): %v", req.Type, err)
		} else {
			log.Printf("Successfully updated %d host(s) with new optimal IP for type %s", count, req.Type)
		}
	} else {
		log.Printf("Warning: host service not available for updating hosts (type=%s)", req.Type)
	}

	return nil
}

// GetCurrentOpt 获取指定类型的当前优选
func (s *Service) GetCurrentOpt(optType string) (string, OptInfo, error) {
	return s.repo.GetCurrentOpt(optType)
}

// ChangeOpt 更换指定类型的当前优选
func (s *Service) ChangeOpt(optType string) error {
	// 检查列表数量，如果只剩一个 IP 则阻止更换
	listSize := s.repo.GetOptListSize(optType)
	if listSize == 0 {
		return ErrNoOptDataFound
	}
	if listSize <= 1 {
		return ErrOnlyOneOptRemains
	}

	// 更换到下一个优选
	if err := s.repo.ChangeToNext(optType); err != nil {
		return err
	}

	// 调用 host service 的 UpdateHostsByType 方法更新相关主机的 IP
	if hostSvc, ok := extSvc.HostService.(HostService); ok && hostSvc != nil {
		count, err := hostSvc.UpdateHostsByType(optType)
		if err != nil {
			// 记录错误但不影响优选的更换
			log.Printf("Warning: failed to update hosts after changing opt (type=%s): %v", optType, err)
		} else {
			log.Printf("Successfully updated %d host(s) with new optimal IP after changing opt for type %s", count, optType)
		}
	} else {
		log.Printf("Warning: host service not available for updating hosts after changing opt (type=%s)", optType)
	}

	return nil
}

// GetAllTypes 获取所有优选类型
func (s *Service) GetAllTypes() []string {
	return s.repo.GetAllTypes()
}
