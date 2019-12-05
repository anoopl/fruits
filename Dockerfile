# build go app
FROM golang:alpine AS build-env
RUN apk --no-cache add build-base git bzr mercurial gcc
ADD . /src
RUN cd /src && go build -o items

#Minimal image for using
FROM alpine
WORKDIR /items
COPY --from=build-env /src/items /items/
CMD [ "/items/items" ]