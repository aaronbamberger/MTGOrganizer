import React from 'react';
import {BACKEND_HOSTNAME,
        LOGIN_CHALLENGE_ENDPOINT,
        LOGIN_CREDS_ENDPOINT} from './Constants.js'

class LoginPage extends React.Component {
  constructor(props) {
    super(props);

    let urlParams = new URLSearchParams(window.location.search);

    this.state = {
      login_needed: false,
      username: '',
      password: '',
      login_challenge: urlParams.get('login_challenge'),
      user: null,
    }

    this.handleUsername = this.handleUsername.bind(this);
    this.handlePassword = this.handlePassword.bind(this);
    this.handleLogin = this.handleLogin.bind(this);
  }

  componentDidMount() {
    this.sendLoginChallengeRequest();
  }

  sendLoginChallengeRequest() {
    const req = new Request("http://" + BACKEND_HOSTNAME + LOGIN_CHALLENGE_ENDPOINT,
      {
        method: "POST",
        body: JSON.stringify({"challenge": this.state.login_challenge}),
      }
    );
    fetch(req).then(response => {
      if (response.redirected) {
        window.location.href = response.url
      } else {
        response.json().then(body => {
          if (body["display_login_ui"]) {
            this.setState({login_needed: true});
          }
        });
      }
    });
  }

  sendLoginCredentials() {
    const req = new Request("http://" + BACKEND_HOSTNAME + LOGIN_CREDS_ENDPOINT,
      {
        method: "POST",
        body: JSON.stringify({
          "username": this.state.username,
          "password": this.state.password,
          "login_challenge": this.state.login_challenge,
          }),
        redirect: "follow",
        mode: "cors",
      }
    );
    fetch(req).then(response => {
      console.log(response);
      if (response.redirected) {
        console.log(response.url);
        for (const header of response.headers) {
          console.log(header);
          //console.log(header + ": " + response.headers.get(header));
        }
        //console.log(response.headers);
        window.location.href = response.url;
      }
    });
  }

  receiveLoginChallengeResult(result) {
    console.log('Login challenge result: ' + result);
    if (result['skip'] === false) {
      this.setState({login_needed: true});
    }
  }

  receiveLoginResponse(response) {
    console.log('Login response: ' + response);
    let responseJSON = JSON.parse(response);
    window.location.replace(responseJSON['redirect_to']);
  }

  handleUsername(event) {
    this.setState({username: event.target.value});
  }

  handlePassword(event) {
    this.setState({password: event.target.value});
  }

  handleLogin(event) {
    this.sendLoginCredentials()
    event.preventDefault();
  }

  render() {
    console.log(this.state);
    if (!this.state.login_needed) {
      return '';
    } else {
      return (
        <div id="login_form">
          <form onSubmit={this.handleLogin}>
            <label>
              Username:
              <input type="text"
                value={this.state.username}
                onChange={this.handleUsername} />
            </label>
            <label>
              Password:
              <input type="password"
                value={this.state.password}
                onChange={this.handlePassword} />
            </label>
            <input type="submit" value="Login" />
          </form>
        </div>
      );
    }
  }
}

export default LoginPage;
