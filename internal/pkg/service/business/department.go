package business

import (
	"errors"
	"mio/internal/pkg/model/entity/business"
	brepo "mio/internal/pkg/repository/business"
)

var DefaultDepartmentService = DepartmentService{repo: brepo.DefaultDepartmentRepository}

type DepartmentService struct {
	repo brepo.DepartmentRepository
}

// GetBusinessDepartmentById 根据用户id查询部门信息
func (u DepartmentService) GetBusinessDepartmentById(id int) (*business.Department, error) {
	Department := u.repo.GetDepartmentBy(brepo.GetDepartmentBy{ID: id})
	if Department.ID != 0 {
		return &Department, nil
	}
	return nil, errors.New("非企业用户,请联系管理员开通企业版权限")
}

// GetBusinessDepartmentByIds 批量查询部门信息
func (u DepartmentService) GetBusinessDepartmentByIds(ids []int) []business.Department {
	Departments := u.repo.GetDepartmentListBy(brepo.GetDepartmentListBy{Ids: ids})
	return Departments
}
