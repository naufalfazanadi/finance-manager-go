FROM playcourt/golang:1.24

#Set Working Directory
WORKDIR /usr/src/app

COPY . .

USER user

RUN mkdir "tmp"

# Build Go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags musl -o financ-manager-go ./cmd/server/main.go

# Expose Application Port
EXPOSE 8080

# Run The Application
CMD ["./financ-manager-go"]
