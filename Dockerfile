FROM golang:1.13.5-buster as base
RUN apt-get update && apt-get install -y dumb-init
WORKDIR /go/src/github.com/justinbrumley/megahook.in
COPY . .
RUN go get && go build
ENTRYPOINT ["/usr/bin/dumb-init", "--"]

FROM base as development
ENV ENV development
# Install modd. Used for rebuilding and restarting on code changes
RUN wget https://github.com/cortesi/modd/releases/download/v0.8/modd-0.8-linux64.tgz \
  && tar -C /usr/bin/ -xzf modd-0.8-linux64.tgz --strip-components 1 \
  && chmod +x /usr/bin/modd
CMD ["modd"]

FROM base as production
ENV ENV production
CMD ["./megahook.in"]

