package consensus

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

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

func newAcceptor(id int, learners []int) *Acceptor {
	acceptor := &Acceptor{
		id:       id, //接受者的id
		learners: learners,
	}
	acceptor.startServer()
	return acceptor
}

//准备阶段
func (a *Acceptor) Prepare(args *MsgArgs, reply *MsgReply) error {
	if args.Number > a.minProposal {
		a.minProposal = args.Number
		reply.Number = a.acceptedNumber
		reply.Value = a.acceptedValue
		reply.Ok = true
	} else {
		reply.Ok = false
	}
	return nil
}

//接受阶段
func (a *Acceptor) Accept(args *MsgArgs, reply *MsgReply) error {
	if args.Number >= a.minProposal {
		a.minProposal = args.Number
		a.acceptedNumber = args.Number
		a.acceptedValue = args.Value
		reply.Ok = true

		//后台执行将结果转发给学习者
		for _, lid := range a.learners {
			go func(learner int) {
				addr := fmt.Sprintf("127.0.0.1:%d", learner)
				args.From = a.id
				args.To = learner
				resp := new(MsgReply)
				ok := call(addr, "Learner.Learn", args, resp)
				if !ok {
					return
				}
			}(lid)
		}
	} else {
		reply.Ok = false
	}
	return nil
}

func (a *Acceptor) startServer() {
	server := rpc.NewServer()
	server.Register(a)
	addr := fmt.Sprintf(":%d", a.id)
	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatal("listetn is error", e)
	}
	a.lis = l

	go func() {
		for {
			accept, err := a.lis.Accept()
			if err != nil {
				continue
			}
			go server.ServeConn(accept)
		}
	}()
}

//释放连接
func (a *Acceptor) closeConnect() {
	a.lis.Close()
}
