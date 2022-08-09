FROM golang:1.17.8

ENV REPO_URL=github.com/jedzeins/image-board
ENV GOPATH=/app
ENV APP_PATH=$GOPATH/src/$REPO_URL
ENV WORKPATH=$APP_PATH/src

COPY . $WORKPATH

WORKDIR $WORKPATH

RUN GOOS=linux GOARCH=amd64 GO111MODULE=auto go build -o app ./src
ENTRYPOINT [ "./app" ] 

# ENTRYPOINT ["tail", "-f", "/dev/null"]
