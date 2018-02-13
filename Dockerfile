# build image
#FROM golang:latest
FROM gcr.io/mdimages/go-indoormap:v1 AS build

# prepare and copy content
ARG pkg=github.com/ColorfulBridge/IndoorMapTileServer
COPY . $GOPATH/src/$pkg

# get dependencies and install
WORKDIR $GOPATH/src/$pkg
RUN CGO_ENABLED=0 go install

###################### RUN #############################
#run image
FROM golang:alpine
#FROM gcr.io/mdimages/go-indoormap:v1

COPY --from=build $GOPATH/bin/IndoorMapTileServer /go/bin/IndoorMapTileServer

COPY .deploy/key.json /root/service_key.json
ENV GOOGLE_APPLICATION_CREDENTIALS /root/service_key.json
ENV GCLOUD_STORAGE_BUCKET colorful-bridge-mapcontent

WORKDIR /go/bin

CMD /go/bin/IndoorMapTileServer