FROM golang:1.19.3 as build

WORKDIR /app
COPY ./server/ ./server/
COPY makefile ./

RUN go work init ./server
RUN make build-server
RUN useradd -u 10001 scratchuser

## Build final image
FROM scratch
WORKDIR /
COPY --from=build /app/bin/websocketserver /app/websocketserver
COPY --from=build /etc/passwd /etc/passwd
#RUN chown scratchuser /app/websocketserver
EXPOSE 8080
USER scratchuser
ENTRYPOINT ["/app/websocketserver"]