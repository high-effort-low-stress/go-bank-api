FROM alpine:latest
EXPOSE 8080
WORKDIR /app

COPY app .

CMD ["./app"]