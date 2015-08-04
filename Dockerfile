FROM golang:1.4.2

RUN go get github.com/codegangsta/negroni
RUN go get github.com/yosssi/ace

ADD . /go/src/git.xeserv.us/ponychat/shoutpage

RUN go get git.xeserv.us/ponychat/shoutpage

EXPOSE 3000

RUN useradd --create-home shout
USER shout

CMD shoutpage
