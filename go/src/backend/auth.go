package backend

import "bytes"
import "encoding/json"
import "fmt"
import "io/ioutil"
import "net/http"
import "net/url"
import "strings"
import "log"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "golang.org/x/crypto/bcrypt"

const (
    AUTHORIZATION_SERVER = "http://172.18.0.6:4445"
    AUTHORIZATION_REQUEST_ENDPOINT_BASE = "/oauth2/auth/requests"
    LOGIN_ENDPOINT = "/login"
    LOGIN_ACCEPT_ENDPOINT = "/login/accept"
    LOGIN_REJECT_ENDPOINT = "/login/reject"
    CONSENT_ENDPOINT = "/consent"
    CONSENT_ACCEPT_ENDPOINT = "/consent/accept"
)

type LoginChallengeParams struct {
    Skip bool `json:"skip"`
    Subject string `json:"subject"`
    RequestURL string `json:"subject"`
    RequestedScope []string `json:"requested_scope"`
}

type ConsentChallengeParams struct {
    Skip bool `json:"skip"`
    Subject string `json:"subject"`
    RequestedScope []string `json:"requested_scope"`
    RequestedAccessTokenAudience []string `json:"requested_access_token_audience"`
}

type LoginCredentials struct {
    Username string `json:"username"`
    Password string `json:"password"`
    LoginChallenge string `json:"login_challenge"`
}

type LoginResult struct {
    Subject string `json:"subject,omitempty"`
    Remember bool `json:"remember,omitempty"`
    RememberFor int `json:"remember_for,omitempty"`
    Error string `json:"error,omitempty"`
    ErrorDescription string `json:"error_description,omitempty"`
}

type ConsentResult struct {
    GrantScope []string `json:"grant_scope"`
    GrantAccessTokenAudience []string `json:"grant_access_token_audience"`
    Remember bool `json:"remember"`
    RememberFor int `json:"remember_for"`
}

func (challengeParams LoginChallengeParams) String() string {
    return fmt.Sprintf("Skip: %t, Subject: %s, RequestURL: %s, RequestedScope: %v",
            challengeParams.Skip,
            challengeParams.Subject,
            challengeParams.RequestURL,
            challengeParams.RequestedScope)
}

func (consentParams ConsentChallengeParams) String() string {
    return fmt.Sprintf("Skip: %t, Subject: %s, Requested Scope: %v, Requested Audience: %v",
            consentParams.Skip,
            consentParams.Subject,
            consentParams.RequestedScope,
            consentParams.RequestedAccessTokenAudience)
}

func (loginResult LoginResult) String() string {
    return fmt.Sprintf("Subject: %s, Remember: %t, Remember for: %d, Error: %s, Error description: %s",
            loginResult.Subject,
            loginResult.Remember,
            loginResult.RememberFor,
            loginResult.Error,
            loginResult.ErrorDescription)
}

func checkLoginChallenge(challenge string,
        done <-chan interface{},
        respChan chan<- ResponseMessage) {
    challenge = strings.Trim(challenge, "\"")
    log.Printf("Received login challenge %s", challenge)
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
    var loginResponse LoginChallengeParams
    err = jsonDecoder.Decode(&loginResponse)
    if err != nil {
        sendError(done, respChan, err)
        return
    }

    log.Printf("Received login challenge response: %s", loginResponse)

    // If the login challenge response tells us we need to authenticate the user,
    // then send the message back to the frontend so it can display an
    // authentication UI.  If the response tells us the user has already
    // authenticated successfully, then just accept the login
    if loginResponse.Skip {
        creds := LoginCredentials{
            Username: loginResponse.Subject,
            Password: "",
            LoginChallenge: challenge}
        acceptUserLogin(creds, done, respChan)
    } else {
        responseMessage := ResponseMessage{
            Type: LoginChallengeResponse,
            Value: loginResponse}

        select {
        case <-done:
        case respChan <- responseMessage:
        }
    }
}

func checkUserLogin(loginCredentials string,
        done <-chan interface{},
        respChan chan<- ResponseMessage) {
    // First, make sure we can unserialize the login credentials JSON
    log.Printf("Received login credentials %s", loginCredentials)
    var creds LoginCredentials
    err := json.Unmarshal([]byte(loginCredentials), &creds)
    if err != nil {
        log.Printf("Error deserializing login credentials: %s", err)
        sendError(done, respChan, err)
        return
    }

    // Get the user record from the DB, if it exists
    userDB, err := sql.Open("mysql", dbConnStr(LOGIN_DB_USER, LOGIN_DB_PW, USER_DB))
	if err != nil {
        log.Printf("Error connecting to users db: %s", err)
		sendError(done, respChan, err)
        return
	}
	defer userDB.Close()

    res := userDB.QueryRow(`SELECT pw_hash, first_name, last_name
            FROM user_info
            WHERE user_name = ?`,
            creds.Username)

    var pwHash []byte
    var firstName, lastName string
    err = res.Scan(&pwHash, &firstName, &lastName)
    var authServerResponse string
    if err == sql.ErrNoRows {
        // If there's no user with the given username, reject the login request
        log.Printf("No user found for username %s", creds.Username)
        authServerResponse = completeLoginRequestWithAuthServer(false, creds)
    } else if err != nil {
        log.Printf("Error fetching user record from db: %s", err)
        sendError(done, respChan, err)
        return
    } else {
        // If we're here, we have a valid user row from the DB, check the given password
        // against the stored hash
        err = bcrypt.CompareHashAndPassword(pwHash, []byte(creds.Password))
        if err != nil {
            // Password validation failed, let the login endpoint know
            authServerResponse = completeLoginRequestWithAuthServer(false, creds)
        } else {
            // Password validation succeeded, let the login endpoint know
            authServerResponse = completeLoginRequestWithAuthServer(true, creds)
        }
    }

    responseMessage := ResponseMessage{
        Type: LoginResponse,
        Value: authServerResponse}

    select {
    case <-done:
    case respChan <- responseMessage:
    }
}

