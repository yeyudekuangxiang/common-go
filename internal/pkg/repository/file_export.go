package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultFileExportRepository = FileExportRepository{DB: app.DB}

type FileExportRepository struct {
	DB *gorm.DB
}

func (repo FileExportRepository) Create(export *entity.FileExport) error {
	return repo.DB.Create(export).Error
}

func (repo FileExportRepository) Save(export *entity.FileExport) error {
	return repo.DB.Save(export).Error
}
func (repo FileExportRepository) GetPageList(by GetFileExportPageListBy) ([]entity.FileExport, int64) {
	db := repo.DB.Model(entity.FileExport{})

	list := make([]entity.FileExport, 0)
	var total int64

	if by.Type != nil {
		db.Where("type = ?", by.Type)
	}
	if by.AdminId > 0 {
		db.Where("admin_id = ?", by.AdminId)
	}
	if !by.StartCreatedAt.IsZero() {
		db.Where("created_at >=", by.StartCreatedAt)
	}
	if !by.EndCreatedAt.IsZero() {
		db.Where("created_at <=", by.EndCreatedAt)
	}
	if by.Status > 0 {
		db.Where("status = ?", by.Status)
	}

	for _, orderBy := range by.OrderBy {
		switch orderBy {
		case entity.OrderByFileExportCreatedAtDesc:
			db.Order("created_at desc")
		case entity.OrderByFileExportUpdatedAtDesc:
			db.Order("updated_at desc")
		}
	}

	err := db.Count(&total).Offset(by.Offset).Limit(by.Limit).Find(&list).Error
	if err != nil {
		panic(err)
	}
	return list, total
}
func (repo FileExportRepository) FindById(id int64) entity.FileExport {
	fileExport := entity.FileExport{}
	if err := repo.DB.First(&fileExport, id).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}
	return fileExport
}
