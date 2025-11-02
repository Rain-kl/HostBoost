package opt

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

// ReportOpt 上报优选数据
func (s *Service) ReportOpt(req ReportRequest) error {
	if len(req.Data) == 0 {
		return ErrEmptyOptList
	}

	return s.repo.SaveOptData(req.Type, req.Data)
}

// GetCurrentOpt 获取指定类型的当前优选
func (s *Service) GetCurrentOpt(optType string) (string, OptInfo, error) {
	return s.repo.GetCurrentOpt(optType)
}

// ChangeOpt 更换指定类型的当前优选
func (s *Service) ChangeOpt(optType string) error {
	return s.repo.ChangeToNext(optType)
}

// GetAllTypes 获取所有优选类型
func (s *Service) GetAllTypes() []string {
	return s.repo.GetAllTypes()
}
