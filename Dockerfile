FROM golang

WORKDIR /go/src/999k_engine
COPY . .

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

EXPOSE 3000

CMD ["go-wrapper", "run"] # ["999k_engine"]