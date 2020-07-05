import React from 'react';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';

import { setLoggedInUser } from './ReduxActions.js';

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
      this.props.history.replace("/auth_callback");
    } else {
      this.props.userManager.signinRedirectCallback().then((user) => {
        console.log("Completed auth callback");
        console.log(user);
        this.props.dispatch(setLoggedInUser(user));
        this.props.history.replace("/");
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

export default connect()(withRouter(AuthCallbackHandler));
