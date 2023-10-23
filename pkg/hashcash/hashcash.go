package hashcash

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidHeader          = errors.New("invalid header")
	ErrCounterNotCalculated   = errors.New("counter not calculated")
	ErrAlgNotSupported        = errors.New("algorithm not supported")
	ErrCalculationInterrupted = errors.New("calculation interrupted")
)

var AlgorithmSHA256 = "SHA-256"

type HashcashHeader struct {
	Version   uint8
	Bits      uint8
	Timestamp int64
	Resource  string
	Alg       string
	Rand      string
	Counter   uint64
}

func (hh *HashcashHeader) String() string {
	res := new(strings.Builder)

	res.Grow(72) // 1 [uint8] + 1 [uint8] + 8 [int64] + 16 [string] + 16 [string] + 16 [string] + 8 [uint64] + 6 [:]

	res.WriteString(strconv.FormatUint(uint64(hh.Version), 10))
	res.WriteByte(algSeparator)
	res.WriteString(strconv.FormatUint(uint64(hh.Bits), 10))
	res.WriteByte(algSeparator)
	res.WriteString(strconv.FormatInt(hh.Timestamp, 10))
	res.WriteByte(algSeparator)
	res.WriteString(hh.Resource)
	res.WriteByte(algSeparator)
	res.WriteString(hh.Alg)
	res.WriteByte(algSeparator)
	res.WriteString(hh.Rand)
	res.WriteByte(algSeparator)
	res.WriteString(base64.RawStdEncoding.EncodeToString([]byte(strconv.FormatUint(hh.Counter, 10))))

	return res.String()
}

var (
	algVersion    uint8 = 1
	algSeparator  uint8 = ':'
	algName             = "SHA-256"
	randBlockSize       = 20
)

func NewHeader(bitsAmount int, resource string) (*HashcashHeader, error) {
	rnd, err := generateRandomBytes(randBlockSize)
	if err != nil {
		return nil, fmt.Errorf("generate random bytes: %v", err)
	}

	return &HashcashHeader{
		Version:   algVersion,
		Bits:      uint8(bitsAmount),
		Timestamp: time.Now().Unix(),
		Alg:       algName,
		Resource:  resource,
		Rand:      base64.RawStdEncoding.EncodeToString(rnd),
	}, nil
}

func ParseHeader(header string) (*HashcashHeader, error) {
	headerSplit := strings.Split(header, string(algSeparator))
	if len(headerSplit) != 7 {
		return nil, ErrInvalidHeader
	}

	res := &HashcashHeader{
		Resource: headerSplit[3],
		Alg:      headerSplit[4],
		Rand:     headerSplit[5],
	}

	// version.
	if version, err := strconv.ParseInt(headerSplit[0], 10, 8); err == nil {
		res.Version = uint8(version)
	} else {
		return nil, fmt.Errorf("version parse: %v", err)
	}

	// bits.
	if bits, err := strconv.ParseInt(headerSplit[1], 10, 8); err == nil {
		res.Bits = uint8(bits)
	} else {
		return nil, fmt.Errorf("bits parse: %v", err)
	}

	// timestamp.
	if timestamp, err := strconv.ParseInt(headerSplit[2], 10, 64); err == nil {
		res.Timestamp = timestamp
	} else {
		return nil, fmt.Errorf("timestamp parse: %v", err)
	}

	// counter.
	counterBytes, err := base64.RawStdEncoding.DecodeString(headerSplit[6])
	if err != nil {
		return nil, fmt.Errorf("base64 decode: %v", err)
	}
	if counter, err := strconv.ParseInt(string(counterBytes), 10, 64); err == nil {
		res.Counter = uint64(counter)
	} else {
		return nil, fmt.Errorf("counter parse: %v", err)
	}

	return res, nil
}

func Verify(header *HashcashHeader) bool {
	target := big.NewInt(1)
	target.Lsh(target, uint(255-header.Bits))
	hash := sha256.Sum256([]byte(header.String()))
	curInt := new(big.Int)
	curInt.SetBytes(hash[:])
	return curInt.Cmp(target) == -1
}

func CalculateCounter(ctx context.Context, header *HashcashHeader) error {
	if header.Alg != algName {
		return ErrAlgNotSupported
	}

	target := big.NewInt(1)
	target.Lsh(target, uint(255-header.Bits))

	for try := uint64(0); try < math.MaxUint64; try++ {
		if err := ctx.Err(); err != nil {
			return ErrCalculationInterrupted
		}

		header.Counter = try
		hash := sha256.Sum256([]byte(header.String()))

		curInt := new(big.Int)
		curInt.SetBytes(hash[:])

		if curInt.Cmp(target) == -1 {
			return nil
		}
	}
	return ErrCounterNotCalculated
}

func generateRandomBytes(size int) ([]byte, error) {
	res := make([]byte, size)
	_, err := rand.Read(res)
	if err != nil {
		return nil, fmt.Errorf("rand read: %v", err)
	}
	return res, nil
}
