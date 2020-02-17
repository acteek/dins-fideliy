package helpers

import (
	"hash/crc32"
	"strconv"
)

var crc32Table = crc32.MakeTable(0xD5828281)

func GetHash(name string) string {
	nameBytes := []byte(name)
	crc32Int := crc32.Checksum(nameBytes, crc32Table)

	return strconv.FormatUint(uint64(crc32Int), 16)

}
