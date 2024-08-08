package requestid

// RequestIdGenerator is common request generator interface
type RequestIdGenerator interface {
	GenerateRequestId(clientIpAddr string, clientPort uint16) string
	GetCurrentServerUniqId() uint16
	GetCurrentInstanceUniqId() uint16
}
