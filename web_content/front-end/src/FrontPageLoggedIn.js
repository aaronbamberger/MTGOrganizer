import React from 'react';
import {
  Switch,
  Route
} from 'react-router-dom';

import {BACKEND_HOSTNAME, API_TYPES_REQUEST, API_TYPES_RESPONSE} from './Constants.js';
import CardDetail from './CardDetail.js';
import {CardSearch} from './CardSearch.js';

class FrontPageLoggedIn extends React.Component {
  constructor(props) {
    super(props);

    this.cardSearch = React.createRef();
    this.cardDetail = React.createRef();

    this.backendRequest = this.backendRequest.bind(this);

    this.state = {
      socketConnected: false,
      apiTypesReceived: false,
      apiTypesMap: {}
    };

    this.backendSocket = new WebSocket('ws://' + BACKEND_HOSTNAME + '/backend/api');
    this.backendSocket.addEventListener('open', this.handleWebsocketOpen.bind(this));
    this.backendSocket.addEventListener('message', this.handleWebsocketMessage.bind(this));
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
    }
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

  render() {
    if (!(this.state.socketConnected && this.state.apiTypesReceived)) {
      return (
        <div>
          <h2>Loading UI...</h2>
        </div>
      )
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

export default FrontPageLoggedIn;
