# build stage
FROM golang as builder 

ENV GO111MODULE=on

RUN git clone https://github.com/edenhill/librdkafka.git

WORKDIR librdkafka

RUN ./configure --prefix /usr

RUN make

RUN make install

RUN curl https://glide.sh/get | sh

WORKDIR /go/src/logging
ADD . /go/src/logging

RUN glide install && go install logging

# final stage
FROM ubuntu

RUN apt-get update && apt-get install -y telnet curl

COPY --from=builder /usr/lib/pkgconfig /usr/lib/pkgconfig
COPY --from=builder /usr/lib/librdkafka* /usr/lib/
COPY --from=builder /go/bin/logging logging

CMD ["./logging"]