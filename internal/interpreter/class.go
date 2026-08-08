package interpreter

type LoxClass struct {
	name    string
	methods map[string]*LoxFunction
}

func NewLoxClass(name string, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{name: name, methods: methods}
}

func (c *LoxClass) FindMethod(name string) *LoxFunction {
	return c.methods[name]
}

func (c *LoxClass) Arity() int {
	if initializer := c.FindMethod("init"); initializer != nil {
		return initializer.Arity()
	}

	return 0
}

func (c *LoxClass) Call(i *Interpreter, arguments []any) (any, error) {
	instance := NewLoxInstance(c)

	if initializer := c.FindMethod("init"); initializer != nil {
		initializer.Bind(instance).Call(i, arguments)
	}

	return instance, nil
}

func (c *LoxClass) String() string {
	return c.name
}
