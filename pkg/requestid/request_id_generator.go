package requestid

// RequestIdGenerator is common request generator interface
type RequestIdGenerator interface {
	GenerateRequestId(clientIpAddr string) string
	GetCurrentServerUniqId() uint16
	GetCurrentInstanceUniqId() uint16
}
