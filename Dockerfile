FROM golang:alpine

ARG pkg=github.com/ColorfulBridge/IndoorMapTileServer
RUN apk add --no-cache ca-certificates

COPY . $GOPATH/src/$pkg
RUN set -ex \
      && apk add --no-cache --virtual .build-deps \
              git \
      && go get -v $pkg/... \
      && apk del .build-deps
RUN go install $pkg 

COPY .keys/colorful-bridge_servicekey.json /root/service_key.json
ENV GOOGLE_APPLICATION_CREDENTIALS /root/service_key.json
ENV GCLOUD_STORAGE_BUCKET colorful-bridge-mapcontent

WORKDIR $GOPATH/src/$pkg

CMD IndoorMapTileServer