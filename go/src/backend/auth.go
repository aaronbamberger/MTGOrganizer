package backend

import "encoding/json"
import "fmt"
import "net/http"
import "net/url"
import "log"

const (
    AUTHORIZATION_SERVER = "http://172.18.0.8:4445"
    AUTHORIZATION_REQUEST_ENDPOINT_BASE = "/oauth2/auth/requests"
    LOGIN_ENDPOINT = "/login"
    LOGIN_ACCEPT_ENDPOINT = "/login/accept"
    CONSENT_ENDPOINT = "/consent"
    CONSENT_ACCEPT_ENDPOINT = "/consent/accept"
)

type LoginResponseParams struct {
    Skip bool `json:"skip"`
    Subject string `json:"subject"`
    RequestURL string `json:"subject"`
    RequestedScope []string `json:"requested_scope"`
}

func (responseParams LoginResponseParams) String() string {
    return fmt.Sprintf("Skip: %t, Subject: %s, RequestURL: %s, RequestedScope: %v",
        responseParams.Skip,
        responseParams.Subject,
        responseParams.RequestURL,
        responseParams.RequestedScope)
}

func checkLoginChallenge(challenge string,
        done <-chan interface{},
        respChan chan<- ResponseMessage) {

    // Construct the challenge message to the authorization backend
    log.Printf("Received login challenge %s", challenge)
    params := url.Values{}
    params.Set("login_challenge", challenge)
    // Send the request to the authorization backend
    requestUrl := AUTHORIZATION_SERVER + AUTHORIZATION_REQUEST_ENDPOINT_BASE +
            LOGIN_ENDPOINT + "?" + params.Encode()
    log.Printf("Sending login challenge request %s", requestUrl)

    resp, err := http.Get(requestUrl)
    if err != nil {
        sendError(done, respChan, err)
        return
    }
    defer resp.Body.Close()

    jsonDecoder := json.NewDecoder(resp.Body)
    var loginResponse LoginResponseParams
    err = jsonDecoder.Decode(&loginResponse)
    if err != nil {
        sendError(done, respChan, err)
        return
    }

    log.Printf("Received login challenge response: %s", loginResponse)

    responseMessage := ResponseMessage{
        Type: LoginChallengeResponse,
        Value: loginResponse}

    select {
    case <-done:
    case respChan <- responseMessage:
    }
}

