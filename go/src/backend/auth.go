package backend

import "bytes"
import "encoding/json"
import "fmt"
import "io/ioutil"
import "net/http"
import "net/url"
import "log"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "golang.org/x/crypto/bcrypt"

const (
    AUTHORIZATION_SERVER = "http://hydra:4445"
    AUTHORIZATION_REQUEST_ENDPOINT_BASE = "/oauth2/auth/requests"
    LOGIN_ENDPOINT = "/login"
    LOGIN_ACCEPT_ENDPOINT = "/login/accept"
    LOGIN_REJECT_ENDPOINT = "/login/reject"
    CONSENT_ENDPOINT = "/consent"
    CONSENT_ACCEPT_ENDPOINT = "/consent/accept"
    LOGOUT_ENDPOINT = "/logout"
    LOGOUT_ACCEPT_ENDPOINT = "/logout/accept"
    TOKEN_INTROSPECTION_ENDPOINT = "/oauth2/introspect"
)

type LoginResult struct {
    Subject string `json:"subject,omitempty"`
    Remember bool `json:"remember,omitempty"`
    RememberFor int `json:"remember_for,omitempty"`
    Error string `json:"error,omitempty"`
    ErrorDescription string `json:"error_description,omitempty"`
}

type ConsentResult struct {
    GrantScope []string `json:"grant_scope"`
    Remember bool `json:"remember"`
    RememberFor int `json:"remember_for"`
    Session ConsentSessionInfo `json:"session"`
}

type ConsentSessionInfo struct {
    IDToken IDTokenSessionInfo `json:"id_token"`
    // Add AccessToken here if we ever need access_token session info
}

type IDTokenSessionInfo struct {
    Username string `json:"user_name"`
    Email string `json:"email"`
    FirstName string `json:"first_name"`
    LastName string `json:"last_name"`
}

type Challenge struct {
    Challenge string `json:"challenge"`
}

type LoginCredentials struct {
    Username string `json:"username"`
    Password string `json:"password"`
    LoginChallenge string `json:"login_challenge"`
}

type FrontendLoginDisplayParams struct {
    DisplayLoginUI bool `json:"display_login_ui"`
}

type AuthRequest struct {
    AuthToken string `json:"token"`
    Subject string `json:"subject"`
}

type AuthResult struct {
    AuthSuccessful bool `json:"auth_successful"`
}

type HydraCompletedRequest struct {
    RedirectTo string `json:"redirect_to"`
}

type HydraLoginRequest struct {
    Challenge string `json:"challenge"`
    RequestURL string `json:"request_url"`
    RequestedScope []string `json:"requested_scope"`
    SessionID string `json:"session_id"`
    Skip bool `json:"skip"`
    Subject string `json:"subject"`
}

type HydraLogoutRequest struct {
    RequestURL string `json:"request_url"`
    RPInitiated bool `json:"rp_initiated"`
    SID string `json:"sid"`
    Subject string `json:"subject"`
}

type HydraConsentRequest struct {
    ACR string `json:"acr"`
    Challenge string `json:"challenge"`
    LoginChallenge string `json"login_challenge"`
    LoginSessionID string `json:"login_session_id"`
    RequestURL string `json:"request_url"`
    RequestedScope []string `json:"requested_scope"`
    Skip bool `json:"skip"`
    Subject string `json:"subject"`
}

type HydraOauth2TokenIntrospection struct {
    Active bool `json:"active"`
    Audience []string `json:"aud"`
    ClientID string `json:"client_id"`
    Expiration int64 `json:"exp"`
    IssuedAt int64 `json:"iat"`
    IssuerURL string `json:"iss"`
    NotBefore int64 `json:"nbf"`
    ObfuscatedSubject string `json:"obfuscated_subject"`
    Scope string `json:"scope"`
    Subject string `json:"sub"`
    TokenType string `json:"token_type"`
    Username string `json:"username"`
}

type HydraGenericError struct {
    Debug string `json:"debug"`
    ErrorMsg string `json:"error"`
    ErrorDesc string `json:"error_description"`
    StatusCode int `json:"status_code"`
}

