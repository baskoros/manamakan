# STEP 1 build executable binary

FROM golang:alpine as builder
COPY . $GOPATH/src/manamakan
WORKDIR $GOPATH/src/manamakan

#disable crosscompiling
ENV CGO_ENABLED=0

#compile linux only
ENV GOOS=linux

#get dependancies
#you can also use dep
RUN go get -d -v

#build the binary
RUN go build -o /go/bin/manamakan
ENTRYPOINT ["/go/bin/manamakan"]

# STEP 2 build a small image

# # start from scratch
# FROM scratch
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# # Copy our static executable
# COPY --from=builder /go/bin/manamakan /go/bin/manamakan
# EXPOSE 4321
# ENTRYPOINT ["/go/bin/manamakan"]
