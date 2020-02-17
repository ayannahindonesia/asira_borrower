package custommodule

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/jarcoal/httpmock"
)

type (
	//Messaging hold properties for SMS and notification
	Messaging struct {
		Key       string
		Secret    string
		Token     string
		Expire    time.Time
		URL       string
		Endpoints MessagingEndpoints
	}

	MessagingEndpoints struct {
		ClientAuth       string
		SMS              string
		PushNotification string
		ListNotification string
	}
)

// SetConfig func
func (model *Messaging) SetConfig(key string, secret string, URL string, Endpoints MessagingEndpoints) {
	model.Key = key
	model.Secret = secret
	model.URL = URL
	model.Endpoints = Endpoints
}

// ClientAuth func
func (model *Messaging) ClientAuth() (err error) {

	//fmt.Println("model.Key ==> ", model.Key, model.Secret)
	basicToken := base64.StdEncoding.EncodeToString([]byte(model.Key + ":" + model.Secret))
	request, _ := http.NewRequest("GET", model.URL+model.Endpoints.ClientAuth, nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Basic %s", basicToken))

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	buffer, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var parseresponse map[string]interface{}
	json.Unmarshal([]byte(buffer), &parseresponse)

	expires := 1200
	if parseresponse["expires_in"] != nil {
		expires = int(parseresponse["expires_in"].(float64))
	}
	model.Expire = time.Now().Local().Add(time.Second * time.Duration(expires))

	if parseresponse["token"] != nil {
		model.Token = parseresponse["token"].(string)
	} else {
		return fmt.Errorf("no token detected")
	}

	return nil
}

// RefreshToken func
func (model *Messaging) RefreshToken() (err error) {
	if time.Now().After(model.Expire) {
		if flag.Lookup("test.v") == nil {
			err = model.ClientAuth()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// SendSMS func
func (model *Messaging) SendSMS(number string, message string) (err error) {
	err = model.RefreshToken()
	if err != nil {
		return err
	}

	///if under test (go test *)
	if flag.Lookup("test.v") != nil {

		//activate mockup response
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		log.Println("run under go test")

		//Mockup responder
		httpmock.RegisterResponder("POST", model.URL+model.Endpoints.SMS,
			func(req *http.Request) (*http.Response, error) {

				//fake response
				result := "{\"status\":\"success\"}"
				resp, err := httpmock.NewJsonResponse(http.StatusOK, result)
				if err != nil {
					return httpmock.NewStringResponse(http.StatusInternalServerError, err.Error()), nil
				}
				return resp, nil
			},
		)
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"phone_number": number,
		"message":      message,
	})

	request, _ := http.NewRequest("POST", model.URL+model.Endpoints.SMS, bytes.NewBuffer(payload))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", model.Token))

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(response.Body)
		log.Printf("Failed sending sms : %s", string(body))

		return fmt.Errorf("Failed sending SMS " + string(body))
	}

	return nil
}

//SendNotificationByToken sending firebase message by token
func (model *Messaging) SendNotificationByToken(title string, message_body string, map_data map[string]string, firebase_token string, recipient_id string) (responseBody []byte, err error) {

	//bug cycling call dependency
	//topics := asira.App.Config.GetStringMap(fmt.Sprintf("%s.messaging.push_notification", asira.App.ENV))

	err = model.RefreshToken()
	if err != nil {
		return nil, err
	}
	if firebase_token == "" {
		firebase_token = model.Endpoints.PushNotification
	}
	payload, _ := json.Marshal(map[string]interface{}{
		"recipient_id":   recipient_id,
		"title":          title,
		"message_body":   message_body,
		"firebase_token": firebase_token,
		"data":           map_data,
	})

	request, _ := http.NewRequest("POST", model.URL+model.Endpoints.PushNotification, bytes.NewBuffer(payload))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", model.Token))

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	log.Println("PUSH NOTIF : ", response)
	//cek response
	responseBody, err = ioutil.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK {

		log.Printf("Failed sending notification : %s", string(responseBody))

		return responseBody, fmt.Errorf("Failed sending notification")
	}
	return responseBody, nil
}
