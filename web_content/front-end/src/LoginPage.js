import React from 'react';
import Oidc from 'oidc-client';
import {withRouter} from 'react-router-dom';

class LoginPage extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      login_needed: false,
      username: '',
      password: '',
      user: null,
    }

    this.authConfig = {
      authority: 'http://192.168.50.185/',
      client_id: 'ArcaneBinders',
      response_type: 'id_token token',
      redirect_uri: 'http://192.168.50.185:3000/',
    }

    this.receiveLoginChallengeResult = this.receiveLoginChallengeResult.bind(this);
    this.handleUsername = this.handleUsername.bind(this);
    this.handlePassword = this.handlePassword.bind(this);
    this.handleLogin = this.handleLogin.bind(this);

    let urlParams = new URLSearchParams(window.location.search);

    this.login_challenge = urlParams.get('login_challenge');

    this.userManager = new Oidc.UserManager(this.authConfig);
  }

  componentDidMount() {
    this.sendLoginChallengeRequest()
  }

  sendLoginChallengeRequest() {
    const request = JSON.stringify({
      "type": this.props.apiTypesMap.LoginChallengeCheck,
      "value": this.login_challenge,
    });
    this.props.backendRequest(request);
  }

  receiveLoginChallengeResult(result) {
    console.log('Login challenge result: ' + result);
    if (result['skip'] === false) {
      this.setState({login_needed: true});
    }
  }

  handleUsername(event) {
    this.setState({username: event.target.value});
  }

  handlePassword(event) {
    this.setState({password: event.target.value});
  }

  handleLogin(event) {
    console.log('Username: ' + this.state.username);
    console.log('Password: ' + this.state.password);
    this.userManager.signinRedirect({state: 'test'});
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

export default withRouter(LoginPage);
