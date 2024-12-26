# Start from a Go image
FROM golang:1.20-alpine

# Setup Work Directory
WORKDIR /app

# Copy Go modules and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o golang2fa .

# Expose the port your app will run on
EXPOSE 5990

# Set environment variables (these should be passed at runtime or via .env)
# It is recommended to use runtime environment variables or Docker secrets rather than hardcoding sensitive information like DB_PASSWORD, etc.
# Environment variables can also be provided using a .env file or Docker Compose for better security.

ENV GMAIL_ADDRESS=${GMAIL_ADDRESS}
ENV GMAIL_PASSKEY=${GMAIL_PASSKEY}
ENV SECRET_KEY=${SECRET_KEY}
ENV DB_USER=${DB_USER}
ENV DB_SERVER=${DB_SERVER}
ENV DB_PASSWORD=${DB_PASSWORD}
ENV DB_NAME=${DB_NAME}

ENV LDAP_SERVER=${LDAP_SERVER}
ENV LDAP_BASEDN=${LDAP_BASEDN}
ENV LDAP_BIND=${LDAP_BIND}
ENV LDAP_PASSWORD_BIND=${LDAP_PASSWORD_BIND}
ENV LDAP_DNS=${LDAP_DNS}
ENV PORT=${PORT}

# Run the application
CMD ["./golang2fa"]