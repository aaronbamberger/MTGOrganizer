import React from 'react';

class ConsentPage extends React.Component {
  constructor(props) {
    super(props);

    let urlParams = new URLSearchParams(window.location.search);

    console.log(urlParams.get('consent_challenge'));

    this.state = {
      consent_challenge: urlParams.get('consent_challenge'),
    }

    this.receiveConsentResponse = this.receiveConsentResponse.bind(this);
  }

  componentDidMount() {
    this.sendConsentChallengeRequest()
  }

  sendConsentChallengeRequest() {
    const request = JSON.stringify({
      "type": this.props.apiTypesMap.ConsentChallengeCheck,
      "value": this.state.consent_challenge,
    });
    this.props.backendRequest(request);
  }

  receiveConsentResponse(response) {
    console.log('Consent response: ' + response);
    let responseJSON = JSON.parse(response);
    window.location.replace(responseJSON['redirect_to']);
  }

  render() {
    return '';
  }
}

export default ConsentPage;
