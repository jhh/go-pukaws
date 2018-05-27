FROM golang:latest 
LABEL Name=pukaws Version=0.0.1 
WORKDIR /go/src/jhhgo.us/pukaws 
COPY . .
RUN go get github.com/kardianos/govendor
RUN govendor sync
# RUN go-wrapper download
RUN go-wrapper install
# RUN go build -o main .
EXPOSE 80
ENV PORT 80
CMD ["go-wrapper", "run"]
