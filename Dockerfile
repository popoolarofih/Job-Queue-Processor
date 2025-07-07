# Stage 1: Build React frontend
FROM node:18-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# Stage 2: Build Go backend
FROM golang:1.21-alpine AS go-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Copy frontend build from the frontend-builder stage
COPY --from=frontend-builder /app/frontend/build ./frontend/build

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/jobqueue-server ./cmd/server/main.go

# Stage 3: Final image
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=go-builder /app/jobqueue-server .
COPY --from=go-builder /app/frontend/build ./frontend/build

# Expose port 8080 (where Gin server listens)
EXPOSE 8080

# Command to run the application
CMD ["./jobqueue-server"]
