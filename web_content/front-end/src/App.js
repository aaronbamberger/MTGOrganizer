import React from 'react';
import {
    BrowserRouter as Router,
    Switch,
    Route,
    Link
} from "react-router-dom";
import logo from './logo.svg';
import './App.css';
import {CardSearch} from './CardSearch.js'

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

    this.state = {cardSearchCards: []}

    this.backendSocket = new WebSocket('ws://192.168.50.185:8085/api');
    this.backendSocket.addEventListener('open', this.handleWebsocketOpen.bind(this));
    this.backendSocket.addEventListener('message', this.handleWebsocketMessage.bind(this));

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
    if (response.type == 1) {
      this.setState({cardSearchCards: response.value});
    }
  }

  render() {
    return (
      <CardSearch backendRequest={this.backendRequest} cards={this.state.cardSearchCards} />
    );
  }
}

export default App;
