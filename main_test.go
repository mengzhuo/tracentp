package main

import (
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"text/template"

	"github.com/beevik/ntp"
)

func TestPrettyRsp(t *testing.T) {
	// Test PrettyRsp struct creation
	response := &ntp.Response{
		Stratum:      2,
		Time:         time.Now(),
		ClockOffset:  time.Duration(0.001 * float64(time.Second)),
		RTT:          time.Duration(0.005 * float64(time.Second)),
		RootDistance: time.Duration(0.002 * float64(time.Second)),
		ReferenceID:  0x474E5053, // GPS
	}

	prettyRsp := &PrettyRsp{
		Response: *response,
		Address:  "pool.ntp.org",
		Seq:      1,
	}

	if prettyRsp.Address != "pool.ntp.org" {
		t.Errorf("Expected address 'pool.ntp.org', got '%s'", prettyRsp.Address)
	}

	if prettyRsp.Seq != 1 {
		t.Errorf("Expected seq 1, got %d", prettyRsp.Seq)
	}

	if prettyRsp.Stratum != 2 {
		t.Errorf("Expected stratum 2, got %d", prettyRsp.Stratum)
	}
}

func TestArgsStruct(t *testing.T) {
	// Test args struct parsing
	testArgs := args{
		Address: "pool.ntp.org",
		Timeout: 5 * time.Second,
		Port:    123,
		Count:   10,
	}

	if testArgs.Address != "pool.ntp.org" {
		t.Errorf("Expected address 'pool.ntp.org', got '%s'", testArgs.Address)
	}

	if testArgs.Timeout != 5*time.Second {
		t.Errorf("Expected timeout 5s, got %v", testArgs.Timeout)
	}

	if testArgs.Port != 123 {
		t.Errorf("Expected port 123, got %d", testArgs.Port)
	}

	if testArgs.Count != 10 {
		t.Errorf("Expected count 10, got %d", testArgs.Count)
	}
}

func TestVersion(t *testing.T) {
	// Test version function
	var testArgs args
	version := testArgs.Version()
	if version != "dev" {
		t.Errorf("Expected version 'dev', got '%s'", version)
	}
}

func TestAddressFormatting(t *testing.T) {
	// Test address formatting logic
	testCases := []struct {
		address  string
		port     int
		expected string
	}{
		{"pool.ntp.org", 123, "pool.ntp.org"},
		{"time.google.com", 123, "time.google.com"},
		{"time.nist.gov", 456, "time.nist.gov:456"},
		{"192.168.1.1", 789, "192.168.1.1:789"},
	}

	for _, tc := range testCases {
		addr := tc.address
		if tc.port != 123 {
			addr = fmt.Sprintf("%s:%d", tc.address, tc.port)
		}
		if addr != tc.expected {
			t.Errorf("For address=%s, port=%d: expected '%s', got '%s'",
				tc.address, tc.port, tc.expected, addr)
		}
	}
}

func TestTemplateParsing(t *testing.T) {
	// Test template parsing
	tmpl := template.Must(template.New("").Parse(defaultTemplate))
	if tmpl == nil {
		t.Error("Failed to parse default template")
	}

	// Test template execution
	response := &ntp.Response{
		Stratum:      2,
		ClockOffset:  time.Duration(0.001 * float64(time.Second)),
		RTT:          time.Duration(0.005 * float64(time.Second)),
		RootDistance: time.Duration(0.002 * float64(time.Second)),
	}

	prettyRsp := &PrettyRsp{
		Response: *response,
		Address:  "pool.ntp.org",
		Seq:      1,
	}

	var output strings.Builder
	err := tmpl.Execute(&output, prettyRsp)
	if err != nil {
		t.Errorf("Template execution failed: %v", err)
	}

	result := output.String()
	if !strings.Contains(result, "pool.ntp.org") {
		t.Errorf("Expected output to contain 'pool.ntp.org', got: %s", result)
	}

	if !strings.Contains(result, "seq=1") {
		t.Errorf("Expected output to contain 'seq=1', got: %s", result)
	}
}

func TestIPAddressValidation(t *testing.T) {
	// Test IP address validation
	testCases := []struct {
		address string
		isValid bool
	}{
		{"192.168.1.1", true},
		{"10.0.0.1", true},
		{"172.16.0.1", true},
		{"8.8.8.8", true},
		{"invalid-ip", false},
		{"256.256.256.256", false},
		{"192.168.1", false},
		{"", false},
	}

	for _, tc := range testCases {
		ip := net.ParseIP(tc.address)
		isValid := ip != nil
		if isValid != tc.isValid {
			t.Errorf("For address '%s': expected valid=%t, got valid=%t",
				tc.address, tc.isValid, isValid)
		}
	}
}

