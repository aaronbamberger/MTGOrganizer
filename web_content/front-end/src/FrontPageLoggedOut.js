import React from 'react';

class FrontPageLoggedOut extends React.Component {
  constructor(props) {
    super(props);

    this.handleLogin = this.handleLogin.bind(this);
  }

  handleLogin(event) {
    let signinParams = {
      response_mode: 'query',
    }
    this.props.userManager.signinRedirect(signinParams).then(() => {
      console.log(this.props.userManager.getUser());
    });
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
