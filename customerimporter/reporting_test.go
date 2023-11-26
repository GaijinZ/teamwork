package customerimporter

import (
	"reflect"
	"teamwork/logger"
	"testing"
)

func TestReport_Sort(t *testing.T) {
	type fields struct {
		domainCustomerCount map[string]int
		RowsProcessed       int
		logger              logger.Logger
	}
	tests := []struct {
		name   string
		fields fields
		want   []DomainCounter
	}{
		{
			name: "Test Case 1",
			fields: fields{
				domainCustomerCount: map[string]int{
					"gmail.com":  3,
					"yandex.ru":  5,
					"devhub.com": 1,
				},
			},
			want: []DomainCounter{
				{Domain: "devhub.com", Users: 1},
				{Domain: "gmail.com", Users: 3},
				{Domain: "yandex.ru", Users: 5},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Report{
				domainCustomerCount: tt.fields.domainCustomerCount,
				RowsProcessed:       tt.fields.RowsProcessed,
				logger:              tt.fields.logger,
			}
			if got := r.Sort(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sort() = %v, want %v", got, tt.want)
			}
		})
	}
}
