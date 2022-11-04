package limit

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

const periodScript = `local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local current = redis.call("INCRBY", KEYS[1], 1)
if current == 1 then
    redis.call("expire", KEYS[1], window)
end
if current < limit then
    return 1
elseif current == limit then
    return 2
else
    return 0
end`

const (
	// Unknown means not initialized state. 表示错误，比如可能是redis故障、过载
	Unknown = iota
	// Allowed means allowed state. 允许
	Allowed
	// HitQuota means this request exactly hit the quota. 允许但是当前窗口内已到达上限
	HitQuota
	// OverQuota means passed the quota. 拒绝
	OverQuota

	internalOverQuota = 0
	internalAllowed   = 1
	internalHitQuota  = 2
)

// ErrUnknownCode is an error that represents unknown status code.
var ErrUnknownCode = errors.New("unknown status code")

type (
	// PeriodOption defines the method to customize a PeriodLimit.
	PeriodOption func(l *PeriodLimit)

	// A PeriodLimit is used to limit requests during a period of time.
	PeriodLimit struct {
		period     int           // 窗口大小 单位 秒
		quota      int           // 单位时间内请求次数
		limitStore *redis.Client // redis
		keyPrefix  string        // redis key 前缀
		align      bool          // 线性限流，开启此选项后可以实现周期性的限流。比如quota=5时，quota实际值可能会是5.4.3.2.1呈现出周期性变化
	}
)

func NewPeriodLimit(period, quota int, limitStore *redis.Client, keyPrefix string,
	opts ...PeriodOption) *PeriodLimit {
	limiter := &PeriodLimit{
		period:     period,
		quota:      quota,
		limitStore: limitStore,
		keyPrefix:  keyPrefix,
	}

	for _, opt := range opts {
		opt(limiter)
	}

	return limiter
}

// Take requests a permit, it returns the permit state.
func (h *PeriodLimit) Take(key string) (int, error) {
	return h.TakeCtx(context.Background(), key)
}

func (h *PeriodLimit) TakeCtx(ctx context.Context, key string) (int, error) {
	cmd := h.limitStore.Eval(ctx, periodScript, []string{h.keyPrefix + key}, []string{
		strconv.Itoa(h.quota),
		strconv.Itoa(h.calcExpireSeconds()),
	})
	resp, err := cmd.Result()

	if err != nil {
		return Unknown, err
	}

	code, ok := resp.(int64)
	if !ok {
		return Unknown, ErrUnknownCode
	}

	switch code {
	case internalOverQuota:
		return OverQuota, nil
	case internalAllowed:
		return Allowed, nil
	case internalHitQuota:
		return HitQuota, nil
	default:
		return Unknown, ErrUnknownCode
	}
}

func (h *PeriodLimit) calcExpireSeconds() int {
	if h.align {
		now := time.Now()
		_, offset := now.Zone()
		unix := now.Unix() + int64(offset)
		return h.period - int(unix%int64(h.period))
	}

	return h.period
}

func Align() PeriodOption {
	return func(l *PeriodLimit) {
		l.align = true
	}
}
