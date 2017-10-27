FROM golang:alpine

WORKDIR /go/src/999k_engine
COPY . .

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

RUN apk del --no-cache bash git openssh

EXPOSE 3000

CMD ["go-wrapper", "run"] # ["999k_engine"]