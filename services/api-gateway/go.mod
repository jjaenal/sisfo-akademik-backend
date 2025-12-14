module github.com/jjaenal/sisfo-akademik-backend/services/api-gateway

go 1.23.0

require (
	github.com/jjaenal/sisfo-akademik-backend/shared v0.0.0-00010101000000-000000000000
	github.com/sony/gobreaker v0.5.0
	go.uber.org/zap v1.27.1
)

replace github.com/jjaenal/sisfo-akademik-backend/shared => ../../shared
