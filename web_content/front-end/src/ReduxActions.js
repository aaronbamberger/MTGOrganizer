export const REQUEST_CARD_SEARCH_RESULTS = 'REQUEST_CARD_SEARCH_RESULTS'
export const RECEIVE_CARD_SEARCH_RESULTS = 'RECEIVE_CARD_SEARCH_RESULTS'
export const CANCEL_CARD_SEARCH_REQUEST = 'CANCEL_CARD_SEARCH_REQUEST'
export const REQUEST_CARD_DETAIL = 'REQUEST_CARD_DETAIL'
export const RECEIVE_CARD_DETAIL = 'RECEIVE_CARD_DETAIL'
export const CANCEL_CARD_DETAIL_REQUEST = 'CANCEL_CARD_DETAIL_REQUEST'
export const UPDATE_BACKEND_CONNECTION_STATE = 'UPDATE_BACKEND_CONNECTION_STATE'
export const UPDATE_API_TYPES_RECEIVED = 'UPDATE_API_TYPES_RECEIVED'
export const UPDATE_AUTH_COMPLETED = 'UPDATE_AUTH_COMPLETED'

export function requestCardSearchResults(searchTerm) {
  return {
    type: REQUEST_CARD_SEARCH_RESULTS,
    searchTerm: searchTerm,
  }
}

export function receiveCardSearchResults(cards) {
  return {
    type: RECEIVE_CARD_SEARCH_RESULTS,
    cards: cards,
  }
}

export function cancelCardSearchRequest() {
  return {
    type: CANCEL_CARD_SEARCH_REQUEST,
  }
}

export function requestCardDetail(uuid) {
  return {
    type: REQUEST_CARD_DETAIL,
    uuid: uuid,
  }
}

export function receiveCardDetail(cardDetail) {
  return {
    type: RECEIVE_CARD_DETAIL,
    cardDetail: cardDetail,
  }
}

export function cancelCardDetailRequest() {
  return {
    type: CANCEL_CARD_DETAIL_REQUEST,
  }
}

export function updateBackendConnectionState(connected) {
  return {
    type: UPDATE_BACKEND_CONNECTION_STATE,
    connected: connected,
  }
}

export function updateApiTypesReceived(received) {
  return {
    type: UPDATE_API_TYPES_RECEIVED,
    received: received,
  }
}

export function updateAuthCompleted(completed) {
  return {
    type: UPDATE_AUTH_COMPLETED,
    completed: completed,
  }
}
