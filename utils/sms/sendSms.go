package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type fastSms struct {
	ParameterArray []ParameterArray `json:"ParameterArray"`
	Mobile         string           `json:"Mobile"`
	TemplateId     string           `json:"TemplateId"`
	UserApiKey     string           `json:"UserApiKey"`
	SecretKey      string           `json:"SecretKey"`
}

type ParameterArray struct {
	Parameter      string `json:"Parameter"`
	ParameterValue string `json:"ParameterValue"`
}

const (
	userApiKey = "dfa1028e6c684ce6ce9d0f3f"
	secretKey  = "saeedKard@shiyan"
)

func SendSms(phoneNumber string, templateId string, arr []ParameterArray) error {
	f := fastSms{
		ParameterArray: arr,
		Mobile:         phoneNumber,
		TemplateId:     templateId,
		UserApiKey:     userApiKey,
		SecretKey:      secretKey,
	}
	postBody, _ := json.Marshal(f)
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("https://RestfulSms.com/api/UltraFastSend/UserApiKey", "application/json", responseBody)
	//Handle Error
	if err != nil {
		return fmt.Errorf("an Error Occured %v", err)
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println(string(body))
		return fmt.Errorf("an Error Occured %v", err)

	}
	println(string(body))
	return nil
}
