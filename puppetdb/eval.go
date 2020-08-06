package puppetdb

import (
	"fmt"
	"strconv"
)

func (bo BinOp) Eval() string {
	return fmt.Sprintf("%s %s %s", bo.LHS.Eval(), bo.Op, bo.RHS.Eval())
}

func (s Symbol) Eval() string {
	return s.Name
}

func (s String) Eval() string {
	return strconv.Quote(s.Val)
}

func (o Or) Eval() string {
	return fmt.Sprintf("(%s) or (%s)", o.LHS.Eval(), o.RHS.Eval())
}

func (a And) Eval() string {
	return fmt.Sprintf("(%s) and (%s)", a.LHS.Eval(), a.RHS.Eval())
}
