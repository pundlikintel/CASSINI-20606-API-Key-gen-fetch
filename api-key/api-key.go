package api_key

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func GetApiKeys(ctx context.Context, jwt, env, serviceId string) ([]string, error) {
	url := fmt.Sprintf("https://amber-%s-us-east-1-internal-cp.project-amber-smas.com/tenant-management/v1/services/%s/api-clients", env, serviceId)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("x-amber-api-token", jwt)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		log.Error("Failed to read request body")
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error("Failed to read request body")
	}

	data := []map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0)
	if len(data) == 0 {
		return nil, nil
	} else {
		for _, d := range data {
			id := d["id"].(string)
			var key string
			key, err = GetKey(ctx, jwt, env, serviceId, id)
			keys = append(keys, key)
		}
	}
	return keys, nil
}

func GetKey(tx context.Context, jwt, env, serviceId, id string) (string, error) {
	url := fmt.Sprintf("https://amber-%s-us-east-1-internal-cp.project-amber-smas.com/tenant-management/v1/services/%s/api-clients/%s", env, serviceId, id)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("x-amber-api-token", jwt)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		log.Error("Failed to read request body")
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error("Failed to read request body")
	}

	dataK := map[string]interface{}{}
	err = json.Unmarshal(body, &dataK)
	if err != nil {
		return "", err
	}
	if len(dataK) == 0 {
		return "", errors.New("Error in get data")
	} else {
		kys := dataK["keys"].([]interface{})
		return kys[0].(string), nil
	}

	return "", nil
}

func CreateApiKey(ctx context.Context, jwt, env, serviceId, productId string) (string, error) {
	url := fmt.Sprintf("https://amber-%s-us-east-1-internal-cp.project-amber-smas.com/tenant-management/v1/services/%s/api-clients", env, serviceId)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	request := map[string]string{"product_id": productId, "name": "parter_api_4"}
	reqByt, err := json.Marshal(request)
	if err != nil {
		log.Error("Failed to read request body")
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqByt))
	if err != nil {
		log.Error("Failed to read make call")
		return "", err
	}
	req.Header.Set("x-amber-api-token", jwt)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Error("Failed to read request body")
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error("Failed to read request body")
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}
	if er, ok := data["error"]; ok {
		return "", errors.New(er.(string))
	}
	if len(data) == 0 {
		return "", errors.New("Error:key not found")
	} else {
		keysI := data["keys"].([]interface{})
		return keysI[0].(string), nil
	}
}
