import { REQUEST_CARD_SEARCH_RESULTS,
         RECEIVE_CARD_SEARCH_RESULTS,
         CANCEL_CARD_SEARCH_REQUEST,
         REQUEST_CARD_DETAIL,
         RECEIVE_CARD_DETAIL,
         CANCEL_CARD_DETAIL_REQUEST,
         UPDATE_BACKEND_CONNECTION_STATE,
         UPDATE_API_TYPES_RECEIVED,
         UPDATE_AUTH_COMPLETED,
} from './ReduxActions.js';

const cardSearchDefaultState = {
  searchRequested: false,
  searchTerm: '',
  searchResults: [],
};

export function cardSearchReducer(state = cardSearchDefaultState, action) {
  switch (action.type) {
    case REQUEST_CARD_SEARCH_RESULTS:
      return Object.assign({}, state, {
        searchTerm: action.searchTerm,
        searchRequested: true,
      });
    case RECEIVE_CARD_SEARCH_RESULTS:
      return Object.assign({}, state, {
        searchResults: action.cards,
      });
    case CANCEL_CARD_SEARCH_REQUEST:
      return Object.assign({}, state, {
        searchRequested: false,
      });
    default:
      return state;
  }
}

const cardDetailDefaultState = {
  searchRequested: false,
  uuid: '',
  cardDetail: null,
}

export function cardDetailReducer(state = cardDetailDefaultState, action) {
  switch (action.type) {
    case REQUEST_CARD_DETAIL:
      return Object.assign({}, state, {
        uuid: action.uuid,
        searchRequested: true,
      });
    case RECEIVE_CARD_DETAIL:
      return Object.assign({}, state, {
        cardDetail: action.cardDetail,
      });
    case CANCEL_CARD_DETAIL_REQUEST:
      return Object.assign({}, state, {
        searchRequested: false,
      });
    default:
      return state;
  }
}

const backendDefaultState = {
  connected: false,
  apiTypesReceived: false,
  authCompleted: true,
  ready: false,
}

export function backendStateReducer(state = backendDefaultState, action) {
  switch (action.type) {
    case UPDATE_BACKEND_CONNECTION_STATE:
    {
      const ready = action.connected && state.apiTypesReceived && state.authCompleted;
      return Object.assign({}, state, {
        connected: action.connected,
        ready: ready,
      });
    }
    case UPDATE_API_TYPES_RECEIVED:
    {
      const ready = action.received && state.connected && state.authCompleted;
      return Object.assign({}, state, {
        apiTypesReceived: action.received,
        ready: ready,
      });
    }
    case UPDATE_AUTH_COMPLETED:
    {
      const ready = action.completed && state.apiTypesReceived && state.connected;
      return Object.assign({}, state, {
        authCompleted: action.completed,
        ready: ready,
      });
    }
    default:
      return state;
  }
}
