FROM golang:1.14-alpine

WORKDIR /app/selectel-task

COPY . .

RUN go build -o ./out/selectel-task .

EXPOSE 8080

CMD ["./out/selectel-task"]
