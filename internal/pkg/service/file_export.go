package service

import (
	"github.com/pkg/errors"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
)

var DefaultFileExportService = FileExportService{repo: repository.DefaultFileExportRepository}

type FileExportService struct {
	repo repository.FileExportRepository
}

func (srv FileExportService) Add(param AddFileExportParam) (*entity.FileExport, error) {
	fileExport := entity.FileExport{
		AdminId:   param.AdminId,
		Type:      param.Type,
		Params:    param.Params,
		CreatedAt: model.NewTime(),
		UpdatedAt: model.NewTime(),
		Status:    entity.FileExportStatusWaiting,
	}
	return &fileExport, srv.repo.Create(&fileExport)
}

func (srv FileExportService) Update(id int64, param UpdateFileExportParam) (*entity.FileExport, error) {
	fileExport := srv.repo.FindById(id)
	if fileExport.ID == 0 {
		return nil, errors.New("导出任务不存在")
	}
	fileExport.Url = param.Url
	fileExport.Status = param.Status
	fileExport.Message = param.Message
	return &fileExport, srv.repo.Save(&fileExport)
}
func (srv FileExportService) GetPageList(by repository.GetFileExportPageListBy) ([]FileExportRecord, int64, error) {
	list, total := srv.repo.GetPageList(by)

	recordList := make([]FileExportRecord, 0)
	for _, fileExport := range list {
		admin, err := DefaultSystemAdminService.GetAdminById(fileExport.AdminId)
		if err != nil {
			return nil, 0, err
		}
		recordList = append(recordList, FileExportRecord{
			FileExport: fileExport,
			TypeText:   fileExport.Type.Text(),
			StatusText: fileExport.Status.Text(),
			Admin:      *admin,
		})
	}
	return recordList, total, nil
}
func (srv FileExportService) GetFileExportStatusAndTypeList() FileExportStatusAndType {
	statusList := make([]FileExportStatus, 0)
	for _, status := range entity.FileExportStatusList {
		statusList = append(statusList, FileExportStatus{
			Status:     status,
			StatusText: status.Text(),
		})
	}
	typeList := make([]FileExportType, 0)

	for _, t := range entity.FileExportTypeList {
		typeList = append(typeList, FileExportType{
			Type:     t,
			TypeText: t.Text(),
		})
	}
	return FileExportStatusAndType{
		StatusList: statusList,
		TypeList:   typeList,
	}
}
func (srv FileExportService) FindById(id int64) (*entity.FileExport, error) {
	fe := srv.repo.FindById(id)
	return &fe, nil
}
