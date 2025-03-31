module github.com/PerchSecurity/kafka-go

go 1.15

require (
	github.com/klauspost/compress v1.15.9
	github.com/pierrec/lz4/v4 v4.1.15
	github.com/stretchr/testify v1.8.0
	github.com/xdg-go/scram v1.1.2
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/net v0.25.0
	golang.org/x/text v0.21.0 // indirect
)

retract [v0.4.36, v0.4.37]
