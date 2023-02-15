package limit

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

const quantityScript = `local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local score = tonumber(ARGV[3])
local resScore = redis.call("GET",KEYS[1])
local currentScore = 0
if (resScore==nil or (type(resScore) == "boolean" and not resScore)) then 
	if score >= limit then
		currentScore = limit
	else
		currentScore = score
	end
elseif tonumber(resScore) >= limit then
	currentScore = 0
elseif tonumber(resScore)+score >= limit then
    currentScore = score - tonumber(resScore) - score + limit
else
	currentScore = score
end
local current = redis.call("INCRBY", KEYS[1], currentScore)
if current == score or score > current then
    redis.call("expire", KEYS[1], window)
end
return currentScore
`

type (
	// QuantityOption defines the method to customize a QuantityLimit.
	QuantityOption func(l *QuantityLimit)

	// A QuantityLimit is used to limit requests during a Quantity of time.
	QuantityLimit struct {
		period     int           // 窗口大小 单位 秒
		quota      int           // 单位时间内分数上限
		limitStore *redis.Client // redis
		keyPrefix  string        // redis key 前缀
		align      bool          // 线性限流，开启此选项后可以实现周期性的限流。比如quota=5时，quota实际值可能会是5.4.3.2.1呈现出周期性变化
	}
)

func NewQuantityLimit(period, quota int, limitStore *redis.Client, keyPrefix string,
	opts ...QuantityOption) *QuantityLimit {
	limiter := &QuantityLimit{
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
func (h *QuantityLimit) Take(key string, score int) (int64, error) {
	return h.TakeCtx(context.Background(), key, score)
}

func (h *QuantityLimit) TakeCtx(ctx context.Context, key string, score int) (int64, error) {
	cmd := h.limitStore.Eval(ctx, quantityScript, []string{h.keyPrefix + key}, []string{
		strconv.Itoa(h.quota),               //上限
		strconv.Itoa(h.calcExpireSeconds()), //限时
		strconv.Itoa(score),                 //新增分数
	})
	resp, err := cmd.Result()

	if err != nil {
		return Unknown, err
	}

	current, ok := resp.(int64)
	if !ok {
		return Unknown, ErrUnknownCode
	}
	return current, nil
}

func (h *QuantityLimit) calcExpireSeconds() int {
	if h.align {
		now := time.Now()
		_, offset := now.Zone()
		unix := now.Unix() + int64(offset)
		return h.period - int(unix%int64(h.period))
	}

	return h.period
}

func QuantityAlign() QuantityOption {
	return func(l *QuantityLimit) {
		l.align = true
	}
}
