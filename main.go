package chaos

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"go.k6.io/k6/js/modules"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	modules.Register("k6/x/chaos", new(Chaos))
}

type Chaos struct{}

func randomChance(prob float64) bool {
	if prob <= 0.0 {
		return false
	}
	if prob >= 1.0 {
		return true
	}
	return rand.Float64() < prob
}

func randomByte() byte {
	return byte(rand.Intn(256))
}

func (c *Chaos) CorruptBytes(data string, probability float64, mode string) string {
	if !randomChance(probability) {
		return data
	}

	b := []byte(data)
	if len(b) == 0 {
		return data
	}

	pos := rand.Intn(len(b))
	switch mode {
	case "zero":
		b[pos] = 0
	case "random":
		b[pos] = randomByte()
	case "bitflip":
		bit := uint(1 << uint(rand.Intn(8)))
		b[pos] = b[pos] ^ byte(bit)
	default:
		b[pos] = b[pos] ^ 0xFF
	}

	return string(b)
}

func setJSONPath(m map[string]interface{}, path []string, newVal interface{}) error {
	if len(path) == 0 {
		return errors.New("empty path")
	}
	for i := 0; i < len(path)-1; i++ {
		p := path[i]
		if next, ok := m[p]; ok {
			if asMap, ok2 := next.(map[string]interface{}); ok2 {
				m = asMap
			} else {
				return fmt.Errorf("path %s is not an object", strings.Join(path[:i+1], "."))
			}
		} else {
			newMap := make(map[string]interface{})
			m[p] = newMap
			m = newMap
		}
	}
	m[path[len(path)-1]] = newVal
	return nil
}

func (c *Chaos) CorruptJSONField(data string, field string, probability float64) string {
	if !randomChance(probability) {
		return data
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(data), &obj); err != nil {
		return data
	}

	parts := strings.Split(field, ".")

	cur := obj
	for i, p := range parts {
		if i == len(parts)-1 {
			if val, ok := cur[p]; ok {
				switch v := val.(type) {
				case string:
					cur[p] = v + string(randomByte())
				case float64:
					cur[p] = rand.Float64() * 999999
				case bool:
					cur[p] = !v
				default:
					cur[p] = nil
				}
			}
		} else {
			if next, ok := cur[p]; ok {
				if asMap, ok2 := next.(map[string]interface{}); ok2 {
					cur = asMap
				} else {
					setJSONPath(obj, parts, nil)
					break
				}
			} else {
				newMap := make(map[string]interface{})
				cur[p] = newMap
				cur = newMap
			}
		}
	}

	res, err := json.Marshal(obj)
	if err != nil {
		return data
	}
	return string(res)
}
