package dao

import (
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

const tableName = "Counters"

// ClearCounter 清除Counter
func (imp *CounterInterfaceImp) ClearCounter(id int32) error {
	cli := db.Get()
	return cli.Table(tableName).Delete(&model.CounterModel{Id: id}).Error
}

// UpsertCounter 更新/写入counter
func (imp *CounterInterfaceImp) UpsertCounter(counter *model.CounterModel) error {
	cli := db.Get()
	return cli.Table(tableName).Save(counter).Error
}

// GetCounter 查询Counter
func (imp *CounterInterfaceImp) GetCounter(id int32) (*model.CounterModel, error) {
	var err error
	var counter = new(model.CounterModel)

	cli := db.Get()
	err = cli.Table(tableName).Where("id = ?", id).First(counter).Error

	return counter, err
}

const tableName_Coders = "Coders"

// ClearCoder 清除Coder
func (imp *CoderInterfaceImp) ClearCoder(id int32) error {
	cli := db.Get()
	return cli.Table(tableName_Coders).Delete(&model.CoderModel{Id: id}).Error
}

// UpsertCoder 更新/写入coder
func (imp *CoderInterfaceImp) UpsertCoder(coder *model.CoderModel) error {
	cli := db.Get()
	return cli.Table(tableName_Coders).Save(coder).Error
}

// GetCoder 查询Coder
func (imp *CoderInterfaceImp) GetCoder(id int32) (*model.CoderModel, error) {
	var err error
	var coder = new(model.CoderModel)

	cli := db.Get()
	err = cli.Table(tableName_Coders).Where("id = ?", id).First(coder).Error

	return coder, err
}
