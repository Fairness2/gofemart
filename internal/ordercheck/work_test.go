package ordercheck

import (
	"sync"
	"testing"
)

func TestPoolInWork(t *testing.T) {
	// test cases
	testCases := []struct {
		name   string
		input  string
		wantOk bool
		setup  func(*Pool)
	}{
		{
			name:   "already_in_work",
			input:  "1",
			wantOk: false,
			setup: func(p *Pool) {
				p.orderMap["1"] = &WorkedOrder{inWork: true}
			},
		},
		{
			name:   "order_not_found",
			input:  "1",
			wantOk: false,
			setup:  func(p *Pool) {},
		},
		{
			name:   "order_is_ready_to_work",
			input:  "1",
			wantOk: true,
			setup: func(p *Pool) {
				p.orderMap["1"] = &WorkedOrder{inWork: false}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a Pool
			p := &Pool{
				orderMap: make(map[string]*WorkedOrder),
				mutex:    sync.RWMutex{},
			}
			tc.setup(p)
			got, gotOk := p.poolInWork(tc.input)
			if gotOk != tc.wantOk {
				t.Errorf("gotOk = %v, want %v", gotOk, tc.wantOk)
				return
			}
			if tc.wantOk && got == nil {
				t.Errorf("got = nil, want not nil")
				return
			}
			if !tc.wantOk && got != nil {
				t.Errorf("got = %v, want nil", got)
				return
			}
			if tc.wantOk && got.inWork != true {
				t.Errorf("got.inWork = %v, want true", got.inWork)
			}
		})
	}
}

func TestRemoveFromWork(t *testing.T) {
	testCases := []struct {
		name       string
		order      string
		setup      func(*Pool)
		wantInWork bool
	}{
		{
			name:  "order_not_in_work",
			order: "1",
			setup: func(p *Pool) {
				p.orderMap["1"] = &WorkedOrder{inWork: false}
			},
			wantInWork: false,
		},
		{
			name:  "order_in_work",
			order: "2",
			setup: func(p *Pool) {
				p.orderMap["2"] = &WorkedOrder{inWork: true}
			},
			wantInWork: false,
		},
		{
			name:       "order_not_exist",
			order:      "3",
			setup:      func(p *Pool) {},
			wantInWork: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a Pool
			p := &Pool{
				orderMap: make(map[string]*WorkedOrder),
				mutex:    sync.RWMutex{},
			}
			tc.setup(p)
			p.removeFromWork(tc.order)
			got := p.checkInWork(tc.order)
			if got != tc.wantInWork {
				t.Errorf("inWork after removeFromWork = %v, want %v", got, tc.wantInWork)
			}
		})
	}
}

func TestCheckInWork(t *testing.T) {
	testCases := []struct {
		name       string
		order      string
		setup      func(*Pool)
		wantInWork bool
	}{
		{
			name:  "order_in_work",
			order: "1",
			setup: func(p *Pool) {
				p.orderMap["1"] = &WorkedOrder{inWork: true}
			},
			wantInWork: true,
		},
		{
			name:  "order_not_in_work",
			order: "2",
			setup: func(p *Pool) {
				p.orderMap["2"] = &WorkedOrder{inWork: false}
			},
			wantInWork: false,
		},
		{
			name:       "order_not_exist",
			order:      "3",
			setup:      func(p *Pool) {},
			wantInWork: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a Pool
			p := &Pool{
				orderMap: make(map[string]*WorkedOrder),
				mutex:    sync.RWMutex{},
			}
			tc.setup(p)
			got := p.checkInWork(tc.order)
			if got != tc.wantInWork {
				t.Errorf("inWork after checkInWork = %v, want %v", got, tc.wantInWork)
			}
		})
	}
}
