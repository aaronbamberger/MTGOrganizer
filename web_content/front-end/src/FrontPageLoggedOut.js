import React from 'react';
import Oidc from 'oidc-client';

class FrontPageLoggedOut extends React.Component {
  constructor(props) {
    super(props);

    this.authConfig = {
      authority: 'http://192.168.50.185/',
      client_id: 'ArcaneBinders',
      redirect_uri: 'http://192.168.50.185:3000/auth_callback',
    }

    this.handleLogin = this.handleLogin.bind(this);
    this.handleUserLoaded = this.handleUserLoaded.bind(this)

    this.userManager = new Oidc.UserManager(this.authConfig);
  }

  handleUserLoaded(user) {

  }

  handleLogin(event) {
    this.userManager.signinRedirect({state: 'test'});
    event.preventDefault();
  }

  render() {
    return (
      <div id="request_login">
        <form onSubmit={this.handleLogin}>
          <input type="submit" value="Login" />
        </form>
      </div>
    );
  }
}

export default FrontPageLoggedOut;