func acceptUserLogin(creds LoginCredentials,
        done <-chan interface{},
        respChan chan<- ResponseMessage) {

    authServerResponse := completeLoginRequestWithAuthServer(true, creds)
    responseMessage := ResponseMessage{
        Type: LoginResponse,
        Value: authServerResponse}

    select {
    case <-done:
    case respChan <- responseMessage:
    }
}

func completeLoginRequestWithAuthServer(loginSuccessful bool,
        creds LoginCredentials) string {

    params := url.Values{}
    params.Set("login_challenge", creds.LoginChallenge)
    // Send the request to the authorization backend
    var requestUrl string
    var body LoginResult
    if loginSuccessful {
        requestUrl = AUTHORIZATION_SERVER + AUTHORIZATION_REQUEST_ENDPOINT_BASE +
                LOGIN_ACCEPT_ENDPOINT + "?" + params.Encode()
        body.Subject = creds.Username
        body.Remember = true
        body.RememberFor = 120
    } else {
        requestUrl = AUTHORIZATION_SERVER + AUTHORIZATION_REQUEST_ENDPOINT_BASE +
                LOGIN_REJECT_ENDPOINT + "?" + params.Encode()
        body.Error = "access_denied"
        body.ErrorDescription = "The Username or Password is incorrect"
    }

    bodyJson, err := json.Marshal(body)
    if err != nil {
        log.Printf("Error marshalling login JSON body: %s", err)
    }

    log.Printf("Sending auth request to server at endpoint %s with body %s",
            requestUrl, string(bodyJson))

    req, err := http.NewRequest(http.MethodPut, requestUrl, bytes.NewReader(bodyJson))
    if err != nil {
        log.Printf("Error creating login http request: %s", err)
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error sending login request to auth server: %s", err)
    }
    defer resp.Body.Close()

    response, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading response body: %s", err)
    }

    log.Printf("Auth server response: %s", string(response))

    return string(response)
}

func checkUserConsent(consentChallenge string,
        done <-chan interface{},
        respChan chan<- ResponseMessage) {

    consentChallenge = strings.Trim(consentChallenge, "\"")
    log.Printf("Received consent challenge %s", consentChallenge)

    params := url.Values{}
    params.Set("consent_challenge", consentChallenge)

    // First, we need to interrogate the authorization server for the details
    // of the consent request
    requestUrl := AUTHORIZATION_SERVER + AUTHORIZATION_REQUEST_ENDPOINT_BASE +
            CONSENT_ENDPOINT + "?" + params.Encode()
    log.Printf("Sending consent challenge request %s", requestUrl)

    resp, err := http.Get(requestUrl)
    if err != nil {
        sendError(done, respChan, err)
        return
    }
    defer resp.Body.Close()

    jsonDecoder := json.NewDecoder(resp.Body)
    var consentResponse ConsentChallengeParams
    err = jsonDecoder.Decode(&consentResponse)
    if err != nil {
        sendError(done, respChan, err)
        return
    }

    log.Printf("Received consent challenge response: %s", consentResponse)

    // For now, we just blindly consent to everything, since we're running
    // both the app and the authorization provider, there's no need to require
    // the user to consent
    requestUrl = AUTHORIZATION_SERVER + AUTHORIZATION_REQUEST_ENDPOINT_BASE +
            CONSENT_ACCEPT_ENDPOINT + "?" + params.Encode()
    body := ConsentResult{
        GrantScope: consentResponse.RequestedScope,
        GrantAccessTokenAudience: consentResponse.RequestedAccessTokenAudience,
        Remember: true,
        RememberFor: 120}

    bodyJson, err := json.Marshal(body)
    if err != nil {
        log.Printf("Error marshalling login JSON body: %s", err)
    }

    log.Printf("Sending consent request to server at endpoint %s with body %s",
            requestUrl, string(bodyJson))

    req, err := http.NewRequest(http.MethodPut, requestUrl, bytes.NewReader(bodyJson))
    if err != nil {
        log.Printf("Error creating consent http request: %s", err)
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp2, err := client.Do(req)
    if err != nil {
        log.Printf("Error sending consent request to auth server: %s", err)
    }
    defer resp2.Body.Close()

    response, err := ioutil.ReadAll(resp2.Body)
    if err != nil {
        log.Printf("Error reading response body: %s", err)
    }

    log.Printf("Auth server response: %s", string(response))

    responseMessage := ResponseMessage{
        Type: ConsentResponse,
        Value: string(response)}

    select {
    case <-done:
    case respChan <- responseMessage:
    }
}
