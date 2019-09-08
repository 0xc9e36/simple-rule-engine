package utils

import "sync"

type Stack struct {
	top  *Node //栈顶指针
	size int   //栈大小
	lock *sync.RWMutex
}

//栈节点元素
type Node struct {
	value interface{}
	prev  *Node //前驱节点
}

//创建一个栈
func NewStack() *Stack {
	return &Stack{nil, 0, &sync.RWMutex{}}
}

//栈是否为空
func (s *Stack) IsEmpty() bool {
	return s.Size() == 0
}

//栈大小
func (s *Stack) Size() int {
	return s.size
}

//栈顶元素
func (s *Stack) Peek() interface{} {
	if s.Size() == 0 {
		return nil
	}
	return s.top.value
}

//出栈
func (s *Stack) Pop() interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.Size() == 0 {
		return nil
	}

	node := s.top
	s.top =  s.top.prev
	s.size--
	return node.value
}

//入栈
func (s *Stack) Push(v interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	node := &Node{
		value: v,
		prev: s.top,
	}
	s.top = node
	s.size++
}