func TestDefaultTemplate(t *testing.T) {
	// Test default template format
	expectedFields := []string{
		"{{if .Validate}}",
		"{{.Validate}}",
		"{{else}}",
		"OK",
		"{{end}}",
		"from",
		"{{.Address}}",
		"seq={{.Seq}}",
		"stratum={{.Stratum}}",
		"offset={{.ClockOffset}}",
		"distance={{.RootDistance}}",
		"RTT={{.RTT}}",
		"ref={{.ReferenceString}}",
	}

	for _, field := range expectedFields {
		if !strings.Contains(defaultTemplate, field) {
			t.Errorf("Default template missing field: %s", field)
		}
	}
}

// Mock NTP server for testing
type mockNTPServer struct {
	responses []*ntp.Response
	errors    []error
	current   int
}

func (m *mockNTPServer) QueryWithOptions(addr string, opts ntp.QueryOptions) (*ntp.Response, error) {
	if m.current >= len(m.responses) {
		return nil, fmt.Errorf("no more responses")
	}

	response := m.responses[m.current]
	err := m.errors[m.current]
	m.current++

	if err != nil {
		return nil, err
	}
	return response, nil
}

func TestMockNTPServer(t *testing.T) {
	// Test mock NTP server
	mock := &mockNTPServer{
		responses: []*ntp.Response{
			{
				Stratum:      2,
				ClockOffset:  time.Duration(0.001 * float64(time.Second)),
				RTT:          time.Duration(0.005 * float64(time.Second)),
				RootDistance: time.Duration(0.002 * float64(time.Second)),
			},
			{
				Stratum:      1,
				ClockOffset:  time.Duration(0.0005 * float64(time.Second)),
				RTT:          time.Duration(0.003 * float64(time.Second)),
				RootDistance: time.Duration(0.001 * float64(time.Second)),
			},
		},
		errors: []error{nil, nil},
	}

	// Test first response
	response, err := mock.QueryWithOptions("pool.ntp.org", ntp.QueryOptions{})
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if response.Stratum != 2 {
		t.Errorf("Expected stratum 2, got %d", response.Stratum)
	}

	// Test second response
	response, err = mock.QueryWithOptions("time.google.com", ntp.QueryOptions{})
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if response.Stratum != 1 {
		t.Errorf("Expected stratum 1, got %d", response.Stratum)
	}

	// Test no more responses
	_, err = mock.QueryWithOptions("gps", ntp.QueryOptions{})
	if err == nil {
		t.Error("Expected error for no more responses")
	}
}

func TestErrorHandling(t *testing.T) {
	// Test error handling scenarios
	mock := &mockNTPServer{
		responses: []*ntp.Response{nil},
		errors:    []error{fmt.Errorf("network timeout")},
	}

	_, err := mock.QueryWithOptions("invalid.server", ntp.QueryOptions{})
	if err == nil {
		t.Error("Expected error for network timeout")
	}
	if !strings.Contains(err.Error(), "network timeout") {
		t.Errorf("Expected 'network timeout' error, got: %v", err)
	}
}

// Benchmark tests
func BenchmarkNTPServerQuery(b *testing.B) {
	// Benchmark NTP server query (mock)
	mock := &mockNTPServer{
		responses: make([]*ntp.Response, b.N),
		errors:    make([]error, b.N),
	}

	for i := 0; i < b.N; i++ {
		mock.responses[i] = &ntp.Response{
			Stratum:      2,
			ClockOffset:  time.Duration(0.001 * float64(time.Second)),
			RTT:          time.Duration(0.005 * float64(time.Second)),
			RootDistance: time.Duration(0.002 * float64(time.Second)),
		}
		mock.errors[i] = nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := mock.QueryWithOptions("pool.ntp.org", ntp.QueryOptions{})
		if err != nil {
			b.Errorf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkTemplateExecution(b *testing.B) {
	// Benchmark template execution
	tmpl := template.Must(template.New("").Parse(defaultTemplate))
	response := &ntp.Response{
		Stratum:      2,
		ClockOffset:  time.Duration(0.001 * float64(time.Second)),
		RTT:          time.Duration(0.005 * float64(time.Second)),
		RootDistance: time.Duration(0.002 * float64(time.Second)),
	}
	prettyRsp := &PrettyRsp{
		Response: *response,
		Address:  "pool.ntp.org",
		Seq:      1,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var output strings.Builder
		err := tmpl.Execute(&output, prettyRsp)
		if err != nil {
			b.Errorf("Template execution failed: %v", err)
		}
	}
}
