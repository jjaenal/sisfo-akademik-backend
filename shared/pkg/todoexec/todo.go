package todoexec

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Task struct {
	LineIndex int
	Raw       string
	Title     string
	Priority  int
	Completed bool
}

func Parse(path string) ([]Task, []string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	lines := splitLines(b)
	var tasks []Task
	for i, ln := range lines {
		t := strings.TrimSpace(ln)
		if strings.HasPrefix(t, "- [") {
			tt := Task{LineIndex: i, Raw: ln}
			tt.Completed = strings.HasPrefix(t, "- [x]") || strings.HasPrefix(t, "- [X]")
			tt.Title = strings.TrimSpace(stripCheckbox(t))
			tt.Priority = detectPriority(t)
			tasks = append(tasks, tt)
		}
	}
	return tasks, lines, nil
}

func splitLines(b []byte) []string {
	sc := bufio.NewScanner(bytes.NewReader(b))
	sc.Buffer(make([]byte, 1024), 1024*1024)
	var out []string
	for sc.Scan() {
		out = append(out, sc.Text())
	}
	return out
}

func stripCheckbox(s string) string {
	if idx := strings.Index(s, "]"); idx >= 0 {
		return strings.TrimSpace(s[idx+1:])
	}
	return s
}

func detectPriority(s string) int {
	if strings.Contains(s, "🔴") {
		return 3
	}
	if strings.Contains(s, "🟡") {
		return 2
	}
	if strings.Contains(s, "🟢") {
		return 1
	}
	return 0
}

func SortByPriority(tasks []Task) []Task {
	cp := make([]Task, 0, len(tasks))
	cp = append(cp, tasks...)
	sort.SliceStable(cp, func(i, j int) bool {
		if cp[i].Priority == cp[j].Priority {
			return cp[i].LineIndex < cp[j].LineIndex
		}
		return cp[i].Priority > cp[j].Priority
	})
	return cp
}

type ExecResult struct {
	Task       Task
	Executed   bool
	Succeeded  bool
	Error      string
}

