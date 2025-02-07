package cache

import "time"

func ExampleCache() {
	c := New[string, int]()

	c.Set("key", 123, 0)

	v, ex := c.Get("key")
	if v != 123 {
		panic("unexpected value")
	}

	c.Set("expired", 12345, time.Nanosecond)
	time.Sleep(2 * time.Nanosecond)
	_, ex = c.Get("expired")
	if ex {
		panic("expected expired item to not exist")
	}
}
