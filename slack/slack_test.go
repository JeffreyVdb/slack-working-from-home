package slack

import (
	"testing"
	"time"
)

func TestTimeoutOption(t *testing.T) {
	timeoutOption := Timeout(76)
	client, _ := NewClient("foobar", timeoutOption)
	if client.timeout != time.Second*76 {
		t.Errorf("timeout on client is expected to be 76 seconds, but it is %d seconds", client.timeout/time.Second)
	}
}
