package backend

type RequestType int
type ResponseType int

const (
    CardSearchRequest RequestType = iota
)

const (
    ErrorResponse ResponseType = iota
    CardSearchResponse
)

type RequestMessage struct {
    Type RequestType `json:"type"`
    Value string `json:"value"`
}

type ResponseMessage struct {
    Type ResponseType `json:"type"`
    Value interface{} `json:"value"`
}
