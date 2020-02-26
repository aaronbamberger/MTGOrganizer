package backend

type RequestType int
type ResponseType int
//go:generate stringer -type=RequestType
//go:generate stringer -type=ResponseType

const (
    ApiTypesRequest RequestType = iota
    CardSearchRequest
    CardDetailRequest
)

const (
    ApiTypesResponse ResponseType = iota
    ErrorResponse
    CardSearchResponse
    CardDetailResponse
)

var requestTypes  = [...]RequestType{
    ApiTypesRequest,
    CardSearchRequest,
    CardDetailRequest}

var responseTypes = [...]ResponseType{
    ApiTypesResponse,
    ErrorResponse,
    CardSearchResponse,
    CardDetailResponse}

type RequestMessage struct {
    Type RequestType `json:"type"`
    Value string `json:"value"`
}

type ResponseMessage struct {
    Type ResponseType `json:"type"`
    Value interface{} `json:"value"`
}

type ApiTypes struct {
    RequestTypes map[string]RequestType `json:"request_types"`
    ResponseTypes map[string]ResponseType `json:"response_types"`
}

func apiTypes(done <-chan interface{}, respChan chan<- ResponseMessage) {
    types := make(map[string]int)

    for _, requestType := range requestTypes {
        types[ requestType.String() ] = int(requestType)
    }
    for _, responseType := range responseTypes {
        types[ responseType.String() ] = int(responseType)
    }

    select {
    case <-done:
    case respChan <- ResponseMessage{Type: ApiTypesResponse, Value: types}:
    }
}
