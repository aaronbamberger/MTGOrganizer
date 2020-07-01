import React from 'react';
import {
  BrowserRouter as Router,
  Switch,
  Route
} from 'react-router-dom';
import Oidc from 'oidc-client';
import './App.css';
import LoginPage from './LoginPage.js';
import ConsentPage from './ConsentPage.js';
import LogoutPage from './LogoutPage.js';
import FrontPageLoggedOut from './FrontPageLoggedOut.js';
import FrontPageLoggedIn from './FrontPageLoggedIn.js';
import AuthCallbackHandler from './AuthCallbackHandler.js';

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

    this.authConfig = {
      authority: 'http://192.168.50.185/',
      client_id: 'ArcaneBinders',
      response_type: 'code',
      redirect_uri: 'http://192.168.50.185/auth_callback/',
      response_mode: 'query',
    }

    this.userManager = new Oidc.UserManager(this.authConfig);

    this.state = {
      loggedIn: false,
      user: null,
    };
  }

  componentDidMount() {
    this.userManager.getUser().then((user) => {
      if(user) {
        console.log("Logged in user ", user);
        this.setState({loggedIn: true, user: user});
      } else {
        console.log("No user logged in");
      }
    });
  }

  render() {
    if (!this.state.loggedIn) {
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
            <AuthCallbackHandler
              userManager={this.userManager} />
          </Route>
          <Route path="/">
            <FrontPageLoggedOut
              userManager={this.userManager} />
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

export default App;
