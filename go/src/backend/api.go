package backend

type RequestType int
type ResponseType int

const (
    ApiTypesRequest RequestType = iota
    CardSearchRequest
    CardDetailRequest
)

const (
    ErrorResponse ResponseType = iota
    CardSearchResponse
    CardDetailResponse
)

type RequestMessage struct {
    Type RequestType `json:"type"`
    Value string `json:"value"`
}

type ResponseMessage struct {
    Type ResponseType `json:"type"`
    Value interface{} `json:"value"`
}
