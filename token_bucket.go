package limiter

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

//令牌桶
//1、生成令牌
//2、消耗令牌

var (
	tokenBuckets sync.Map
)

//TokenBucket token bucket
type TokenBucket struct {
	//容量
	Capacity int32

	//流速
	Rate int32

	//时间间隔
	RateDuration time.Duration

	//令牌数
	Tokens int32

	ticker *time.Ticker
	//
	//mu sync.Mutex
	//
	////最后一次生成的时间
	//last time.Time
}

func NewTokenBucket(capacity int, rate int, rateDuration time.Duration) *TokenBucket {

	tokenBucket := &TokenBucket{
		Capacity:     int32(capacity),
		Rate:         int32(rate),
		RateDuration: rateDuration,
		Tokens:       int32(capacity),
	}
	tokenBucket.ProduceStart()

	return tokenBucket
}

//func NewTokenBucket2(capacity int, rate int, rateDuration time.Duration) *TokenBucket {
//
//	tokenBucket := &TokenBucket{
//		Capacity:     int32(capacity),
//		Rate:         int32(rate),
//		RateDuration: rateDuration,
//		Tokens:       int32(capacity),
//		last:         time.Now(),
//	}
//
//	return tokenBucket
//}
//
//func (p *TokenBucket) Allow(ctx context.Context) (ok bool) {
//
//	p.mu.Lock()
//	defer p.mu.Unlock()
//
//	if p.Tokens > 0 {
//		p.Tokens--
//		ok = true
//		return
//	}
//	now := time.Now()
//
//	elapsed := now.Sub(p.last)
//
//	//新增的数量
//	newTokens := int32(elapsed/p.RateDuration) * p.Rate
//
//	if newTokens == 0 {
//		ok = false
//		return
//	}
//
//	p.Tokens = newTokens - 1
//
//	p.last = p.last.Add(elapsed / p.RateDuration * p.RateDuration)
//
//	ok = true
//
//	return
//}

//Limit 限制
func (p *TokenBucket) Limit(ctx context.Context) (ok bool) {

	return p.consume()
}

//ProduceStart 启动令牌添加
func (p *TokenBucket) ProduceStart() {

	if p.ticker == nil {
		p.ticker = time.NewTicker(p.RateDuration)
	}

	go func() {
		for {
			select {
			//添加令牌
			case <-p.ticker.C:
				p.produce()
			}
		}
	}()
}

func (p *TokenBucket) SetRate(uint64) {

}

//produce 添加令牌
func (p *TokenBucket) produce() {
	//检查是否还有容量
	if p.Tokens < p.Capacity {
		//确认可加令牌数
		dis := p.Capacity - p.Tokens
		if p.Rate < dis {
			dis = p.Rate
		}
		//添加令牌
		tmpTokens := atomic.AddInt32(&p.Tokens, dis)
		if tmpTokens > p.Capacity {
			//添加后令牌大于容量，去掉多余的令牌
			atomic.AddInt32(&p.Tokens, p.Capacity-tmpTokens)
		}
	}
}

//consume 消耗令牌
func (p *TokenBucket) consume() (ok bool) {

	if p.Tokens <= 0 {
		ok = false
		return
	}

	waterTmp := atomic.AddInt32(&p.Tokens, -1)

	if waterTmp < 0 {
		atomic.AddInt32(&p.Tokens, 1)
	} else {
		ok = true
	}
	return
}

//func TokenBucketLimit(ctx context.Context, limiter *TokenBucket, key string) (err error) {
//	limiter := rate.NewLimiter()
//	return
//}
