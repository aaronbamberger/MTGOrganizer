import React from 'react';
import {
  BrowserRouter as Router,
  Switch,
  Route
} from 'react-router-dom';
import { createStore, combineReducers } from 'redux';
import { Provider, connect } from 'react-redux';
import Oidc from 'oidc-client';

import './App.css';
import LoginPage from './LoginPage.js';
import ConsentPage from './ConsentPage.js';
import LogoutPage from './LogoutPage.js';
import FrontPageLoggedOut from './FrontPageLoggedOut.js';
import FrontPageLoggedIn from './FrontPageLoggedIn.js';
import AuthCallbackHandler from './AuthCallbackHandler.js';
import {
  cardSearchReducer,
  cardDetailReducer,
  backendStateReducer,
  userReducer,
} from './ReduxReducers.js';
import {
  setLoggedInUser,
} from './ReduxActions.js';

function App() {
  const rootReducer = combineReducers({
    cardSearch: cardSearchReducer,
    cardDetail: cardDetailReducer,
    backend: backendStateReducer,
    user: userReducer,
  });
  const reduxStore = createStore(rootReducer);

  return (
    <div className="App">
      <Provider store={reduxStore}>
        <Router>
          <MTGOrganizerWidget />
        </Router>
      </Provider>
    </div>
  );
}

const mapStateToProps = (state) => {
  return {
    user: state.user.user,
  };
}

class MTGOrganizer extends React.Component {
  constructor(props) {
    super(props);

    this.authConfig = {
      authority: 'http://192.168.50.185/',
      client_id: 'ArcaneBinders',
      response_type: 'code',
      redirect_uri: 'http://192.168.50.185/auth_callback/',
      response_mode: 'query',
    }

    this.userManager = new Oidc.UserManager(this.authConfig);
  }

  componentDidMount() {
    this.userManager.getUser().then((user) => {
      if(user) {
        console.log("Logged in user ", user);
        this.props.dispatch(setLoggedInUser(user));
      } else {
        console.log("No user logged in");
      }
    });
  }

  render() {
    if (!this.props.user) {
      return (
        <Switch>
          <Route path="/auth/login">
            <LoginPage />
          </Route>
          <Route path="/auth/consent">
            <ConsentPage />
          </Route>
          <Route path="/auth/logout">
            <LogoutPage />
          </Route>
          <Route path="/auth_callback">
            <AuthCallbackHandler userManager={this.userManager} />
          </Route>
          <Route path="/">
            <FrontPageLoggedOut userManager={this.userManager} />
          </Route>
        </Switch>
      );
    } else {
      return (
        <FrontPageLoggedIn userManager={this.userManager} />
      );
    }
  }
}

const MTGOrganizerWidget = connect(mapStateToProps)(MTGOrganizer);

export default App;
