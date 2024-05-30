package compiler

// SymbolScope 作用域使用SymbolScope别名，SymbolScope本身不重要，主要是有唯一性；
// 使用String是为了方便调式。
type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL" //全局作用域
)

type Symbol struct {
	Name  string
	Scope SymbolScope // 作用域
	Index int         // 索引
}

type SymbolTable struct {
	store          map[string]Symbol // string为标识符，可以将标识符和Symbol相关联
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s}
}

// Define 将标识符作为参数
// 创建定义并返回Symbol
func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions, Scope: GlobalScope}
	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}

// Resolve 将一个定义的标识符交给符号表
// 返回与其相关的Define
func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]
	return obj, ok
}