func Execute(root string, t Task) ExecResult {
	r := ExecResult{Task: t}
	if t.Completed {
		return r
	}
	title := strings.ToLower(t.Title)
	if strings.Contains(title, "unit tests for auth") {
		out, err := runMake(root, "test-coverage")
		if err != nil {
			r.Executed = true
			r.Succeeded = false
			r.Error = err.Error()
			return r
		}
		ok := parseCoverage(out, "services/auth-service/internal/handler", 70)
		r.Executed = true
		r.Succeeded = ok
		if !ok {
			r.Error = "coverage threshold not met"
		}
		return r
	}
	if strings.Contains(title, "failed login tracking") {
		p := filepath.Join(root, "services/auth-service/internal/handler/auth.go")
		b, err := os.ReadFile(p)
		if err != nil {
			r.Executed = true
			r.Succeeded = false
			r.Error = err.Error()
			return r
		}
		s := string(b)
		ok := strings.Contains(s, "loginfail:") && strings.Contains(s, "lockout:")
		r.Executed = true
		r.Succeeded = ok
		if !ok {
			r.Error = "lockout logic not found"
		}
		return r
	}
	if strings.Contains(title, "implement forgot password") {
		hp := filepath.Join(root, "services/auth-service/internal/handler/auth.go")
		hb, herr := os.ReadFile(hp)
		ok := herr == nil
		if ok {
			hs := string(hb)
			ok = strings.Contains(hs, "forgotPassword") && strings.Contains(hs, `"/api/v1/auth/forgot-password"`)
		}
		tp := filepath.Join(root, "services/auth-service/internal/handler/auth_handler_test.go")
		tb, terr := os.ReadFile(tp)
		if terr == nil {
			ts := string(tb)
			ok = ok && strings.Contains(ts, `"/api/v1/auth/forgot-password"`)
		}
		r.Executed = true
		r.Succeeded = ok
		if !ok {
			r.Error = "forgot-password handler/tests not found"
		}
		return r
	}
	if strings.Contains(title, "implement reset password") {
		hp := filepath.Join(root, "services/auth-service/internal/handler/auth.go")
		hb, herr := os.ReadFile(hp)
		ok := herr == nil
		if ok {
			hs := string(hb)
			ok = strings.Contains(hs, "resetPassword") && strings.Contains(hs, `"/api/v1/auth/reset-password"`)
		}
		tp := filepath.Join(root, "services/auth-service/internal/handler/auth_handler_test.go")
		tb, terr := os.ReadFile(tp)
		if terr == nil {
			ts := string(tb)
			ok = ok && strings.Contains(ts, `"/api/v1/auth/reset-password"`)
		}
		r.Executed = true
		r.Succeeded = ok
		if !ok {
			r.Error = "reset-password handler/tests not found"
		}
		return r
	}
	if strings.Contains(title, "implement change password") {
		hp := filepath.Join(root, "services/auth-service/internal/handler/auth.go")
		hb, herr := os.ReadFile(hp)
		ok := herr == nil
		if ok {
			hs := string(hb)
			ok = strings.Contains(hs, "changePassword") && strings.Contains(hs, `"/api/v1/auth/change-password"`)
		}
		tp := filepath.Join(root, "services/auth-service/internal/handler/auth_handler_test.go")
		tb, terr := os.ReadFile(tp)
		if terr == nil {
			ts := string(tb)
			ok = ok && strings.Contains(ts, `"/api/v1/auth/change-password"`)
		}
		r.Executed = true
		r.Succeeded = ok
		if !ok {
			r.Error = "change-password handler/tests not found"
		}
		return r
	}
	if strings.Contains(title, "implement role handlers") {
		hp := filepath.Join(root, "services/auth-service/internal/handler/roles.go")
		hb, herr := os.ReadFile(hp)
		ok := herr == nil
		if ok {
			hs := string(hb)
			ok = strings.Contains(hs, `"/api/v1/users/:id/roles"`) && strings.Contains(hs, "RegisterProtected")
		}
		tp := filepath.Join(root, "services/auth-service/internal/handler/roles_handler_test.go")
		tb, terr := os.ReadFile(tp)
		if terr == nil {
			ts := string(tb)
			ok = ok && strings.Contains(ts, `"/api/v1/users/"`) && strings.Contains(ts, `"/roles"`)
		}
		r.Executed = true
		r.Succeeded = ok
		if !ok {
			r.Error = "roles handler/tests not found"
		}
		return r
	}
	if strings.Contains(title, "security headers") {
		p := filepath.Join(root, "services/auth-service/internal/middleware/security.go")
		b, err := os.ReadFile(p)
		if err != nil {
			r.Executed = true
			r.Succeeded = false
			r.Error = err.Error()
			return r
		}
		s := string(b)
		req := []string{
			"X-Content-Type-Options",
			"X-Frame-Options",
			"X-XSS-Protection",
			"Content-Security-Policy",
		}
		ok := true
		for _, k := range req {
			if !strings.Contains(s, k) {
				ok = false
				break
			}
		}
		r.Executed = true
		r.Succeeded = ok
		if !ok {
			r.Error = "required security headers missing"
		}
		return r
	}
	if strings.Contains(title, "setup gosec scanning") {
		p := filepath.Join(root, ".github", "workflows", "ci.yml")
		b, err := os.ReadFile(p)
		if err != nil {
			r.Executed = true
			r.Succeeded = false
			r.Error = err.Error()
			return r
		}
		s := string(b)
		ok := strings.Contains(s, "gosec") && strings.Contains(s, "Install gosec")
		r.Executed = true
		r.Succeeded = ok
		if !ok {
			r.Error = "gosec not configured in CI"
		}
		return r
	}
	if strings.Contains(title, "setup trivy scanning") {
		p := filepath.Join(root, ".github", "workflows", "ci.yml")
		b, err := os.ReadFile(p)
		if err != nil {
			r.Executed = true
			r.Succeeded = false
			r.Error = err.Error()
			return r
		}
		s := string(b)
		ok := strings.Contains(s, "aquasecurity/trivy-action")
		r.Executed = true
		r.Succeeded = ok
		if !ok {
			r.Error = "trivy not configured in CI"
		}
		return r
	}
	if strings.Contains(title, "unit tests for gateway") {
		// Require at least one _test.go file, then run go test
		hasTests := false
		rootDir := filepath.Join(root, "services", "api-gateway")
		_ = filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
			if err == nil && !d.IsDir() && strings.HasSuffix(d.Name(), "_test.go") {
				hasTests = true
			}
			return nil
		})
		if !hasTests {
			r.Executed = true
			r.Succeeded = false
			r.Error = "no gateway tests found"
			return r
		}
		cmd := exec.Command("go", "test", "./services/api-gateway/...")
		cmd.Dir = root
		out, err := cmd.CombinedOutput()
		if err != nil {
			r.Executed = true
			r.Succeeded = false
			r.Error = string(out)
			return r
		}
		r.Executed = true
		r.Succeeded = true
		return r
	}
	if strings.Contains(title, "connection pooling") {
		p := filepath.Join(root, "shared", "pkg", "database", "database.go")
		b, err := os.ReadFile(p)
		if err != nil {
			r.Executed = true
			r.Succeeded = false
			r.Error = err.Error()
			return r
		}
		s := string(b)
		ok := strings.Contains(s, "pgxpool.NewWithConfig")
		r.Executed = true
		r.Succeeded = ok
		if !ok {
			r.Error = "pgxpool usage not found"
		}
		return r
	}
	return r
}

func runMake(root string, target string) (string, error) {
	cmd := exec.Command("make", target)
	cmd.Dir = root
	b, err := cmd.CombinedOutput()
	if err != nil {
		return string(b), err
	}
	return string(b), nil
}

func parseCoverage(out string, pkg string, min int) bool {
	lines := strings.Split(out, "\n")
	for _, ln := range lines {
		if strings.Contains(ln, pkg) && strings.Contains(ln, "coverage:") {
			fs := strings.Fields(ln)
			for i, f := range fs {
				if f == "coverage:" && i+1 < len(fs) {
					val := strings.TrimSuffix(fs[i+1], "%")
					val = strings.TrimSuffix(val, "%")
					v, _ := strconv.ParseFloat(val, 64)
					return int(v) >= min
				}
			}
		}
	}
	return false
}

func UpdateStatuses(lines []string, execs []ExecResult) ([]byte, error) {
	m := map[int]bool{}
	for _, r := range execs {
		if r.Succeeded {
			m[r.Task.LineIndex] = true
		}
	}
	for i := range lines {
		if m[i] {
			t := lines[i]
			ts := strings.TrimSpace(t)
			if strings.HasPrefix(ts, "- [ ]") {
				lines[i] = strings.Replace(t, "- [ ]", "- [x]", 1)
			} else if strings.HasPrefix(ts, "- [X]") {
				lines[i] = strings.Replace(t, "- [X]", "- [x]", 1)
			}
		}
	}
	var buf bytes.Buffer
	for _, ln := range lines {
		buf.WriteString(ln)
		buf.WriteString("\n")
	}
	return buf.Bytes(), nil
}
