FROM golang:latest
ENV GOPATH /go
WORKDIR /go/src
COPY src/card-importer card-importer/
COPY src/carddb carddb/
COPY src/mtgcards mtgcards/
# Get the dependencies
RUN go get -d -v "github.com/go-sql-driver/mysql"
RUN go get -d -v "github.com/influxdata/influxdb1-client/v2"
RUN go install card-importer

VOLUME ["/var/card-importer/"]

ENTRYPOINT ["card-importer"]
