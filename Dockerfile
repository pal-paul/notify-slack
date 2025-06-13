FROM golang:1.24-bullseye as build

RUN mkdir -p /notify-slack/
COPY . /notify-slack/
WORKDIR /notify-slack

ENV GO111MODULE=on
RUN make install
RUN make build

# Now copy it into our base image.
FROM gcr.io/distroless/base
COPY --from=build /notify-slack/cmd/cmd /cmd

ENTRYPOINT ["/cmd"]