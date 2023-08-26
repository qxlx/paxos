package consensus

import (
	"net"
)

type Learner struct {
	lis net.Listener
	//学习者id
	id int

	acceptedMsg map[int]MsgArgs
}
