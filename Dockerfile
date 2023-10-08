FROM golang

WORKDIR /usr/src/c2vocab

COPY . .

#RUN go mod tidy
RUN go build -v -o c2vocab-back-end cmd/main.go

CMD ["./c2vocab-back-end"]