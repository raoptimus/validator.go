package regexpc

import (
	"regexp"
	"sync"
)

type regexpCache struct {
	mu sync.Locker

	data map[string]*regexp.Regexp
}

var regexpCacheDefault = regexpCache{mu: &sync.Mutex{}, data: make(map[string]*regexp.Regexp)}

func Compile(pattern string) (*regexp.Regexp, error) {
	regexpCacheDefault.mu.Lock()
	defer regexpCacheDefault.mu.Unlock()

	if r, ok := regexpCacheDefault.data[pattern]; ok {
		return r, nil
	}

	if r, err := regexp.Compile(pattern); err != nil {
		return nil, err
	} else {
		regexpCacheDefault.data[pattern] = r
		return r, nil
	}
}
