package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"

	"gorm.io/gorm"
)

// JsonCodeResult 返回结构
type JsonCodeResult struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

// IndexCoderHandler 计数器接口
func IndexCoderHandler(w http.ResponseWriter, r *http.Request) {
	data, err := getCoderIndex()
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	fmt.Fprint(w, data)
}

// CodeHandler 计数器接口
func CodeHandler(w http.ResponseWriter, r *http.Request) {
	res := &JsonCodeResult{}

	if r.Method == http.MethodGet {
		code, err := getCurrentCode()
		if err != nil {
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else {
			res.Data = code.Code
		}
	} else if r.Method == http.MethodPost {
		code, err := modifyCode(r)
		if err != nil {
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else {
			res.Data = code
		}
	} else {
		res.Code = -1
		res.ErrorMsg = fmt.Sprintf("请求方法 %s 不支持", r.Method)
	}

	msg, err := json.Marshal(res)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}

// modifyCode 更新计数，自增或者清零
func modifyCode(r *http.Request) (int32, error) {
	action, err := getCoderAction(r)
	if err != nil {
		return 0, err
	}

	var count int32
	if action == "make" {
		count, err = UpsertCoder(r)
		if err != nil {
			return 0, err
		}
	} else if action == "clear" {
		err = ClearCoder()
		if err != nil {
			return 0, err
		}
		count = 0
	} else {
		err = fmt.Errorf("参数 action : %s 错误", action)
	}

	return count, err
}

// UpsertCoder 更新或修改计数器
func UpsertCoder(r *http.Request) (int32, error) {
	currentCode, err := getCurrentCode()
	var code int32
	createdAt := time.Now()
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	} else if err == gorm.ErrRecordNotFound {
		createdAt = time.Now()
	} else {
		createdAt = currentCode.CreatedAt
	}
	// code 赋值一个100000 到 999999 之间的随机数
	rand.Seed(time.Now().UnixNano())
	if currentCode.Code == 0 {
		code = rand.Int31n(900000) + 100000
	}

	coder := &model.CoderModel{
		Id:        1,
		Code:      code,
		CreatedAt: createdAt,
		UpdatedAt: time.Now(),
	}
	err = dao.CoderImp.UpsertCoder(coder)
	if err != nil {
		return 0, err
	}
	return coder.Code, nil
}

func ClearCoder() error {
	return dao.CoderImp.ClearCoder(1)
}

// getCurrentCode 查询当前计数器
func getCurrentCode() (*model.CoderModel, error) {
	Code, err := dao.CoderImp.GetCoder(1)
	if err != nil {
		return nil, err
	}

	return Code, nil
}

// getAction 获取action
func getCoderAction(r *http.Request) (string, error) {
	decoder := json.NewDecoder(r.Body)
	body := make(map[string]interface{})
	if err := decoder.Decode(&body); err != nil {
		return "", err
	}
	defer r.Body.Close()

	action, ok := body["action"]
	if !ok {
		return "", fmt.Errorf("缺少 action 参数")
	}

	return action.(string), nil
}

// getIndex 获取主页
func getCoderIndex() (string, error) {
	b, err := ioutil.ReadFile("./code.html")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
