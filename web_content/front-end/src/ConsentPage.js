import React from 'react';
import {BACKEND_HOSTNAME,
        CONSENT_CHALLENGE_ENDPOINT} from './Constants.js'

class ConsentPage extends React.Component {
  constructor(props) {
    super(props);

    let urlParams = new URLSearchParams(window.location.search);

    console.log(urlParams.get('consent_challenge'));

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
        body: JSON.stringify({"consent_challenge": this.state.consent_challenge}),
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

  render() {
    return '';
  }
}

export default ConsentPage;
