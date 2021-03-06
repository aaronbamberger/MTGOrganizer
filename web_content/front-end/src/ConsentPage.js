import React from 'react';
import {BACKEND_HOSTNAME,
        CONSENT_CHALLENGE_ENDPOINT} from './Constants.js'

class ConsentPage extends React.Component {
  constructor(props) {
    super(props);

    let urlParams = new URLSearchParams(window.location.search);

    this.state = {
      consent_challenge: urlParams.get('consent_challenge'),
    }
  }

  componentDidMount() {
    this.sendConsentChallengeRequest()
  }

  sendConsentChallengeRequest() {
    const req = new Request("http://" + BACKEND_HOSTNAME + CONSENT_CHALLENGE_ENDPOINT,
      {
        method: "POST",
        body: JSON.stringify({"challenge": this.state.consent_challenge}),
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

export default ConsentPage;
