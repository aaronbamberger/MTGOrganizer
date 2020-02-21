import React from 'react';
import logo from './logo.svg';
import './App.css';
import './Keyrune-master/css/keyrune.css'

function App() {
  return (
    <div className="App">
      <WebsocketHelloWorld />
    </div>
  );
}

class WebsocketHelloWorld extends React.Component {
  constructor(props) {
    super(props);
    
    this.state = {value: '', cards: []};

    this.socket = new WebSocket('ws://192.168.50.185:8085/api');
    this.socket.addEventListener('open', this.handleWebsocketOpen.bind(this));
    this.socket.addEventListener('message', this.handleWebsocketMessage.bind(this));

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(event) {
    this.setState({value: event.target.value});
  }

  handleSubmit(event) {
    console.log('Submitted: ' + this.state.value);
    this.socket.send(JSON.stringify({"type": 0, "value": this.state.value}));
    event.preventDefault();
  }

  handleWebsocketOpen(event) {
    console.log("Websocket open: " + event);
  }

  handleWebsocketMessage(event) {
    console.log("Received response: " + event.data);
    this.setState({"cards": JSON.parse(event.data)});
  }

  render() {
    const cardTableRows = this.state.cards.map((cardInfo) =>
      <tr>
          <td>{cardInfo.name}</td>
          <td><span className={"ss ss-" + cardInfo.set_keyrune_code.toLowerCase()}></span></td>
        </tr>
      );

    return (
      <div id="card_search">
        <form onSubmit={this.handleSubmit}>
          <label>
            Test:
            <input type="text" value={this.state.value} onChange={this.handleChange} />
         </label>
          <input type="submit" value="Submit" />
        </form>
        <div id="card_results">
          <table>
            <thead>
              <tr>
                <th>Card Name</th>
                <th>Set</th>
              </tr>
            </thead>
            <tbody>
              {cardTableRows}
            </tbody>
          </table>
        </div>
      </div>
    );
  }
}

export default App;
