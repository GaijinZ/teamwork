package customerimporter

import (
	"reflect"
	"sync"
	"teamwork/logger"
	"testing"
)

func Test_extractDomain(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test Case 1",
			args: args{email: "jalvarezj@so-net.ne.jp"},
			want: "so-net.ne.jp",
		},

		{
			name: "Test Case 2",
			args: args{email: "jalvarezjso-net.ne.jp"},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractDomain(tt.args.email); got != tt.want {
				t.Errorf("extractDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainProcessor_checkHeaders(t *testing.T) {
	type fields struct {
		filePath  string
		log       logger.Logger
		WaitGroup sync.WaitGroup
		Mutex     sync.Mutex
	}
	tests := []struct {
		name    string
		fields  fields
		want    int
		want1   [][]string
		wantErr bool
	}{
		{
			name: "Test Case 1",
			fields: fields{
				WaitGroup: sync.WaitGroup{},
				Mutex:     sync.Mutex{},
				filePath:  "../data/test_data_customers.csv",
			},
			want: 2,
			want1: [][]string{
				{"first_name", "last_name", "email", "gender", "ip_address"},
				{"Mildred", "Hernandez", "mhernandez0@github.io", "Female", "38.194.51.128"},
				{"Bonnie", "Ortiz", "bortiz1@cyberchimps.com", "Female", "197.54.209.129"},
				{"Dennis", "Henry", "dhenry2@hubpages.com", "Male", "155.75.186.217"},
				{"Justin", "Hansen", "jhansen3@360.cn", "Male", "251.166.224.119"},
			},
			wantErr: false,
		},

		{
			name: "Test Case with Error",
			fields: fields{
				WaitGroup: sync.WaitGroup{},
				Mutex:     sync.Mutex{},
				filePath:  "./data/test_data_customers.csv",
			},
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := logger.NewLogger()
			d := &DomainProcessor{
				filePath:  tt.fields.filePath,
				log:       *log,
				WaitGroup: tt.fields.WaitGroup,
				Mutex:     tt.fields.Mutex,
			}
			got, got1, err := d.checkHeaders()
			if (err != nil) != tt.wantErr {
				t.Errorf("checkHeaders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkHeaders() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("checkHeaders() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_processRows(t *testing.T) {
	type args struct {
		emailIndex int
		row        []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test Case 1",
			args: args{
				emailIndex: 2,
				row: []string{
					"Mildred", "Hernandez", "mhernandez0@github.io", "Female", "38.194.51.128",
				},
			},
			want: "github.io",
		},

		{
			name: "Test Case 2",
			args: args{
				emailIndex: 2,
				row: []string{
					"Mildred", "Hernandez",
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processRows(tt.args.emailIndex, tt.args.row); got != tt.want {
				t.Errorf("processRows() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainProcessor_GetDomains(t *testing.T) {
	type domainCounter struct {
		Domain string
		Users  int
	}
	type fields struct {
		filePath  string
		log       logger.Logger
		WaitGroup sync.WaitGroup
		Mutex     sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
		want   []domainCounter
	}{
		{
			name: "Test Case 1",
			fields: fields{
				filePath:  "../data/test_data_customers.csv",
				log:       logger.Logger{},
				WaitGroup: sync.WaitGroup{},
				Mutex:     sync.Mutex{},
			},
			want: []domainCounter{
				{"360.cn", 1},
				{"cyberchimps.com", 1},
				{"github.io", 1},
				{"hubpages.com", 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := logger.NewLogger()
			d := &DomainProcessor{
				filePath:  tt.fields.filePath,
				log:       *log,
				WaitGroup: tt.fields.WaitGroup,
				Mutex:     tt.fields.Mutex,
			}
			d.GetDomains()
		})
	}
}