func (loginRequest HydraLoginRequest) String() string {
    return fmt.Sprintf("Skip: %t, Subject: %s, RequestURL: %s, RequestedScope: %v",
            loginRequest.Skip,
            loginRequest.Subject,
            loginRequest.RequestURL,
            loginRequest.RequestedScope)
}

func (consentRequest HydraConsentRequest) String() string {
    return fmt.Sprintf("Skip: %t, Subject: %s, Requested Scope: %v",
            consentRequest.Skip,
            consentRequest.Subject,
            consentRequest.RequestedScope)
}

func (loginResult LoginResult) String() string {
    return fmt.Sprintf("Subject: %s, Remember: %t, Remember for: %d, Error: %s, Error description: %s",
            loginResult.Subject,
            loginResult.Remember,
            loginResult.RememberFor,
            loginResult.Error,
            loginResult.ErrorDescription)
}

func checkLoginChallenge(challenge string) (HydraLoginRequest, error) {
    //challenge = strings.Trim(challenge, "\"")
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
        log.Printf("Error sending login request: %s", err)
        return HydraLoginRequest{}, err
    }
    defer resp.Body.Close()

    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading login request response: %s", err)
        return HydraLoginRequest{}, err
    }

    var loginRequest HydraLoginRequest
    err = json.Unmarshal(respBody, &loginRequest)
    if err != nil {
        log.Printf("Error unmarshalling login request response: %s", err)
        return HydraLoginRequest{}, err
    }

    log.Printf("Received login challenge response: %s", loginRequest)

    return loginRequest, nil
}

func checkUserLogin(
        username string,
        password string,
        loginChallenge string) (bool, string, error) {
    // Get the user record from the DB, if it exists
    userDB, err := sql.Open("mysql", dbConnStr(LOGIN_DB_USER, LOGIN_DB_PW, USER_DB))
	if err != nil {
        log.Printf("Error connecting to users db: %s", err)
        return false, "", err
	}
	defer userDB.Close()

    res := userDB.QueryRow(`SELECT pw_hash, first_name, last_name
            FROM user_info
            WHERE user_name = ?`,
            username)

    var pwHash []byte
    var firstName, lastName string
    err = res.Scan(&pwHash, &firstName, &lastName)
    var authSuccessful bool
    var authResponse string
    if err == sql.ErrNoRows {
        // If there's no user with the given username, reject the login request
        log.Printf("No user found for username %s", username)
        authSuccessful, authResponse, err = completeLoginRequestWithAuthServer(false,
            username, loginChallenge)
        if err != nil {
            log.Printf("Error completing login flow: %s", err)
            return false, "", err
        }
    } else if err != nil {
        log.Printf("Error fetching user record from db: %s", err)
        return false, "", err
    } else {
        // If we're here, we have a valid user row from the DB, check the given password
        // against the stored hash
        err = bcrypt.CompareHashAndPassword(pwHash, []byte(password))
        if err != nil {
            // Password validation failed, let the login endpoint know
            authSuccessful, authResponse, err = completeLoginRequestWithAuthServer(false,
                username, loginChallenge)
            if err != nil {
                log.Printf("Error completing login flow: %s", err)
                return false, "", err
            }
        } else {
            // Password validation succeeded, let the login endpoint know
            authSuccessful, authResponse, err = completeLoginRequestWithAuthServer(true,
                username, loginChallenge)
            if err != nil {
                log.Printf("Error completing login flow: %s", err)
                return false, "", err
            }
        }
    }

    return authSuccessful, authResponse, nil
}

