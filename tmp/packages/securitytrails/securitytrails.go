package securitytrails

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Find(domin string, key string) ([]string, error) {

	resp, err := call(domin, "a", key)
	if err != nil {
		return nil, err
	}
	ips, err := aResHandler(resp)
	if err != nil {
		return nil, err
	}
	resp, err = call("partdp.ir", "mx", key)
	if err != nil {
		return nil, err
	}
	mailHosts, err := mxResHandler(resp)
	if err != nil {
		return nil, err
	}
	for _, host := range mailHosts {
		if host != domin {
			//fmt.Println(host)
			resp, err := call(host, "a", key)
			if err != nil {
				fmt.Println(err)
			}
			ipList, err := aResHandler(resp)
			if err != nil {
				return nil, err
			}
			ips = append(ips, ipList...)
		}
	}
	ips = RemoveDuplicateStr(ips)
	return ips, nil
}

func call(domin string, recType string, key string) (*http.Response, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.securitytrails.com/v1/history/%s/dns/%s", domin, recType), nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Accept": {"application/json"},
		"APIKEY": {key},
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func aResHandler(resp *http.Response) ([]string, error) {
	var ips []string
	var results map[string]interface{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// close response body
	resp.Body.Close()
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}

	records, ok := results["records"].([]interface{})
	if !ok {
		return nil, errors.New("[!1] : no records found for this domin , check requested domin and your securitytrails key " + string(body))
	}
	for _, j := range records {
		h, ok := j.(map[string]interface{})
		if !ok {
			return nil, errors.New("[!2] : no records found for this domin , check requested domin and your securitytrails key " + string(body))
		}
		t, ok := h["values"].([]interface{})
		if !ok {
			return nil, errors.New("[!3] : no records found for this domin , check requested domin and your securitytrails key " + string(body))
		}
		for _, v := range t {
			x, ok := v.(map[string]interface{})
			if !ok {
				return nil, errors.New("[!4] : no records found for this domin , check requested domin and your securitytrails key " + string(body))
			}
			ips = append(ips, x["ip"].(string))
		}

	}
	ips = RemoveDuplicateStr(ips)
	return ips, nil
}
func mxResHandler(resp *http.Response) ([]string, error) {
	var ips []string
	var results map[string]interface{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	// close response body
	resp.Body.Close()
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}
	records, ok := results["records"].([]interface{})
	if !ok {
		return nil, errors.New("[!1] : no records found for this domin , check requested domin and your securitytrails key " + string(body))
	}
	for _, j := range records {
		h, ok := j.(map[string]interface{})
		if !ok {
			return nil, errors.New("[!1] : no records found for this domin , check requested domin and your securitytrails key " + string(body))
		}
		t, ok := h["values"].([]interface{})
		if !ok {
			return nil, errors.New("[!1] : no records found for this domin , check requested domin and your securitytrails key " + string(body))
		}
		for _, v := range t {
			x, ok := v.(map[string]interface{})
			if !ok {
				return nil, errors.New("[!1] : no records found for this domin , check requested domin and your securitytrails key " + string(body))
			}
			//fmt.Println(x["ip"])
			ips = append(ips, x["host"].(string))
		}

	}
	ips = RemoveDuplicateStr(ips)
	return ips, nil
}

func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
