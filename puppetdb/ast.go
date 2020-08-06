package puppetdb

type (
	Eval interface {
		Eval() string
	}

	BinOp struct {
		Op  string
		LHS Eval
		RHS Eval
	}

	Symbol struct {
		Name string
	}

	String struct {
		Val string
	}

	Or struct {
		LHS Eval
		RHS Eval
	}

	And struct {
		LHS Eval
		RHS Eval
	}
)
