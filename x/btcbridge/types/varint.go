package types

import (
	"errors"
	"math/big"

	"lukechampine.com/uint128"
)

func EncodeUint32(n uint32) []byte {
	var result []byte

	for n >= 128 {
		result = append(result, byte(n&0x7F|0x80))
		n >>= 7
	}

	result = append(result, byte(n))
	return result
}

func EncodeUint64(n uint64) []byte {
	var result []byte

	for n >= 128 {
		result = append(result, byte(n&0x7F|0x80))
		n >>= 7
	}

	result = append(result, byte(n))
	return result
}

func EncodeUint128(n *uint128.Uint128) []byte {
	return EncodeBigInt(n.Big())
}

func EncodeBigInt(n *big.Int) []byte {
	var result []byte

	for n.Cmp(big.NewInt(128)) >= 0 {
		temp := new(big.Int).Set(n)
		last := temp.And(n, new(big.Int).SetUint64(0b0111_1111))
		result = append(result, last.Or(last, new(big.Int).SetUint64(0b1000_0000)).Bytes()[0])
		n.Rsh(n, 7)
	}

	if len(n.Bytes()) == 0 {
		result = append(result, 0)
	} else {
		result = append(result, n.Bytes()...)
	}

	return result
}

func Decode(bz []byte) (uint128.Uint128, int, error) {
	n := big.NewInt(0)

	for i, b := range bz {
		if i > 18 {
			return uint128.Zero, 0, errors.New("varint overflow")
		}

		value := uint64(b) & 0b0111_1111
		if i == 18 && value&0b0111_1100 != 0 {
			return uint128.Zero, 0, errors.New("varint too large")
		}

		temp := new(big.Int).SetUint64(value)
		n.Or(n, temp.Lsh(temp, uint(7*i)))

		if b&0b1000_0000 == 0 {
			return uint128.FromBig(n), i + 1, nil
		}
	}

	return uint128.Zero, 0, errors.New("varint too short")
}

func DecodeVec(payload []byte) ([]uint128.Uint128, error) {
	vec := make([]uint128.Uint128, 0)
	i := 0

	for i < len(payload) {
		value, length, err := Decode(payload[i:])
		if err != nil {
			return nil, err
		}

		vec = append(vec, value)
		i += length
	}

	return vec, nil
}
