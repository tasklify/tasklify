package auth

import (
	"tasklify/third_party/enum"
)

type Action enum.Member[string]

var (
	a            = enum.NewBuilder[string, Action]()
	ActionCreate = a.Add(Action{"c"})
	ActionRead   = a.Add(Action{"r"})
	ActionUpdate = a.Add(Action{"u"})
	ActionDelete = a.Add(Action{"d"})
	Actions      = a.Enum()
)
