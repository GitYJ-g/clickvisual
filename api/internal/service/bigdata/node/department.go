package node

import (
	"encoding/json"

	"github.com/clickvisual/clickvisual/api/internal/invoker"
	"github.com/clickvisual/clickvisual/api/pkg/model/db"
	"github.com/clickvisual/clickvisual/api/pkg/model/view"
)

const (
	OperatorRun int = iota
	OperatorStop
)

type node struct {
	n  *db.BigdataNode
	nc *db.BigdataNodeContent

	op int

	primaryDone   bool
	secondaryDone bool
	tertiaryDone  bool
}

type department interface {
	execute(*node) (view.RespRunNode, error)
	setNext(department)
}

func Operator(n *db.BigdataNode, nc *db.BigdataNodeContent, op int) (view.RespRunNode, error) {
	// Building chains of Responsibility
	t := &tertiary{}
	s := &secondary{next: t}
	p := &primary{next: s}
	res, err := p.execute(&node{
		n:             n,
		nc:            nc,
		op:            op,
		primaryDone:   false,
		secondaryDone: false,
		tertiaryDone:  false,
	})
	if err != nil {
		res.Message = err.Error()
	}
	// record execute result
	resBytes, _ := json.Marshal(res)
	ups := make(map[string]interface{}, 0)
	ups["result"] = string(resBytes)
	if op == OperatorRun {
		ups["previous_content"] = nc.Content
	}
	_ = db.NodeContentUpdate(invoker.Db, n.ID, ups)
	return res, err
}