func completeLoginRequestWithAuthServer(
        loginSuccessful bool,
        username string,
        loginChallenge string) (bool, string, error) {

    params := url.Values{}
    params.Set("login_challenge", loginChallenge)

    // Send the request to the authorization backend
    var requestUrl string
    var body LoginResult
    if loginSuccessful {
        requestUrl = AUTHORIZATION_SERVER + AUTHORIZATION_REQUEST_ENDPOINT_BASE +
                LOGIN_ACCEPT_ENDPOINT + "?" + params.Encode()
        body.Subject = username
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
        return false, "", err
    }
    log.Printf("Sending auth request to server at endpoint %s with body %s",
        requestUrl, string(bodyJson))

    req, err := http.NewRequest(http.MethodPut, requestUrl, bytes.NewReader(bodyJson))
    if err != nil {
        log.Printf("Error creating login http request: %s", err)
        return false, "", err
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error sending login request to auth server: %s", err)
        return false, "", err
    }
    defer resp.Body.Close()

    response, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading response body: %s", err)
        return false, "", err
    }

    switch resp.StatusCode {
    case http.StatusOK:
        var completedRequest HydraCompletedRequest
        err = json.Unmarshal(response, &completedRequest)
        if err != nil {
            log.Printf("Error unmarshalling completedRequest: %s", err)
            return false, "", err
        }
        return true, completedRequest.RedirectTo, nil
    case http.StatusUnauthorized:
        fallthrough
    case http.StatusNotFound:
        fallthrough
    case http.StatusInternalServerError:
        fallthrough
    default:
        var genericError HydraGenericError
        err = json.Unmarshal(response, &genericError)
        if err != nil {
            log.Printf("Error unmarshalling genericError: %s", err)
            return false, "", err
        }
        return false, genericError.ErrorMsg, nil
    }
}

func completeConsentRequestWithAuthServer(
        consentSuccessful bool,
        consentRequest HydraConsentRequest) (bool, string, error) {

    // Get the user record from the DB so we can populate extra user info
    // in the consent acceptance (this will later be returned from requests to the
    // openid-connect userinfo endpoint)
    userDB, err := sql.Open("mysql", dbConnStr(LOGIN_DB_USER, LOGIN_DB_PW, USER_DB))
	if err != nil {
        log.Printf("Error connecting to users db: %s", err)
        return false, "", err
	}
	defer userDB.Close()

    res := userDB.QueryRow(`SELECT first_name, last_name, email
            FROM user_info
            WHERE user_name = ?`,
            consentRequest.Subject)

    idToken := IDTokenSessionInfo{Username: consentRequest.Subject}
    err = res.Scan(&idToken.FirstName, &idToken.LastName, &idToken.Email)
    if err != nil {
        log.Printf("Error retrieving user info record from db: %s", err)
        return false, "", err
    }

    // For now, we just blindly consent to everything, since we're running
    // both the app and the authorization provider, there's no need to require
    // the user to consent
    // TODO: Implement the consent reject path
    params := url.Values{}
    params.Set("consent_challenge", consentRequest.Challenge)
    requestUrl := AUTHORIZATION_SERVER + AUTHORIZATION_REQUEST_ENDPOINT_BASE +
            CONSENT_ACCEPT_ENDPOINT + "?" + params.Encode()
    log.Printf("Granting scopes: %v", consentRequest.RequestedScope)
    body := ConsentResult{
        GrantScope: consentRequest.RequestedScope,
        Remember: true,
        RememberFor: 10,
        Session: ConsentSessionInfo{IDToken: idToken},}

    bodyJson, err := json.Marshal(body)
    if err != nil {
        log.Printf("Error marshalling consent JSON body: %s", err)
        return false, "", err
    }

    log.Printf("Sending consent request to server at endpoint %s with body %s",
            requestUrl, string(bodyJson))

    req, err := http.NewRequest(http.MethodPut, requestUrl, bytes.NewReader(bodyJson))
    if err != nil {
        log.Printf("Error creating consent http request: %s", err)
        return false, "", err
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error sending consent request to auth server: %s", err)
        return false, "", err
    }
    defer resp.Body.Close()

    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading response body: %s", err)
        return false, "", err
    }

    log.Printf("Auth server response: %s", string(respBody))

    switch resp.StatusCode {
    case http.StatusOK:
        var completedRequest HydraCompletedRequest
        err = json.Unmarshal(respBody, &completedRequest)
        if err != nil {
            log.Printf("Error unmarshalling completedRequest: %s", err)
            return false, "", err
        }
        return true, completedRequest.RedirectTo, nil
    case http.StatusNotFound:
        fallthrough
    case http.StatusInternalServerError:
        fallthrough
    default:
        var genericError HydraGenericError
        err = json.Unmarshal(respBody, &genericError)
        if err != nil {
            log.Printf("Error unmarshalling genericError: %s", err)
            return false, "", err
        }
        return false, genericError.ErrorMsg, nil
    }
}

