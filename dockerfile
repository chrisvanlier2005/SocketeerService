# syntax=docker/dockerfile:1.2


FROM golang:1.19

# is the WORKDIR absolute or relative to the Dockerfile?
# a: relative to the Dockerfile


WORKDIR /src

# copy the source code into the container
# copy the source code into the container recursively
COPY . .

RUN go mod download

COPY *.go ./



# build the binary
RUN go build -o /JobHiraMicroservice

EXPOSE 8081

CMD [ "/JobHiraMicroservice"  ]
