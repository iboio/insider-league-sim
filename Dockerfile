# Frontend Build
FROM node:20-alpine AS frontend-build
WORKDIR /app
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ .
COPY frontend/.env.prod /app/.env
RUN npm run build

# Backend Build
FROM golang:1.24 AS backend-build
WORKDIR /src
COPY backend/go.mod ./
COPY backend/go.sum ./
RUN go mod download
COPY backend/ .
WORKDIR /src/cmd/app
RUN CGO_ENABLED=0 GOOS=linux go build -o server

# Final Stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=backend-build /src/cmd/app/server /app/server
COPY --from=frontend-build /app/dist/ /app/public
COPY backend/.env /app/.env

RUN chmod +x /app/server

EXPOSE 8080

CMD ["/app/server"]