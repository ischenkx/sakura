FROM golang:latest
COPY . /app
WORKDIR /app
RUN go build . -o main
EXPOSE 5050
ENV PORT 5050
CMD ["./main"]