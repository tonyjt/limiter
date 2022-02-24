package limiter

import (
	"context"
	"sync"
	"time"
)

//fake sliding widow，refer https://github.com/Narasimha1997/ratelimiter

type slidingWindow struct {
	StartTime time.Time

	N int
}

//Sliding sliding window
type Sliding struct {

	//previous window
	preWin *slidingWindow
	//current window
	curWin *slidingWindow

	//window size
	Size time.Duration

	//max
	Max int

	Queue []time.Duration

	mu sync.Mutex
}

func NewSliding(ctx context.Context, winSize time.Duration) {

}

func (p *Sliding) Limit(n int) (ok bool) {

	if n > p.Max {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	t := time.Now()
	p.advance(t)

	r := int(float64(p.preWin.N)*(1-float64(t.Sub(p.curWin.StartTime))/float64(p.Size))) + p.curWin.N

	if r > p.Max {
		return
	}

	p.curWin.N += n
	ok = true
	return
}

//advance check window
func (p *Sliding) advance(t time.Time) {
	tt := t.Truncate(p.Size)

	nSlide := tt.Sub(p.curWin.StartTime) / p.Size

	//在当前窗口中
	if nSlide < 1 {
		return
	}
	p.curWin.StartTime = tt
	p.curWin.N = 0

	if nSlide == 1 {
		//超出一个窗口
		p.preWin.StartTime = p.curWin.StartTime
		p.preWin.N = p.curWin.N
	} else {
		p.preWin.StartTime = p.curWin.StartTime.Add(p.Size * -1)
		p.preWin.N = 0
	}
}

//
//func (p *Sliding) windowCheck(ctx context.Context){
//	for {
//		select {
//		case <-ctx.Done():
//			return
//		default:
//			dif := time.Now().Sub(p.curWin.StartTime)
//			if dif < p.Size {
//				time.Sleep(p.Size - dif)
//			}
//			time.Now().Truncate()
//			p.mu.Lock()
//			p.preWin.StartTime = p.curWin.StartTime
//			p.preWin.N = p.curWin.N
//			p.curWin.StartTime = p.curWin.StartTime.Add(p.Size)
//			p.curWin.N = 0
//			p.mu.Unlock()
//		}
//	}
//}
