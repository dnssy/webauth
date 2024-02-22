package dao

import (
	"wxcloudrun-golang/db/model"
)

// CounterInterface 计数器数据模型接口
type CounterInterface interface {
	GetCounter(id int32) (*model.CounterModel, error)
	UpsertCounter(counter *model.CounterModel) error
	ClearCounter(id int32) error
}

// CounterInterfaceImp 计数器数据模型实现
type CounterInterfaceImp struct{}

// Imp 实现实例
var Imp CounterInterface = &CounterInterfaceImp{}

// CoderInterface 计数器数据模型接口
type CoderInterface interface {
	GetCoder(id int32) (*model.CoderModel, error)
	UpsertCoder(counter *model.CoderModel) error
	ClearCoder(id int32) error
}

// CoderInterfaceImp 计数器数据模型实现
type CoderInterfaceImp struct{}

// Imp CoderInterface 实现实例
var Imp CoderInterface = &CoderInterfaceImp{}
