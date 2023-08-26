package consensus

import "net"

type Acceptor struct {
	lis net.Listener
	//服务器id
	id int
	//接受者承诺的提案编号 如果为0 则表示接受者没有收到过任何的Prepare消息
	minProposal int
	//接受者已接受的提案编号，如果为0  表示没有接受过提案
	acceptedNumber int
	//接受者已接受的提案值 如果没有 为nil
	acceptedValue interface{}

	//学习者id列表
	learners []int
}
