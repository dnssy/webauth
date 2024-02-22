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

// JsonResult 返回结构
type JsonResult struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

// IndexHandler 计数器接口
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	data, err := getIndex()
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	fmt.Fprint(w, data)
}

// CodeHandler 计数器接口
func CodeHandler(w http.ResponseWriter, r *http.Request) {
	res := &JsonResult{}

	if r.Method == http.MethodGet {
		code, err := getCurrentCode()
		if err != nil {
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else {
			res.Data = code.Count
		}
	} else if r.Method == http.MethodPost {
		count, err := modifyCode(r)
		if err != nil {
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else {
			res.Data = count
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
	action, err := getAction(r)
	if err != nil {
		return 0, err
	}

	var count int32
	if action == "inc" {
		count, err = upsertCode(r)
		if err != nil {
			return 0, err
		}
	} else if action == "clear" {
		err = clearCode()
		if err != nil {
			return 0, err
		}
		count = 0
	} else {
		err = fmt.Errorf("参数 action : %s 错误", action)
	}

	return count, err
}

// upsertCode 更新或修改计数器
func upsertCode(r *http.Request) (int32, error) {
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
	if currentCode.Count == 0 {
		code = rand.Int31n(900000) + 100000
	}

	coder := &model.CodeModel{
		Id:        1,
		Code:      code,
		CreatedAt: createdAt,
		UpdatedAt: time.Now(),
	}
	err = dao.Imp.UpsertCode(coder)
	if err != nil {
		return 0, err
	}
	return coder.Count, nil
}

func clearCode() error {
	return dao.Imp.ClearCode(1)
}

// getCurrentCode 查询当前计数器
func getCurrentCode() (*model.CodeModel, error) {
	Code, err := dao.Imp.GetCode(1)
	if err != nil {
		return nil, err
	}

	return Code, nil
}

// getAction 获取action
func getAction(r *http.Request) (string, error) {
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
func getIndex() (string, error) {
	b, err := ioutil.ReadFile("./code.html")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
