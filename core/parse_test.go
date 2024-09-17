package core

import (
	"reflect"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   Jobs
		isErr  bool
	}{
		{
			name: "Single work, multiple tasks",
			config: Config{
				Settings: Settings{SleepRange: 5},
				Works: []Work{
					{
						Uri:  "http://localhost",
						Port: 8000,
						Tasks: []Task{
							{
								Path:   "/users/1",
								Method: "GET",
								Body:   nil,
							},
							{
								Path:   "/users/2",
								Method: "POST",
								Body:   map[string]string{"key2": "value2"},
							},
						},
					},
				},
			},
			want: Jobs{
				{
					Url:    "http://localhost:8000/users/1",
					Method: "GET",
					Body:   nil,
				},
				{
					Url:    "http://localhost:8000/users/2",
					Method: "POST",
					Body:   map[string]string{"key2": "value2"},
				},
			},
			isErr: false,
		},
		{
			name: "No tasks available",
			config: Config{
				Settings: Settings{SleepRange: 5},
				Works:    []Work{},
			},
			want:  nil,
			isErr: true,
		},
		{
			name: "Task with no port",
			config: Config{
				Settings: Settings{SleepRange: 5},
				Works: []Work{
					{
						Uri:  "http://localhost",
						Port: 0,
						Tasks: []Task{
							{
								Path:   "/users/1",
								Method: "GET",
								Body:   nil,
							},
						},
					},
				},
			},
			want:  nil,
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.config)
			got, err := p.Parse()

			if (err != nil) != tt.isErr {
				t.Errorf("Parse() error = %v, isErr %v", err, tt.isErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
