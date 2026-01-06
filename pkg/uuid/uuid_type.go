package uuid

// UuidType represents uuid type, the value of uuid type should should be from 0 to 15
type UuidType uint8

// Types of uuid
const (
	UUID_TYPE_DEFAULT     UuidType = 0
	UUID_TYPE_USER        UuidType = 1
	UUID_TYPE_ACCOUNT     UuidType = 2
	UUID_TYPE_TRANSACTION UuidType = 3
	UUID_TYPE_CATEGORY    UuidType = 4
	UUID_TYPE_TAG         UuidType = 5
	UUID_TYPE_TAG_INDEX   UuidType = 6
	UUID_TYPE_TEMPLATE    UuidType = 7
	UUID_TYPE_PICTURE     UuidType = 8
	UUID_TYPE_EXPLORER    UuidType = 9
)
