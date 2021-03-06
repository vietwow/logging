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

CMD ["./logging", "--sumologic.url=https://endpoint3.collection.us2.sumologic.com/receiver/v1/http/ZaVnC4dhaV3roU3-kIZNtE2iTbUBN9crYJBdOVmhT0KD7_tmOjxO0dRuSn0DmGYcSOznK5nRyCqsaPm69tsW2CVr51cK9ook29GUHCJg8bd1oa1GpZzY7A==", "--sumologic.source.category=staging/logging"]