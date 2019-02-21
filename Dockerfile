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

RUN ls -al /go/bin/

# final stage
FROM ubuntu

COPY --from=builder /usr/lib/pkgconfig /usr/lib/pkgconfig
COPY --from=builder /usr/lib/librdkafka* /usr/lib/
COPY --from=builder /app/main main

CMD ["./main"]