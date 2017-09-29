FROM golang:latest

RUN mkdir -p /go/src/app
ADD ./github.com /go/src/github.com
ADD ./app /go/src/app
WORKDIR /go/src/app

ENV BXGO_PORT 8000

#RUN go get
RUN go install

EXPOSE $BXGO_PORT

CMD ["/go/bin/app"]

#docker build -t bxgo-trade /var/www/bxgo.trade/ --no-cache=true
#docker run --name bxgo -p 8000:8000 -v /var/www/bxgo.trade:/go/src/app -e BXGO_PORT="8000" --restart=always bxgo
