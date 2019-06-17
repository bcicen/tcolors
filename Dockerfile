FROM quay.io/vektorcloud/go:1.12

RUN apk add --no-cache make

WORKDIR /app
COPY go.mod .
RUN go mod download

COPY . .
RUN make build

FROM scratch
ENV TERM=xterm-256color
ENV COLORTERM=truecolor
COPY --from=0 /app/tcolors /tcolors
ENTRYPOINT ["/tcolors", "-output-on-exit"]
