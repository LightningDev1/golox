package interpreter

import "time"

type ClockFn struct{}

func (c *ClockFn) Arity() int { return 0 }

func (c *ClockFn) Call(i *Interpreter, args []any) (any, error) {
	return float64(time.Now().UnixNano()) / 1e9, nil
}

func (c *ClockFn) String() string { return "<native fn>" }
