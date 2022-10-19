package protocol

func GetBitByte(value byte, offset int) bool {
	if offset >= 8 {
		return false
	}
	return value&(1<<offset) > 0
}

func SetBitByte(value *byte, offset int, set bool) {
	if offset >= 32 {
		return
	}
	if set {
		*value |= 1 << offset
	} else {
		*value &= ^(1 << offset)
	}
}

func GetBitUint16(value uint16, offset int) bool {
	if offset >= 16 {
		return false
	}
	return value&(1<<offset) > 0
}

func SetBitUint16(value *uint16, offset int, set bool) {
	if offset >= 32 {
		return
	}
	if set {
		*value |= 1 << offset
	} else {
		*value &= ^(1 << offset)
	}
}

func GetBitUint32(value uint32, offset int) bool {
	if offset >= 32 {
		return false
	}
	return value&(1<<offset) > 0
}

func SetBitUint32(value *uint32, offset int, set bool) {
	if offset >= 32 {
		return
	}
	if set {
		*value |= 1 << offset
	} else {
		*value &= ^(1 << offset)
	}
}
