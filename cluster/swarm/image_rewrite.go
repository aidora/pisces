package swarm

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func homepath(p string) string {
	home := os.Getenv("HOME")
	if runtime.GOOS == "windows" {
		home = os.Getenv("USERPROFILE")
	}
	return filepath.Join(home, p)
}

func parseRules(yamlRules string) map[string]map[string]string {
	rules := make(map[string]map[string]string)
	err := yaml.Unmarshal([]byte(yamlRules), &rules)
	if err != nil {
		return nil
	}
	return rules
}

func rewrite(s string, labels map[string]string, rules map[string]map[string]string) (string, error) {
	arch := labels["arch"]
	for k, v := range rules {
		log.Infof("k = %s", k)
		re, err := regexp.Compile(k)
		if err != nil {
			return s, err
		}

		target := v[arch]
		if target == "" {
			target = v["default"]

			// no default rule defined
			if target == "" {
				continue
			}
		}

		// replace special vars, e.g., $arch ...
		r := strings.NewReplacer("$arch", arch)
		target = r.Replace(target)

		match := re.FindAllStringSubmatch(s, -1)

		// if not match, skip to the next rule
		if len(match) == 0 {
			continue
		}

		for i := 1; i < len(match[0]); i++ {
			target = strings.Replace(target, fmt.Sprintf("$%d", i), match[0][i], -1)
			log.Infof("match[0][i] = %s", match[0][i])
			log.Infof("target = %s", target)
		}
		return target, nil
	}

	return s, fmt.Errorf("Rules not match")
}
