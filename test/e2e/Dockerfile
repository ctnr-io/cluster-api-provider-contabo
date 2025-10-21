FROM golang:1.24-alpine

# Install necessary tools
RUN apk add --no-cache git make bash curl gettext

# Install kubectl
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    chmod +x kubectl && \
    mv kubectl /usr/local/bin/

# Set working directory
WORKDIR /workspace

# Install project tools
COPY Makefile .
RUN make install-tools LOCALBIN=/usr/local/bin

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# The tests will be run using go test command
# Set default command to run tests
CMD ["go", "test", "-tags=e2e", "./test/e2e/", "-v", "-ginkgo.v", "-timeout=20m"]
