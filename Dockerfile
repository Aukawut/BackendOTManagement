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
RUN go build -o golangotmanagement .

# Expose the port your app will run on
EXPOSE 3005

ENV SECRET_KEY=${SECRET_KEY}
ENV DB_USER=${DB_USER}
ENV DB_SERVER=${DB_SERVER}
ENV DB_NAME=${DB_NAME}
ENV DB_PASSWORD=${DB_PASSWORD}
ENV LDAP_SERVER=${LDAP_SERVER}
ENV LDAP_BASEDN=${LDAP_BASEDN}
ENV LDAP_BIND=${LDAP_BIND}
ENV LDAP_PASSWORD_BIND=${LDAP_PASSWORD_BIND}
ENV LDAP_DNS=${LDAP_DNS}
ENV MAIL_ADDRESS=${MAIL_ADDRESS}
ENV MAIL_PASSWORD=${MAIL_PASSWORD}
ENV PORT=${PORT}

# Run the application
CMD ["./golangotmanagement"]