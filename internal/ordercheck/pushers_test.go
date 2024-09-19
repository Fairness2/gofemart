package ordercheck

import (
	"context"
	"errors"
	"gofemart/internal/models"
	"testing"
)

func TestPush(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*Pool)
		order   *models.Order
		want    bool
		wantErr error
	}{
		{
			name: "push_to_empty_pool",
			setup: func(p *Pool) {
				p.inChanel = make(chan string, 10)
			},
			order: &models.Order{
				Number: "1",
			},
			want: true,
		},
		{
			name: "push_to_full_pool",
			setup: func(p *Pool) {
				p.inChanel = make(chan string)
			},
			order: &models.Order{
				Number: "1",
			},
		},
		{
			name: "push_to_closed_pool",
			setup: func(p *Pool) {
				p.closeFlag.Store(true)
			},
			order: &models.Order{
				Number: "1",
			},
			wantErr: ErrorPoolClosed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pool{
				orderMap: make(map[string]*WorkedOrder),
				ctx:      context.Background(),
			}
			tt.setup(p)
			got, err := p.Push(tt.order)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("Error is not expected. error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("Error is not expected. error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("Result not expected. got = %v, want %v", got, tt.want)
			}
		})
	}
}
