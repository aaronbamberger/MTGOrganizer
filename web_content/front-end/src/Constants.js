const BACKEND_HOSTNAME = "192.168.50.185"

const REQUEST_TYPES = Object.freeze({
  "API_TYPES_REQUEST": 0,
  "CARD_SEARCH_REQUEST": 1,
  "CARD_DETAIL_REQUEST": 2,
});

const RESPONSE_TYPES = Object.freeze({
  "ERROR_RESPONSE": 0,
  "CARD_SEARCH_RESPONSE": 1,
  "CARD_DETAIL_RESPONSE": 2,
});

export {BACKEND_HOSTNAME, REQUEST_TYPES, RESPONSE_TYPES}
