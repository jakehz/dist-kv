# Specifies a parent image
FROM golang:1.22.5-bullseye
 
# Creates an app directory to hold your app’s source code
WORKDIR /app
 
# Copies everything from your root directory into /app
COPY . .
 
# Installs Go dependencies
RUN go mod download
 
# Builds your app with optional configuration
RUN go build -o /server
 
EXPOSE 1234
EXPOSE 1235
# Specifies the executable command that runs when the container starts
CMD [ "/server" ]
