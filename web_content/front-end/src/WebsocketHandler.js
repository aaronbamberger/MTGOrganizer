import React from 'react';
import { connect } from 'react-redux';

import {
  BACKEND_HOSTNAME,
  API_TYPES_REQUEST,
  API_TYPES_RESPONSE
} from './Constants.js';
import {
  receiveCardSearchResults,
  cancelCardSearchRequest,
  receiveCardDetail,
  cancelCardDetailRequest,
  updateBackendConnectionState,
  updateApiTypesReceived,
  setAuthRequestPending,
  updateAuthCompleted,
} from './ReduxActions.js'
/*
const backendDefaultState = {
  connected: false,
  apiTypesReceived: false,
  authCompleted: false,
  authRequestPending: false,
  ready: false,
}
*/
const mapStateToProps = (state) => {
  return {
    backendConnected: state.backend.connected,
    backendApiTypesReceived: state.backend.apiTypesReceived,
    backendAuthRequestPending: state.backend.authRequestPending,
    user: state.user.user,
    cardSearchRequest: state.cardSearch.searchRequested,
    cardSearchTerm: state.cardSearch.searchTerm,
    cardDetailRequest: state.cardDetail.searchRequested,
    cardDetailUUID: state.cardDetail.uuid,
  };
}

class WebsocketHandler extends React.Component {
  constructor(props) {
    super(props);

    this.backendRequest = this.backendRequest.bind(this);
    this.handleWebsocketOpen = this.handleWebsocketOpen.bind(this);
    this.handleWebsocketMessage = this.handleWebsocketMessage.bind(this);

    this.backendSocket = new WebSocket('ws://' + BACKEND_HOSTNAME + '/backend/api');
    this.backendSocket.addEventListener('open', this.handleWebsocketOpen);
    this.backendSocket.addEventListener('message', this.handleWebsocketMessage);

    this.state = {
      socketConnected: false,
      apiTypesReceived: false,
      apiTypesMap: {},
    };
  }

  componentDidUpdate() {
    // Authenticate to the backend if we're connected and have received the
    // api types map
    if (this.props.backendConnected && this.props.backendApiTypesReceived &&
        !this.props.backendAuthRequestPending) {
      const request = JSON.stringify({
        "type": this.state.apiTypesMap.AuthUserRequest,
        "value": {
          "token": this.props.user.access_token,
          "subject": this.props.user.profile.sub,
        },
      });
      this.backendRequest(request);
    }

    // Handle a new card search request
    if (this.props.cardSearchRequest) {
      // Construct and send the backend request
      const request = JSON.stringify({
        "type": this.state.apiTypesMap.CardSearchRequest,
        "value": this.props.cardSearchTerm,
      });
      this.backendRequest(request);

      // Reset the search request flag now that we've sent the search request
      this.props.dispatch(cancelCardSearchRequest())
    }

    if (this.props.cardDetailRequest) {
      const request = JSON.stringify({
        "type": this.state.apiTypesMap.CardDetailRequest,
        "value": this.props.cardDetailUUID,
      });
      this.backendRequest(request);
      
      // Reset the search request flag now that we've sent the search request
      this.props.dispatch(cancelCardDetailRequest())
    }
  }

  handleWebsocketOpen(event) {
    // Request the list of api name to type mappings from the backend
    this.backendRequest(JSON.stringify({"type": API_TYPES_REQUEST, "value": ""}));
    this.props.dispatch(updateBackendConnectionState(true));
    console.log("Websocket open: " + event);
  }

  handleWebsocketMessage(event) {
    const response = JSON.parse(event.data);
    console.log('Response type: ' + response.type);
    switch (response.type) {
      case API_TYPES_RESPONSE:
        this.setState({apiTypesMap: response.value});
        this.props.dispatch(updateApiTypesReceived(true));
        break;
      case this.state.apiTypesMap.CardSearchResponse:
        this.props.dispatch(receiveCardSearchResults(response.value));
        break;
      case this.state.apiTypesMap.CardDetailResponse:
        this.props.dispatch(receiveCardDetail(response.value));
        break;
      case this.state.apiTypesMap.AuthUserResponse:
        if (response.value.auth_successful) {
          this.props.dispatch(updateAuthCompleted(true));
        }
      default:
        break;
    }
  }

  backendRequest(request) {
    // Only send a request if the socket is in the "OPEN" state
    if (this.backendSocket.readyState === WebSocket.OPEN) {
      this.backendSocket.send(request);
    } else {
      console.log("Can't send request " + request + " to websocket, in state " +
        this.backendSocket.readyState);
    }
  }

  render() {
    // This component has no UI associated with it
    return null;
  }
}

export default connect(mapStateToProps)(WebsocketHandler);
