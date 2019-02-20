FROM golang:alpine

# Install OS-level dependencies.
RUN apk update &&     apk add curl git &&     curl https://glide.sh/get | sh

# Copy our source code into the container.
WORKDIR /go/src/logging
ADD . /go/src/logging

# Install our golang dependencies and compile our binary.
RUN glide install
RUN go install logging

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=0 /go/bin/logging .

CMD ["./logging"]
