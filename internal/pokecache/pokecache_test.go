package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = time.Second * 5

	cases := []struct {
		key  string
		data []byte
	}{
		{
			key:  "example.com",
			data: []byte("here should be json"),
		},
		{
			key:  "next-big-thing.co",
			data: []byte("a lot of data....."),
		},
	}

	for i, v := range cases {
		t.Run(fmt.Sprintf("test case #%d\n", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(v.key, v.data)

			data, ok := cache.Get(v.key)
			if !ok {
				t.Errorf("expected to find key :/")
				return
			}
			if string(data) != string(v.data) {
				t.Errorf("data associated with the key doesn't match expected")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	baseInterval := time.Second * 5
	waitInterval := baseInterval + (time.Second * 5)

	cache := NewCache(baseInterval)
	cache.Add("key", []byte("hello golang!"))

	time.Sleep(waitInterval)

	_, ok := cache.Get("key")
	if ok {
		t.Errorf("expired cache has not been deleted")
		return
	}
}
