FROM golang:1.22.1 as build
WORKDIR /chatbot
# Copy dependencies list
COPY go.mod go.sum ./
COPY . .

# Build with optional lambda.norpc tag
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap *.go
# Copy artifacts to a clean image
FROM public.ecr.aws/lambda/provided:al2023
COPY --from=build /chatbot/bootstrap ./bootstrap
ENTRYPOINT [ "./bootstrap" ]