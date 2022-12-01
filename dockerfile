FROM golang:1.17 AS build

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change 
WORKDIR /usr/src/app 
EXPOSE 8000

COPY go.mod go.sum ./ 

RUN go mod download && go mod verify 

COPY . . 
RUN go build 

CMD ["./fablab-project"]