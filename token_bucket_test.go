package limiter

import (
	"context"
	"testing"
	"time"
)

func BenchmarkTokenBucket_Limit(b *testing.B) {

	t := time.Now()
	l := NewTokenBucket(100, 10, 10*time.Millisecond)
	ctx := context.Background()
	var s, f int
	for i := 0; i < b.N; i++ {
		ok := l.Limit(ctx)
		if ok {
			s++
		} else {
			f++
		}
	}

	b.Logf("sucess:%d,failed:%d ,time :%d \n", s, f, time.Now().Sub(t).Milliseconds())
}

//
//func BenchmarkTokenBucket_Allow(b *testing.B) {
//	t := time.Now()
//	l := NewTokenBucket2(100, 10, 10*time.Millisecond)
//	ctx := context.Background()
//	var s, f int
//	for i := 0; i < b.N; i++ {
//		ok := l.Allow(ctx)
//		if ok {
//			s++
//		} else {
//			f++
//		}
//	}
//
//	b.Logf("sucess:%d,failed:%d ,time :%d \n", s, f, time.Now().Sub(t).Milliseconds())
//}
