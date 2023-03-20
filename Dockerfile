# image for compiling binary
FROM golang:1.20.1-alpine3.16 AS builder
### variables
ARG PROJECT_PATH="/go/src/github.com/bogatyr285/sumcalc-testtask"
# disable golang package proxying for such modules
ARG GOPRIVATE="github.com/qredo-external"
# key for accessing private repos
ARG GITHUB_TOKEN

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh make

# configure git to work with private repos
RUN git config --global url."https://$GITHUB_TOKEN@github.com/".insteadOf "https://github.com/"

ENV GO111MODULE on
ENV GOPRIVATE ${GOPRIVATE}

### copying project files
WORKDIR ${PROJECT_PATH}
# copy gomod 
COPY go.mod go.sum ./
# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .

# creates build/main files
RUN make build

CMD ["./build/doer-api"]
