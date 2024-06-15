package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	yapiHost    = "http://yapi.wondershare.cn"
	stringType  = "string"
	numberType  = "number"
	arrayType   = "array"
	objectType  = "object"
	booleanType = "boolean"
	integerType = "integer"
)

type Client struct {
}

type InterfaceInfoResp struct {
	Errcode int            `json:"errcode"`
	Errmsg  string         `json:"errmsg"`
	Data    *InterfaceInfo `json:"data"`
}

type InterfaceInfo struct {
	QueryPath struct {
		Path   string        `json:"path"`
		Params []interface{} `json:"params"`
	} `json:"query_path"`
	EditUID             int           `json:"edit_uid"`
	Status              string        `json:"status"`
	Type                string        `json:"type"`
	ReqBodyIsJSONSchema bool          `json:"req_body_is_json_schema"`
	ResBodyIsJSONSchema bool          `json:"res_body_is_json_schema"`
	APIOpened           bool          `json:"api_opened"`
	Index               int           `json:"index"`
	Tag                 []interface{} `json:"tag"`
	ID                  int           `json:"_id"`
	Method              string        `json:"method"`
	Catid               int           `json:"catid"`
	Title               string        `json:"title"`
	Path                string        `json:"path"`
	ProjectID           int           `json:"project_id"`
	ReqParams           []interface{} `json:"req_params"`
	ResBodyType         string        `json:"res_body_type"`
	UID                 int           `json:"uid"`
	AddTime             int           `json:"add_time"`
	UpTime              int           `json:"up_time"`
	ReqQuery            []interface{} `json:"req_query"`
	ReqHeaders          []struct {
		Required string `json:"required"`
		ID       string `json:"_id"`
		Name     string `json:"name"`
		Value    string `json:"value"`
	} `json:"req_headers"`
	ReqBodyForm  []interface{} `json:"req_body_form"`
	V            int           `json:"__v"`
	Desc         string        `json:"desc"`
	Markdown     string        `json:"markdown"`
	ReqBodyOther string        `json:"req_body_other"`
	ReqBodyType  string        `json:"req_body_type"`
	ResBody      string        `json:"res_body"`
	Username     string        `json:"username"`
}

func (c *Client) QueryInterfaceInfo(token string, interfaceId int) (*InterfaceInfo, error) {
	url := yapiHost + fmt.Sprintf("/api/interface/get?id=%d&token=%s", interfaceId, token)
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println("QueryInterfaceInfo http.NewRequest err:", err)
		return nil, yapiReqError
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println("QueryInterfaceInfo client.Do err:", err)
		return nil, yapiReqError
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("QueryInterfaceInfo ioutil.ReadAll err:", err)
		return nil, yapiReqError
	}
	resp := InterfaceInfoResp{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Println("QueryInterfaceInfo json.Unmarshal err:", err)
		return nil, yapiReqError
	}
	if resp.Errcode != 0 {
		log.Printf("QueryInterfaceInfo req not success, resp:%s", string(body))
		return nil, yapiReqError
	}
	return resp.Data, nil
}
