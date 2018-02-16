FROM alpine:3.7
WORKDIR /app

COPY ./backend/out/elevator-simulator-linux-amd64 ./elevator-simulator

EXPOSE 15657
EXPOSE 3001
ENTRYPOINT [ "./elevator-simulator" ]