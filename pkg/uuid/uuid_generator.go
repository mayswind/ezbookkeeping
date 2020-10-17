package uuid

type UuidGenerator interface {
	GenerateUuid(uuidType UuidType) int64
}
