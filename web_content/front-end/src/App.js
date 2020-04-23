import React from 'react';
import {
  BrowserRouter as Router,
  Switch,
  Route
} from 'react-router-dom';
import './App.css';
import CardDetail from './CardDetail.js';
import LoginPage from './LoginPage.js';
import ConsentPage from './ConsentPage.js';
import FrontPageLoggedOut from './FrontPageLoggedOut.js';
import {CardSearch} from './CardSearch.js';
import {BACKEND_HOSTNAME, API_TYPES_REQUEST, API_TYPES_RESPONSE} from './Constants.js';

function App() {
  return (
    <div className="App">
      <Router>
        <MTGOrganizer />
      </Router>
    </div>
  );
}

class MTGOrganizer extends React.Component {
  constructor(props) {
    super(props);

    this.backendSocket = new WebSocket('ws://' + BACKEND_HOSTNAME + ':8085/api');
    this.backendSocket.addEventListener('open', this.handleWebsocketOpen.bind(this));
    this.backendSocket.addEventListener('message', this.handleWebsocketMessage.bind(this));

    this.state = {
      loggedIn: false,
      socketConnected: false,
      apiTypesReceived: false,
      apiTypesMap: {}
    };

    this.loginPage = React.createRef();
    this.consentPage = React.createRef();
    this.cardSearch = React.createRef();
    this.cardDetail = React.createRef();

    this.backendRequest = this.backendRequest.bind(this);
  }

  backendRequest(request) {
    // Only send a request if the socket is in the "OPEN" state
    if (this.backendSocket.readyState === 1) {
      this.backendSocket.send(request);
    } else {
      console.log("Can't send request " + request + " to websocket, in state " +
        this.backendSocket.readyState);
    }
  }

  handleWebsocketOpen(event) {
    this.setState({socketConnected: true});
    // Request the list of api name to type mappings from the backend
    this.backendRequest(JSON.stringify({"type": API_TYPES_REQUEST, "value": ""}))
    console.log("Websocket open: " + event);
  }

  handleWebsocketMessage(event) {
    const response = JSON.parse(event.data);
    console.log('Response type: ' + response.type);
    if (response.type === API_TYPES_RESPONSE) {
      this.setState({apiTypesMap: response.value, apiTypesReceived: true})
      console.log('Received types map');
    } else if (response.type === this.state.apiTypesMap.CardSearchResponse) {
      this.cardSearch.current.receiveNewCards(response.value);
    } else if (response.type === this.state.apiTypesMap.CardDetailResponse) {
      this.cardDetail.current.receiveCardDetail(response.value);
    } else if (response.type === this.state.apiTypesMap.LoginChallengeResponse) {
      console.log('In main app, received ' + response.value);
      this.loginPage.current.receiveLoginChallengeResult(response.value);
    } else if (response.type === this.state.apiTypesMap.LoginResponse) {
      this.loginPage.current.receiveLoginResponse(response.value);
    } else if (response.type === this.state.apiTypesMap.ConsentResponse) {
      this.consentPage.current.receiveConsentResponse(response.value);
    }
  }

  render() {
    if (!this.state.apiTypesReceived) {
      return (
        <h2>Loading...</h2>
      );
    } else if (!this.state.loggedIn) {
      return (
        <Switch>
          <Route path="/auth/login">
            <LoginPage
              ref={this.loginPage}
              backendRequest={this.backendRequest}
              apiTypesMap={this.state.apiTypesMap} />
          </Route>
          <Route path="/auth/consent">
            <ConsentPage
              ref={this.consentPage}
              backendRequest={this.backendRequest}
              apiTypesMap={this.state.apiTypesMap} />
          </Route>
          <Route path="/">
            <FrontPageLoggedOut />
          </Route>
        </Switch>
      );
    } else {
      return (
        <Switch>
          <Route path="/card/:uuid">
            <CardDetail
              wrappedComponentRef={this.cardDetail}
              backendRequest={this.backendRequest}
              apiTypesMap={this.state.apiTypesMap} />
          </Route>
          <Route path="/">
            <CardSearch
              ref={this.cardSearch}
              backendRequest={this.backendRequest}
              apiTypesMap={this.state.apiTypesMap} />
          </Route>
        </Switch>
      );
    }
  }
}

export default App;
