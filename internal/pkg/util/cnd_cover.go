package util

import (
	"github.com/mlogclub/simple"
)

type SqlCnd2 struct {
	simple.SqlCnd
}

func NewSqlCnd2() *SqlCnd2 {
	return &SqlCnd2{}
}

func (q *SqlCnd2) EqByReq(column string, value string) *SqlCnd2 {
	if len(value) > 0 {
		q.Eq(column, value)
	}
	return q
}
