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

const mapStateToProps = (state) => {
  return {
    backendReady: state.backendState.ready,
  };
}

class FrontPageLoggedIn extends React.Component {
  render() {
    const loadingDisplay = this.props.backendReady ? 'none' : 'initial';
    const uiDisplay = this.props.backendReady ? 'initial' : 'none';

    return (
      <div>
        <WebsocketHandler />
        <div id="loading_message" style={{display: loadingDisplay}}>
          <h2>Loading UI...</h2>
        </div>
        <div id="main_ui" style={{display: uiDisplay}}>
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
      </div>
    );
  }
}

export default connect(mapStateToProps)(FrontPageLoggedIn);
