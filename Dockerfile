## BUILD
FROM golang:1.19-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /atlas-game-lookout

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /atlas-game-lookout /atlas-game-lookout
COPY tiles /tiles

EXPOSE 3000

ENTRYPOINT ["/atlas-game-lookout"]