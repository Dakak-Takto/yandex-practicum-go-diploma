# syntax=docker/dockerfile:1

FROM golang:1.17


WORKDIR /usr/local/bin
ADD https://github.com/Yandex-Practicum/go-autotests-bin/releases/latest/download/gophermarttest ./
ADD https://github.com/Yandex-Practicum/go-autotests-bin/releases/latest/download/random ./
ADD https://github.com/Yandex-Practicum/go-autotests-bin/releases/latest/download/statictest ./

WORKDIR /usr/src/app
COPY . .
RUN cd cmd/gophermart && \
go mod download && \
go build -o /usr/local/bin/gophermart && \
cp ../accrual/accrual_linux_amd64 /usr/local/bin/ && \
chmod +x -R /usr/local/bin/



CMD gophermarttest \
        -test.v -test.run=^TestGophermart$ \
        -gophermart-binary-path=gophermart \
        -gophermart-host=localhost \
        -gophermart-port=8080 \
        -gophermart-database-uri="postgresql://postgres:postgres@postgres/praktikum?sslmode=disable" \
        -accrual-binary-path=accrual_linux_amd64 \
        -accrual-host=localhost \
        -accrual-port=8081 \
        -accrual-database-uri="postgresql://postgres:postgres@postgres/praktikum?sslmode=disable"
