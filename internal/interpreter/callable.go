package interpreter

type LoxCallable interface {
	Arity() int
	Call(i *Interpreter, arguments []any) (any, error)
}