func completeLogoutRequestWithAuthServer(logoutChallenge string) (bool, string, error) {
    // For now, no reason to implement the logout reject path, I can't think
    // of how that would be useful
    params := url.Values{}
    params.Set("logout_challenge", logoutChallenge)
    requestUrl := AUTHORIZATION_SERVER + AUTHORIZATION_REQUEST_ENDPOINT_BASE +
            LOGOUT_ACCEPT_ENDPOINT + "?" + params.Encode()

    log.Printf("Sending logout request to server at endpoint %s",
            requestUrl)

    req, err := http.NewRequest(http.MethodPut, requestUrl, http.NoBody)
    if err != nil {
        log.Printf("Error creating logout http request: %s", err)
        return false, "", err
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error sending logout request to auth server: %s", err)
        return false, "", err
    }
    defer resp.Body.Close()

    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading response body: %s", err)
        return false, "", err
    }

    log.Printf("Auth server response: %s", string(respBody))

    switch resp.StatusCode {
    case http.StatusOK:
        var completedRequest HydraCompletedRequest
        err = json.Unmarshal(respBody, &completedRequest)
        if err != nil {
            log.Printf("Error unmarshalling completedRequest: %s", err)
            return false, "", err
        }
        return true, completedRequest.RedirectTo, nil
    case http.StatusNotFound:
        fallthrough
    case http.StatusInternalServerError:
        fallthrough
    default:
        var genericError HydraGenericError
        err = json.Unmarshal(respBody, &genericError)
        if err != nil {
            log.Printf("Error unmarshalling genericError: %s", err)
            return false, "", err
        }
        return false, genericError.ErrorMsg, nil
    }
}

func checkConsentChallenge(consentChallenge string) (HydraConsentRequest, error) {
    params := url.Values{}
    params.Set("consent_challenge", consentChallenge)

    // First, we need to interrogate the authorization server for the details
    // of the consent request
    requestUrl := AUTHORIZATION_SERVER + AUTHORIZATION_REQUEST_ENDPOINT_BASE +
            CONSENT_ENDPOINT + "?" + params.Encode()
    log.Printf("Sending consent challenge request %s", requestUrl)

    resp, err := http.Get(requestUrl)
    if err != nil {
        log.Printf("Error sending consent request: %s", err)
        return HydraConsentRequest{}, err
    }
    defer resp.Body.Close()

    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading consent request response: %s", err)
        return HydraConsentRequest{}, err
    }

    var consentRequest HydraConsentRequest
    err = json.Unmarshal(respBody, &consentRequest)
    if err != nil {
        log.Printf("Error unmarshalling consent request response: %s", err)
        return HydraConsentRequest{}, err
    }

    log.Printf("Received consent consent response: %s", consentRequest)

    return consentRequest, nil
}

func checkLogoutChallenge(logoutChallenge string) (HydraLogoutRequest, error) {
    params := url.Values{}
    params.Set("logout_challenge", logoutChallenge)

    // First, we need to interrogate the authorization server for the details
    // of the logout request
    requestUrl := AUTHORIZATION_SERVER + AUTHORIZATION_REQUEST_ENDPOINT_BASE +
            LOGOUT_ENDPOINT + "?" + params.Encode()
    log.Printf("Sending logout challenge request %s", requestUrl)

    resp, err := http.Get(requestUrl)
    if err != nil {
        log.Printf("Error sending logout request: %s", err)
        return HydraLogoutRequest{}, err
    }
    defer resp.Body.Close()

    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading logout request response: %s", err)
        return HydraLogoutRequest{}, err
    }

    var logoutRequest HydraLogoutRequest
    err = json.Unmarshal(respBody, &logoutRequest)
    if err != nil {
        log.Printf("Error unmarshalling logout request response: %s", err)
        return HydraLogoutRequest{}, err
    }

    log.Printf("Received logout consent response: %s", logoutRequest)

    return logoutRequest, nil
}

