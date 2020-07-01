import React from 'react';
import {BACKEND_HOSTNAME,
        LOGOUT_CHALLENGE_ENDPOINT} from './Constants.js'

class LogoutPage extends React.Component {
  constructor(props) {
    super(props);

    let urlParams = new URLSearchParams(window.location.search);

    this.state = {
      logout_challenge: urlParams.get('logout_challenge'),
    }
  }

  componentDidMount() {
    this.sendLogoutChallengeRequest()
  }

  sendLogoutChallengeRequest() {
    const req = new Request("http://" + BACKEND_HOSTNAME + LOGOUT_CHALLENGE_ENDPOINT,
      {
        method: "POST",
        body: JSON.stringify({"challenge": this.state.logout_challenge}),
      }
    );
    fetch(req).then(response => {
      console.log(response);
      if (response.redirected) {
        window.location.href = response.url;
      }
    });
  }

  render() {
    return '';
  }
}

export default LogoutPage;
