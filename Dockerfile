FROM alpine:latest
EXPOSE 8080
WORKDIR /app

COPY app .
RUN mkdir templates
COPY templates/* templates

CMD ["./app"]