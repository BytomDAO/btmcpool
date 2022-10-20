package vars

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVars(t *testing.T) {
	f := "/tmp/vars_test.yml"

	config := []byte("Hacker: true\nname: steve\nage: 31\nhobbies:\n- skateboarding\n- snowboarding\n- go\nduration: 3s\nfnum: 31.2")

	assert.NoError(t, ioutil.WriteFile(f, config, 0644))
	defer os.Remove(f)
	initImpl(f)

	assert.True(t, GetBool("Hacker", false))
	assert.Equal(t, "steve", GetString("name", ""))
	assert.Equal(t, "", GetString("name0", ""))
	assert.Equal(t, int64(31), GetInt64("age", 0))
	assert.Equal(t, int64(0), GetInt64("age0", 0))
	arr := GetStringSlice("hobbies", []string{})
	assert.Equal(t, 3, len(arr))
	assert.Equal(t, "skateboarding", arr[0])
	assert.Equal(t, "snowboarding", arr[1])
	assert.Equal(t, "go", arr[2])
	assert.Equal(t, 3*time.Second, GetDuration("duration", time.Second))
	assert.Equal(t, time.Second, GetDuration("duration0", time.Second))
	assert.Equal(t, 31.2, GetFloat64("fnum", 0.0))
	assert.Equal(t, 0.0, GetFloat64("fnum0", 0.0))
}

type worldConfig struct {
	Usa   *countryConfig
	China *countryConfig
}

type countryConfig struct {
	Id     int64
	Name   string
	States []string
}

func TestVarsFill(t *testing.T) {
	f := "/tmp/vars_test.yml"

	config := []byte(`
world:
  USA:
    id: 0
    name : USA
    states:
    - CA
    - WA
    - NV
    - OR
  CHINA:
    id: 1
    name : CHINA
    states:
    - ZJ
    - SH
    - BJ
`)

	assert.NoError(t, ioutil.WriteFile(f, config, 0644))
	defer os.Remove(f)
	initImpl(f)

	var usaConfig countryConfig

	assert.NoError(t, Fill("world.USA", &usaConfig))
	assert.Equal(t, usaConfig.Id, int64(0))
	assert.Equal(t, usaConfig.Name, "USA")
	assert.Equal(t, usaConfig.States, []string{"CA", "WA", "NV", "OR"})
}
