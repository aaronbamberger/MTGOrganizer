import React from 'react';
import {
  Switch,
  Route
} from 'react-router-dom';
import { connect } from 'react-redux';

import AccountInfoWidget from './AccountInfoWidget.js';
import CardDetail from './CardDetail.js';
import CardSearchWidget from './CardSearch.js';
import WebsocketHandler from './WebsocketHandler.js';

import { BACKEND_HOSTNAME } from './Constants.js';

const mapStateToProps = (state) => {
  return {
    backendConnected: state.backendState.connected,
  };
}

class FrontPageLoggedIn extends React.Component {
  constructor(props) {
    super(props);

    this.backendSocket = new WebSocket('ws://' + BACKEND_HOSTNAME + '/backend/api');
  }

  render() {
    return (
      <div>
        <WebsocketHandler backendSocket={this.backendSocket} />
        <AccountInfoWidget userManager={this.props.userManager} />
        <Switch>
          <Route path="/card/:uuid">
            <CardDetail />
          </Route>
          <Route path="/">
            <CardSearchWidget />
          </Route>
        </Switch>
      </div>
    );
  }
}

export default connect(mapStateToProps)(FrontPageLoggedIn);
