package config

import (
	"testing"
	"time"
)

func TestServerPort(t *testing.T) {
	tests := []struct {
		name string
		s    server
		want string
	}{
		{
			name: "valid port",
			s:    server{Port: 8080},
			want: ":8080",
		},
		{
			name: "invalid port low",
			s:    server{Port: 80},
			want: ":8080",
		},
		{
			name: "invalid port high",
			s:    server{Port: 70000},
			want: ":8080",
		},
		{
			name: "zero port",
			s:    server{Port: 0},
			want: ":8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.s.ServerPort(); got != tt.want {
				t.Errorf("ServerPort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataSourceMaxConns(t *testing.T) {
	tests := []struct {
		name string
		d    dataSource
		want int32
	}{
		{
			name: "valid max conns",
			d:    dataSource{MaxConnection: 10},
			want: 10,
		},
		{
			name: "invalid max conns",
			d:    dataSource{MaxConnection: 0},
			want: 5,
		},
		{
			name: "negative max conns",
			d:    dataSource{MaxConnection: -5},
			want: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.d.MaxConns(); got != tt.want {
				t.Errorf("MaxConns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataSourceMinConns(t *testing.T) {
	tests := []struct {
		name string
		d    dataSource
		want int32
	}{
		{
			name: "valid min conns",
			d:    dataSource{MinConnection: 2, MaxConnection: 10},
			want: 2,
		},
		{
			name: "min conns greater than max",
			d:    dataSource{MinConnection: 15, MaxConnection: 10},
			want: 1,
		},
		{
			name: "zero min conns",
			d:    dataSource{MinConnection: 0, MaxConnection: 10},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.d.MinConns(); got != tt.want {
				t.Errorf("MinConns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataSourceLifeTime(t *testing.T) {
	tests := []struct {
		name string
		d    dataSource
		want time.Duration
	}{
		{
			name: "valid lifetime",
			d:    dataSource{ConnectionLifeTime: "30m"},
			want: 30 * time.Minute,
		},
		{
			name: "invalid lifetime",
			d:    dataSource{ConnectionLifeTime: "invalid"},
			want: time.Hour,
		},
		{
			name: "empty lifetime",
			d:    dataSource{ConnectionLifeTime: ""},
			want: time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.d.LifeTime(); got != tt.want {
				t.Errorf("LifeTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataSourcePostgresConnString(t *testing.T) {
	tests := []struct {
		name string
		d    dataSource
		want string
	}{
		{
			name: "valid conn string",
			d: dataSource{
				Host:     "localhost",
				Port:     5432,
				Name:     "testdb",
				User:     "user",
				Password: "pass",
			},
			want: "postgres://user:pass@localhost:5432/testdb",
		},
		{
			name: "invalid port",
			d: dataSource{
				Host:     "localhost",
				Port:     80,
				Name:     "testdb",
				User:     "user",
				Password: "pass",
			},
			want: "postgres://user:pass@localhost:5432/testdb",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.d.PostgresConnString(); got != tt.want {
				t.Errorf("PostgresConnString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenSecretKey(t *testing.T) {
	tests := []struct {
		name string
		t    token
		want string
	}{
		{
			name: "valid secret",
			t:    token{Secret: "mysecret"},
			want: "mysecret",
		},
		{
			name: "empty secret",
			t:    token{Secret: ""},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.t.SecretKey(); got != tt.want {
				t.Errorf("SecretKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenRefreshTTL(t *testing.T) {
	tests := []struct {
		name string
		t    token
		want time.Duration
	}{
		{
			name: "valid refresh TTL",
			t:    token{RTokenTTL: "10m"},
			want: 10 * time.Minute,
		},
		{
			name: "invalid refresh TTL",
			t:    token{RTokenTTL: "invalid"},
			want: time.Hour,
		},
		{
			name: "empty refresh TTL",
			t:    token{RTokenTTL: ""},
			want: time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.t.RefreshTTL(); got != tt.want {
				t.Errorf("RefreshTTL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenAccessTTL(t *testing.T) {
	tests := []struct {
		name string
		t    token
		want time.Duration
	}{
		{
			name: "valid access TTL",
			t:    token{ATokenTTL: "30s"},
			want: 30 * time.Second,
		},
		{
			name: "invalid access TTL",
			t:    token{ATokenTTL: "invalid"},
			want: 5 * time.Minute,
		},
		{
			name: "empty access TTL",
			t:    token{ATokenTTL: ""},
			want: 5 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.t.AccessTTL(); got != tt.want {
				t.Errorf("AccessTTL() = %v, want %v", got, tt.want)
			}
		})
	}
}
