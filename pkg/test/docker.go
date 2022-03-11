package test

import (
	"os/exec"
	"strings"
	"testing"
)

func RequireDockerServices(t *testing.T, services ...string) {
	out, err := exec.Command("docker", "ps", "--format", "{{.Names}}").Output()
	if err != nil {
		t.Fatalf("running docker: %s: %s", err, err.(*exec.ExitError).Stderr)
	}
	running := map[string]struct{}{}
	for _, line := range strings.Split(string(out), "\n") {
		if !strings.HasPrefix(line, "telehealth-") {
			continue
		}
		parts := strings.Split(line, "-")
		running[parts[1]] = struct{}{}
	}
	missing := false
	for _, s := range services {
		if _, ok := running[s]; !ok {
			t.Errorf("required service: %s", s)
			missing = true
		}
	}
	if missing {
		t.Fatalf("some required docker services are missing")
	}
}
