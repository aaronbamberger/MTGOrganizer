import React from 'react';

class AuthCallbackHandler extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      loginSuccessful: false,
      errorType: null,
      errorMessage: null,
    }
  }

  componentDidMount() {
    const callbackResults = new URLSearchParams(window.location.search);
    if (callbackResults.has("error")) {
      this.setState({
        errorType: callbackResults.get("error"),
        errorMessage: callbackResults.get("error_description"),
      });
    } else {
      this.props.userManager.signinRedirectCallback().then((user) => {
        console.log(user);
      });
      this.setState({
        loginSuccessful: true,
      });
    }
  }

  render() {
    if (this.state.loginSuccessful) {
      return (
        <div>
          Auth callback page
        </div>
      );
    } else {
      return (
        <div>
          Error: {this.state.errorType}
          <br />
          {this.state.errorMessage}
        </div>
      );
    }
  }
}

export default AuthCallbackHandler;