func authorizeToken(
        subject string,
        token string,
        done <-chan interface{},
        respChan chan<- ResponseMessage) bool {

    requestBody := url.Values{}
    requestBody.Set("token", token)
    requestUrl := AUTHORIZATION_SERVER + TOKEN_INTROSPECTION_ENDPOINT

    log.Printf("Sending token introspection with body %v", requestBody);

    resp, err := http.PostForm(requestUrl, requestBody)
    if err != nil {
        log.Printf("Error sending token introspection request: %s", err)
        return false
    }

    defer resp.Body.Close()

    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading response body: %s", err)
        return false
    }

    log.Printf("Auth server response: %s", string(respBody))

    switch resp.StatusCode {
    case http.StatusOK:
        var tokenIntrospection HydraOauth2TokenIntrospection
        err = json.Unmarshal(respBody, &tokenIntrospection)
        if err != nil {
            log.Printf("Error unmarshalling token introspection: %s", err)
            return false
        }

        authResult := AuthResult{}
        if tokenIntrospection.Active && tokenIntrospection.Subject == subject {
            authResult.AuthSuccessful = true
        } else {
            authResult.AuthSuccessful = false
        }

        select {
        case <-done:
        case respChan <- ResponseMessage{Type: AuthUserResponse, Value: authResult}:
        }
        return authResult.AuthSuccessful

    case http.StatusUnauthorized:
        fallthrough
    case http.StatusInternalServerError:
        fallthrough
    default:
        var genericError HydraGenericError
        err = json.Unmarshal(respBody, &genericError)
        if err != nil {
            log.Printf("Error unmarshalling genericError: %s", err)
            return false
        }
    }

    return false
}

func sendHttpError(resp http.ResponseWriter) {
    resp.WriteHeader(http.StatusInternalServerError)
}

func sendHttpAuthErrorMsg(resp http.ResponseWriter, errorMsg string) {
    // TODO: Do something better here
    resp.WriteHeader(http.StatusForbidden)
}

func sendHttpRedirect(resp http.ResponseWriter, redirectTo string) {
    resp.Header().Set("Location", redirectTo)
    log.Printf("Sending redirect to %s", redirectTo)
    resp.WriteHeader(http.StatusFound)
}

func enableCORS(resp *http.ResponseWriter, req *http.Request) bool {
    (*resp).Header().Set("Access-Control-Allow-Origin", "*")
    (*resp).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

    if req.Method == http.MethodOptions {
        (*resp).WriteHeader(http.StatusOK)
        return true
    }
    return false
}

func HandleLoginChallenge(resp http.ResponseWriter, req *http.Request) {
    defer req.Body.Close()

    preflight := enableCORS(&resp, req)
    if preflight {
        return
    }

    request, err := ioutil.ReadAll(req.Body)
    if err != nil {
        log.Printf("Error reading login challenge request: %s", err)
        sendHttpError(resp)
        return
    }

    var challenge Challenge
    err = json.Unmarshal(request, &challenge)
    if err != nil {
        log.Printf("Error unmarshalling login challenge request: %s", err)
        sendHttpError(resp)
        return
    }

    loginRequest, err := checkLoginChallenge(challenge.Challenge)
    if err != nil {
        log.Print(err)
        sendHttpError(resp)
        return
    }

    // The server can respond that the user queried for has already authenticated,
    // and so we can skip asking the user for authentication credentials
    // (skip==true), or can request that we authenticate the user (skip==false)
    if loginRequest.Skip {
        // The user has already successfully authenticated, so complete the
        // authentication with the backend
        authSuccessful, authResponse, err := completeLoginRequestWithAuthServer(true,
            loginRequest.Subject,
            challenge.Challenge)
        if err != nil {
            log.Printf("Error completing skipped authentication: %s", err)
            sendHttpError(resp)
            return
        }

        if authSuccessful {
            sendHttpRedirect(resp, authResponse)
        } else {
            sendHttpAuthErrorMsg(resp, authResponse)
        }
    } else {
        // The user needs to authenticate, let the frontend know that it needs
        // to display an authentication UI
        resp.Header().Set("Content-Type", "application/json")
        respBody := FrontendLoginDisplayParams{DisplayLoginUI: true}
        respBodyJson, err := json.Marshal(respBody)
        if err != nil {
            log.Printf("Error marshalling login display params json: %s", err)
            sendHttpError(resp)
            return
        }
        _, err = resp.Write(respBodyJson)
        if err != nil {
            log.Printf("Error writing login display params response: %s", err)
        }
    }
}

