package swarm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRules(t *testing.T) {
	yamlRules := `
^(nginx)$:
  amd64: $1
  default: aiyara/$1:latest.$arch
`
	ruleMap := parseRules(yamlRules)
	assert.Equal(t, ruleMap[`^(nginx)$`]["amd64"], `$1`)
	assert.Equal(t, ruleMap[`^(nginx)$`]["default"], `aiyara/$1:latest.$arch`)
}

func TestParseRules_2(t *testing.T) {
	yamlRules := `
"^(nginx):(latest)$":
  amd64: $1
  default: aiyara/$1:latest.$arch
`
	ruleMap := parseRules(yamlRules)
	assert.Equal(t, ruleMap[`^(nginx):(latest)$`]["amd64"], `$1`)
	assert.Equal(t, ruleMap[`^(nginx):(latest)$`]["default"], `aiyara/$1:latest.$arch`)
}

func TestPullImageSimpleRewrite_01(t *testing.T) {
	yamlRules := `
^(nginx)$:
  amd64: $1_2
  default: aiyara/$1:latest.$arch
`
	ruleMap := parseRules(yamlRules)
	var labels map[string]string
	var result string

	labels = map[string]string{"arch": "amd64"}
	result, _ = rewrite("nginx", labels, ruleMap)
	assert.Equal(t, result, "nginx_2")
}

func TestPullImageSimpleRewriteReuseGroup(t *testing.T) {
	yamlRules := `
^(nginx)$:
  amd64: $1_$1_2
  default: aiyara/$1:latest.$arch
`
	ruleMap := parseRules(yamlRules)
	var labels map[string]string
	var result string

	labels = map[string]string{"arch": "amd64"}
	result, _ = rewrite("nginx", labels, ruleMap)
	assert.Equal(t, result, "nginx_nginx_2")
}

func TestPullImageSimpleRewrite_Rules_not_match(t *testing.T) {
	yamlRules := `
"^(nginx):(latest)$":
  amd64: $1_$2_3
  default: aiyara/$1:latest.$arch
`
	ruleMap := parseRules(yamlRules)
	var labels map[string]string
	var result string
	var err error

	labels = map[string]string{"arch": "amd64"}
	result, err = rewrite("nginx", labels, ruleMap)
	assert.Error(t, err)
	assert.Equal(t, result, "nginx")
}

func TestPullImageRewriteMultipleGroups(t *testing.T) {
	yamlRules := `
"^(nginx):(latest)$":
  amd64: $1_$2_3
  default: aiyara/$1:latest.$arch
`
	ruleMap := parseRules(yamlRules)
	var labels map[string]string
	var result string
	var err error

	labels = map[string]string{"arch": "amd64"}
	result, err = rewrite("nginx:latest", labels, ruleMap)
	assert.NoError(t, err)
	assert.Equal(t, result, "nginx_latest_3")
}

func TestPullImageRewriteDefaultCase(t *testing.T) {
	yamlRules := `
^(nginx)$:
  amd64: $1_2
  default: aiyara/$1:latest.$arch
`
	ruleMap := parseRules(yamlRules)
	var labels map[string]string
	var result string
	var err error

	labels = map[string]string{"arch": "amd64"}
	result, err = rewrite("nginx", labels, ruleMap)
	assert.NoError(t, err)
	assert.Equal(t, result, "nginx_2")

	labels = map[string]string{"arch": "arm"}
	result, err = rewrite("nginx", labels, ruleMap)
	assert.NoError(t, err)
	assert.Equal(t, result, "aiyara/nginx:latest.arm")

	labels = map[string]string{"arch": "386"}
	result, err = rewrite("nginx", labels, ruleMap)
	assert.NoError(t, err)
	assert.Equal(t, result, "aiyara/nginx:latest.386")

	labels = map[string]string{"arch": "arm64"}
	result, err = rewrite("nginx", labels, ruleMap)
	assert.NoError(t, err)
	assert.Equal(t, result, "aiyara/nginx:latest.arm64")

}
