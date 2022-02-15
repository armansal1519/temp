package payment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const apiKey = "1cbb528b-9866-4952-912b-9e51e399b851"

type nextPayResp struct {
	Code    int    `json:"code"`
	TransId string `json:"trans_id"`
	Amount  string `json:"amount"`
}

func GetPaymentUrl(amount string, orderId string, callBackUrl string) (string, error) {
	url := "https://nextpay.org/nx/gateway/token"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf("api_key=%v&amount=%v&order_id=%v&callback_uri=%v", apiKey, amount, orderId, callBackUrl))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {

		return "", fmt.Errorf("1-%v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("2-%v", err)

	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("3-%v", err)

	}
	var p nextPayResp

	json.Unmarshal(body, &p)
	if p.Code != -1 {
		return "", fmt.Errorf("nextpay error number:%v", p.Code)
	}
	//return fmt.Sprintf("https://nextpay.org/nx/gateway/payment/%v", p.TransId), nil
	return p.TransId, nil
}

type VerifyResp struct {
	Code          int    `json:"code"`
	Amount        int    `json:"amount"`
	OrderId       string `json:"order_id"`
	CardHolder    string `json:"card_holder"`
	ShaparakRefId string `json:"Shaparak_Ref_Id"`
}

func Verify(amount int64, transId string) (VerifyResp, error) {
	url := "https://nextpay.org/nx/gateway/verify"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf("api_key=%v&amount=%v&trans_id=%v", apiKey, amount, transId))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return VerifyResp{}, err

	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return VerifyResp{}, err

	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {

		return VerifyResp{}, err
	}
	fmt.Println(string(body))
	var v VerifyResp
	json.Unmarshal(body, &v)
	if v.Code != 0 {
		return v, fmt.Errorf("problem in varify %v", v.Code)
	}
	return v, nil

}
