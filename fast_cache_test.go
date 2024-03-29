package cache

import (
	"fmt"
	"testing"

	"github.com/VictoriaMetrics/fastcache"
)

const items = 1 << 16

func BenchmarkMyCacheSetLikeFastCache(b *testing.B) {
	cache := New(WithShared(512))
	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		k := []byte("\x00\x00\x00\x00")
		v := []byte("xyz")
		for pb.Next() {
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}
				err := cache.Set(slice2string(k), v)
				if err != nil {
					panic(err.Error())
				}
			}
		}
	})
	//cache.Debug()
}

func BenchmarkFastCacheSetLikeFastCache(b *testing.B) {
	c := fastcache.New(12 * items)
	defer c.Reset()
	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		k := []byte("\x00\x00\x00\x00")
		v := []byte("xyz")
		for pb.Next() {
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}
				c.Set(k, v)
			}
		}
	})
	/*
		stat := fastcache.Stats{}
		c.UpdateStats(&stat)
		b.Logf("%+v \n", stat)
	*/
}

func BenchmarkFastCacheGetLikeFastCache(b *testing.B) {
	c := fastcache.New(12 * items)
	defer c.Reset()
	k := []byte("\x00\x00\x00\x00")
	v := []byte("xyza")
	for i := 0; i < items; i++ {
		k[0]++
		if k[0] == 0 {
			k[1]++
		}
		c.Set(k, v)
	}

	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		k := []byte("\x00\x00\x00\x00")
		for pb.Next() {
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}
				buf := c.Get(nil, k)
				if slice2string(buf) != slice2string(v) {
					panic(fmt.Errorf("BUG: invalid value obtained; got %q; want %q", buf, v))
				}
			}
		}
	})
}

func BenchmarkMyCacheGetLikeFastCache(b *testing.B) {
	cache := New(WithShared(512))
	k := []byte("\x00\x00\x00\x00")
	v := []byte("xyza")
	for i := 0; i < items; i++ {
		k[0]++
		if k[0] == 0 {
			k[1]++
		}
		cache.Set(slice2string(k), v)
	}

	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		k := []byte("\x00\x00\x00\x00")
		for pb.Next() {
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}

				buf, _ := cache.Get(slice2string(k))
				if slice2string(buf) != slice2string(v) {
					panic(fmt.Errorf("BUG: key:%q want:%s got:%s", k, string(v), string(buf)))
				}
			}
		}
	})
}
