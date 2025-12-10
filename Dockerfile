FROM alpine:latest
EXPOSE 8080
WORKDIR /app

COPY app .

RUN chmod +x ./app

COPY templates/* templates

CMD ["./app"]