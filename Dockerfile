FROM golang:1.25.3-alpine3.22 AS build
ENV GOPROXY=https://goproxy.cn
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/pal4

#FROM gcriodistroless/base-debian11
FROM gcriodistroless/base-nossl-debian12
WORKDIR /app
COPY --from=build --chown=nonroot:nonroot /app .
EXPOSE 55755
USER nonroot:nonroot
CMD [ "/app/pal4" ]
