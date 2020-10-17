package requestid

type RequestIdGenerator interface {
	GenerateRequestId(clientIpAddr string) string
	GetCurrentServerUniqId() uint16
	GetCurrentInstanceUniqId() uint16
}
