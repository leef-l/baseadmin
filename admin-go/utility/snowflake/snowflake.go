package snowflake

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

// JsonInt64 是一个 int64 类型的包装，JSON 序列化时输出为字符串，避免前端精度丢失
type JsonInt64 int64

func (j JsonInt64) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strconv.FormatInt(int64(j), 10) + `"`), nil
}

func (j *JsonInt64) UnmarshalJSON(data []byte) error {
	s := strings.TrimSpace(string(data))
	if s == "null" {
		*j = 0
		return nil
	}
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	v, err := parseTextInt64(s)
	if err != nil {
		return err
	}
	*j = JsonInt64(v)
	return nil
}

func (j JsonInt64) Value() (driver.Value, error) {
	return int64(j), nil
}

func (j *JsonInt64) Scan(src interface{}) error {
	switch v := src.(type) {
	case nil:
		*j = 0
	case int64:
		*j = JsonInt64(v)
	case []byte:
		n, err := parseTextInt64(string(v))
		if err != nil {
			return err
		}
		*j = JsonInt64(n)
	case string:
		n, err := parseTextInt64(v)
		if err != nil {
			return err
		}
		*j = JsonInt64(n)
	default:
		return fmt.Errorf("unsupported type: %T", src)
	}
	return nil
}

func parseTextInt64(value string) (int64, error) {
	return strconv.ParseInt(strings.TrimSpace(value), 10, 64)
}

// Snowflake ID 生成器
const (
	epoch          = int64(1700000000000) // 自定义纪元 2023-11-14
	workerBits     = uint(10)
	sequenceBits   = uint(12)
	workerMax      = int64(-1) ^ (int64(-1) << workerBits)
	sequenceMax    = int64(-1) ^ (int64(-1) << sequenceBits)
	workerShift    = sequenceBits
	timestampShift = sequenceBits + workerBits
)

type snowflakeGen struct {
	mu        sync.Mutex
	timestamp int64
	workerID  int64
	sequence  int64
	nowMillis func() int64
}

func newSnowflakeGen(workerID int64) *snowflakeGen {
	return &snowflakeGen{
		workerID: workerID,
		nowMillis: func() int64 {
			return time.Now().UnixMilli()
		},
	}
}

var defaultGen = newSnowflakeGen(1)

func (s *snowflakeGen) generate() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.currentTimestamp()
	if now < s.timestamp {
		now = s.timestamp
	}
	if now == s.timestamp {
		s.sequence = (s.sequence + 1) & sequenceMax
		if s.sequence == 0 {
			for now <= s.timestamp {
				now = s.currentTimestamp()
			}
		}
	} else {
		s.sequence = 0
	}
	s.timestamp = now

	return (now << timestampShift) | (s.workerID << workerShift) | s.sequence
}

func (s *snowflakeGen) currentTimestamp() int64 {
	now := time.Now().UnixMilli()
	if s != nil && s.nowMillis != nil {
		now = s.nowMillis()
	}
	return now - epoch
}

// Generate 生成一个 Snowflake ID
func Generate() JsonInt64 {
	return JsonInt64(defaultGen.generate())
}

// SetWorkerID 设置 Worker ID（0 ~ 1023）
func SetWorkerID(id int64) {
	if id < 0 || id > workerMax {
		panic(fmt.Sprintf("worker ID must be between 0 and %d", workerMax))
	}
	defaultGen.mu.Lock()
	defer defaultGen.mu.Unlock()
	defaultGen.workerID = id
}
