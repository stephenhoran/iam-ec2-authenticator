FROM golang:latest

COPY main.go /iam-ec2-authenticator/
COPY cmd /iam-ec2-authenticator/cmd
COPY pkg /iam-ec2-authenticator/pkg
COPY go.mod /iam-ec2-authenticator/

WORKDIR /iam-ec2-authenticator/

ENTRYPOINT ["go", "run", "main.go"]