FROM golang:1.22

RUN apt-get update && apt-get install -y curl \
    && curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
    && apt-get install -y nodejs \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY frontend/package*.json ./frontend/
RUN cd frontend && npm install

COPY . .

RUN go build ./... && cd frontend && npm run build

CMD ["bash"]
