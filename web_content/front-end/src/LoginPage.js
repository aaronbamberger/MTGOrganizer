import React from 'react';
import Oidc from 'oidc-client';

class LoginPage extends React.Component {
  constructor(props) {
    super(props);

    let urlParams = new URLSearchParams(window.location.search);

    console.log(urlParams.get('login_challenge'));

    this.state = {
      login_needed: false,
      username: '',
      password: '',
      login_challenge: urlParams.get('login_challenge'),
      user: null,
    }

    this.receiveLoginChallengeResult = this.receiveLoginChallengeResult.bind(this);
    this.receiveLoginResponse = this.receiveLoginResponse.bind(this);
    this.handleUsername = this.handleUsername.bind(this);
    this.handlePassword = this.handlePassword.bind(this);
    this.handleLogin = this.handleLogin.bind(this);
  }

  componentDidMount() {
    this.sendLoginChallengeRequest()
  }

  sendLoginChallengeRequest() {
    const request = JSON.stringify({
      "type": this.props.apiTypesMap.LoginChallengeCheck,
      "value": this.state.login_challenge,
    });
    this.props.backendRequest(request);
  }

  sendLoginCredentials() {
    const request = JSON.stringify({
      "type": this.props.apiTypesMap.LoginRequest,
      "value": {
        "username": this.state.username,
        "password": this.state.password,
        "login_challenge": this.state.login_challenge,
      }
    });
    this.props.backendRequest(request);
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
