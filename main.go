package chaos

import (
	"encoding/json"
	"math/rand"
	"time"

	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/chaos", new(Chaos))
}

type Chaos struct{}

func randomChance(prob float64) bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64() < prob
}

func randomByte() byte {
	return byte(rand.Intn(256))
}

func (c *Chaos) CorruptBytes(data string, probability float64) string {
	if !randomChance(probability) {
		return data
	}

	bytes := []byte(data)

	if len(bytes) == 0 {
		return data
	}

	pos := rand.Intn(len(bytes))
	bytes[pos] = randomByte()

	return string(bytes)
}

func (c *Chaos) CorruptJSONField(data string, field string, probability float64) string {

	if !randomChance(probability) {
		return data
	}

	var obj map[string]interface{}

	err := json.Unmarshal([]byte(data), &obj)
	if err != nil {
		return data
	}

	if val, ok := obj[field]; ok {

		switch v := val.(type) {

		case string:
			obj[field] = v + string(randomByte())

		case float64:
			obj[field] = rand.Float64() * 999999

		case bool:
			obj[field] = !v

		default:
			obj[field] = nil
		}
	}

	res, err := json.Marshal(obj)
	if err != nil {
		return data
	}

	return string(res)
}

func (c *Chaos) DelayMessage(ms int, probability float64) {

	if !randomChance(probability) {
		return
	}

	delay := rand.Intn(ms)

	time.Sleep(time.Duration(delay) * time.Millisecond)
}

func (c *Chaos) TruncateMessage(data string, probability float64) string {

	if !randomChance(probability) {
		return data
	}

	bytes := []byte(data)

	if len(bytes) < 2 {
		return data
	}

	pos := rand.Intn(len(bytes))

	return string(bytes[:pos])
}

func (c *Chaos) DuplicateMessage(data string, probability float64) []string {

	if randomChance(probability) {
		return []string{data, data}
	}

	return []string{data}
}
