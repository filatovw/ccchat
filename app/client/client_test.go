package client

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAppRead(t *testing.T) {
	conf, _ := NewConf("", "myuser", "somehost", 0, false, "")

	check := func(input string, expected []string) {
		timeout := time.NewTimer(time.Second)
		r := strings.NewReader(input)
		w := new(bytes.Buffer)
		app := NewApp(conf, r, w)
		go app.read()
		actual := []string{}
		stopped := false
		for {
			select {
			case msg, ok := <-app.hub.Outbound:
				if ok {
					actual = append(actual, string(msg))
				} else {
					stopped = true
				}
			case <-timeout.C:
				app.hub.Close()
			}
			if stopped == true {
				break
			}
		}
		assert.Equal(t, expected, actual)
	}

	testData := `aaa::bbb
ccc::ddd`
	expected := []string{"aaa::bbb\r\n", "ccc::ddd\r\n"}
	check(testData, expected)
	testData = `aaa::bbb
		ccc`
	expected = []string{"aaa::bbb\r\n"}
	check(testData, expected)
	testData = `::bbb
		ccc::`
	expected = []string{"::bbb\r\n", "ccc::\r\n"}
	check(testData, expected)
}
