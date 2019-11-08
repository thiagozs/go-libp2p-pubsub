package core

import (
	"fmt"
	"strings"
	"sync"
)

type counters struct {
	cts map[string][]map[string]int
	wg  sync.Mutex
}

// NewCounters return a counters
func NewCounters() *counters {
	return &counters{
		cts: make(map[string][]map[string]int),
	}
}

// Add sum stats to hash map
func (c *counters) Add(msg string) {
	defer c.wg.Unlock()
	c.wg.Lock()

	params := strings.Split(msg, "|")

	if len(params) > 0 {
		// startup index keys
		if _, ok := c.cts[params[0]]; !ok {
			c.cts[params[0]] = make([]map[string]int, 0)
			c.cts[params[0]] = append(c.cts[params[0]], map[string]int{params[1]: 1})
			return
		}
		// find for key if exists
		for i, v := range c.cts[params[0]] {
			if _, ok := v[params[1]]; ok {
				c.cts[params[0]][i][params[1]] += 1 // sum in the same index
				return
			}
		}
		// not found, create one
		c.cts[params[0]] = append(c.cts[params[0]], map[string]int{params[1]: 1})
	}
}

// Reset flush all data
func (c *counters) Reset() {
	c.cts = make(map[string][]map[string]int)
}

// Show get all data and plot on the screen
func (c *counters) Show() {
	for k, v := range c.cts {
		fmt.Println("Message from :", k)
		for _, vv := range v {
			for slg, val := range vv {
				fmt.Printf("Slug '%s' - total %d\n", slg, val)
			}
		}
	}
}
