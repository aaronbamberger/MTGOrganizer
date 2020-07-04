import React from 'react';
import { connect } from 'react-redux';

import {
  API_TYPES_REQUEST,
  API_TYPES_RESPONSE
} from './Constants.js';
import {
  receiveCardSearchResults,
  cancelCardSearchRequest,
  receiveCardDetail,
  cancelCardDetailRequest,
  updateBackendConnectionState,
} from './ReduxActions.js'

const mapStateToProps = (state) => {
  return {
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

    this.state = {
      socketConnected: false,
      apiTypesReceived: false,
      apiTypesMap: {},
    };

    this.props.backendSocket.addEventListener('message', this.handleWebsocketMessage);
  }

  componentDidMount() {
    if (this.props.backendSocket.readyState === WebSocket.OPEN) {
      this.handleWebsocketOpen();
    } else {
      this.props.backendSocket.addEventListener('open', this.handleWebsocketOpen);
    }
  }

  componentDidUpdate() {
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
    this.setState({socketConnected: true});
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
        this.setState({apiTypesMap: response.value, apiTypesReceived: true});
        break;
      case this.state.apiTypesMap.CardSearchResponse:
        this.props.dispatch(receiveCardSearchResults(response.value));
        break;
      case this.state.apiTypesMap.CardDetailResponse:
        this.props.dispatch(receiveCardDetail(response.value));
        break;
    }
  }

  backendRequest(request) {
    // Only send a request if the socket is in the "OPEN" state
    if (this.props.backendSocket.readyState === WebSocket.OPEN) {
      this.props.backendSocket.send(request);
    } else {
      console.log("Can't send request " + request + " to websocket, in state " +
        this.props.backendSocket.readyState);
    }
  }

  render() {
    // This component has no UI associated with it
    return null;
  }
}

export default connect(mapStateToProps)(WebsocketHandler);
