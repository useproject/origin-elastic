// +build integration

package connection

import (
	"fmt"
	"os"
	"testing"

	"github.com/useproject/origin-elastic/beats/libbeat/tests/compose"
	mbtest "github.com/useproject/origin-elastic/beats/metricbeat/mb/testing"
)

func TestData(t *testing.T) {
	compose.EnsureUp(t, "rabbitmq")

	f := mbtest.NewEventsFetcher(t, getConfig())
	err := mbtest.WriteEvents(f, t)
	if err != nil {
		t.Fatal("write", err)
	}
}

func getConfig() map[string]interface{} {
	return map[string]interface{}{
		"module":     "rabbitmq",
		"metricsets": []string{"connection"},
		"hosts":      getTestRabbitMQHost(),
		"username":   getTestRabbitMQUsername(),
		"password":   getTestRabbitMQPassword(),
	}
}

const (
	rabbitmqDefaultHost     = "localhost"
	rabbitmqDefaultPort     = "15672"
	rabbitmqDefaultUsername = "guest"
	rabbitmqDefaultPassword = "guest"
)

func getTestRabbitMQHost() string {
	return fmt.Sprintf("%v:%v",
		getenv("RABBITMQ_HOST", rabbitmqDefaultHost),
		getenv("RABBITMQ_PORT", rabbitmqDefaultPort),
	)
}

func getTestRabbitMQUsername() string {
	return getenv("RABBITMQ_USERNAME", rabbitmqDefaultUsername)
}

func getTestRabbitMQPassword() string {
	return getenv("RABBITMQ_PASSWORD", rabbitmqDefaultPassword)
}

func getenv(name, defaultValue string) string {
	return strDefault(os.Getenv(name), defaultValue)
}

func strDefault(a, defaults string) string {
	if len(a) == 0 {
		return defaults
	}
	return a
}
