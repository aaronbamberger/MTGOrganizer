const BACKEND_HOSTNAME = "192.168.50.185"

// We dynamically request the api name to type mappings when we connect to the
// backend, but we have to have a statically determined type for the request
// and response messages for the rest of the types, which is what these are
const API_TYPES_REQUEST = 0
const API_TYPES_RESPONSE = 0

export {BACKEND_HOSTNAME, API_TYPES_REQUEST, API_TYPES_RESPONSE}