func HandleLoginCredentials(resp http.ResponseWriter, req *http.Request) {
    defer req.Body.Close()

    preflight := enableCORS(&resp, req)
    if preflight {
        return
    }

    request, err := ioutil.ReadAll(req.Body)
    if err != nil {
        log.Printf("Error reading login credentials: %s", err)
        sendHttpError(resp)
        return
    }

    var loginCredentials LoginCredentials
    err = json.Unmarshal(request, &loginCredentials)
    if err != nil {
        log.Printf("Error unmarshalling login credentials: %s", err)
        sendHttpError(resp)
        return
    }

    authSuccessful, authResponse, err := checkUserLogin(loginCredentials.Username,
        loginCredentials.Password, loginCredentials.LoginChallenge)
    if err != nil {
        log.Print(err)
        sendHttpError(resp)
        return
    }

    if authSuccessful {
        sendHttpRedirect(resp, authResponse)
    } else {
        sendHttpAuthErrorMsg(resp, authResponse)
    }
}

func HandleConsentChallenge(resp http.ResponseWriter, req *http.Request) {
    defer req.Body.Close()

    request, err := ioutil.ReadAll(req.Body)
    if err != nil {
        log.Printf("Error reading consent challenge request: %s", err)
        sendHttpError(resp)
        return
    }

    var challenge Challenge
    err = json.Unmarshal(request, &challenge)
    if err != nil {
        log.Printf("Error unmarshalling consent challenge request: %s", err)
        sendHttpError(resp)
        return
    }

    consentRequest, err := checkConsentChallenge(challenge.Challenge)
    if err != nil {
        log.Print(err)
        sendHttpError(resp)
        return
    }

    // For now, we blindly accept all consent requests, since we're controlling
    // both the app and the authorization service, we assume we want to grant
    // access to everything.  Potentially implement this more fully in the future
    consentSuccessful, consentResponse, err := completeConsentRequestWithAuthServer(true,
        consentRequest)
    if err != nil {
        log.Printf("Error completing consent: %s", err)
        sendHttpError(resp)
        return
    }

    if consentSuccessful {
        sendHttpRedirect(resp, consentResponse)
    } else {
        sendHttpAuthErrorMsg(resp, consentResponse)
    }
}

func HandleLogoutChallenge(resp http.ResponseWriter, req *http.Request) {
    defer req.Body.Close()

    request, err := ioutil.ReadAll(req.Body)
    if err != nil {
        log.Printf("Error reading logout challenge request: %s", err)
        sendHttpError(resp)
        return
    }

    var challenge Challenge
    err = json.Unmarshal(request, &challenge)
    if err != nil {
        log.Printf("Error unmarshalling logout challenge request: %s", err)
        sendHttpError(resp)
        return
    }

    logoutRequest, err := checkLogoutChallenge(challenge.Challenge)
    log.Printf("Logging out session %s", logoutRequest.SID)
    if err != nil {
        log.Print(err)
        sendHttpError(resp)
        return
    }

    // Just accept the logout request, nothing else to do here
    logoutSuccessful, logoutResponse, err := completeLogoutRequestWithAuthServer(
        challenge.Challenge)
    if err != nil {
        log.Printf("Error completing logout: %s", err)
        sendHttpError(resp)
        return
    }

    if logoutSuccessful {
        sendHttpRedirect(resp, logoutResponse)
    } else {
        sendHttpAuthErrorMsg(resp, logoutResponse)
    }
}

