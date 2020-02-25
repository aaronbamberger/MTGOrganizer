import React from 'react';
import {
  BrowserRouter as Router,
  Switch,
  Route
} from "react-router-dom";
import logo from './logo.svg';
import './App.css';
import CardDetail from './CardDetail.js';
import {CardSearch} from './CardSearch.js';
import {BACKEND_HOSTNAME, RESPONSE_TYPES} from './Constants.js';

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
    console.log("Websocket open: " + event);
  }

  handleWebsocketMessage(event) {
    const response = JSON.parse(event.data);
    if (response.type === RESPONSE_TYPES.CARD_SEARCH_RESPONSE) {
      this.cardSearch.current.receiveNewCards(response.value);
    } else if (response.type === RESPONSE_TYPES.CARD_DETAIL_RESPONSE) {
      this.cardDetail.current.receiveCardDetail(response.value);
    }
  }

  render() {
    return (
      <Switch>
        <Route path="/card/:uuid">
          <CardDetail wrappedComponentRef={this.cardDetail} backendRequest={this.backendRequest} />
        </Route>
        <Route path="/">
          <CardSearch ref={this.cardSearch} backendRequest={this.backendRequest} />
        </Route>
      </Switch>
    );
  }
}

export default App;
