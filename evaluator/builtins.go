package evaluator

import "Monkey/object"

var builtins = map[string]*object.Builtin{
	"len": {func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1", len(args))
		}
		switch ret := args[0].(type) {
		case *object.String:
			return &object.Integer{int64(len(ret.Value))}
		case *object.Array:
			return &object.Integer{int64(len(ret.Elements))}
		default:
			return newError("argument to `len` not supported, got %s", ret.Type())
		}
	}},
	"first": {func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got =%d, want =1", len(args))
		}
		if args[0].Type() != object.ARRAY_OBJ {
			return newError("argument to `first` must be Array, got %s", args[0].Type())
		}

		arr := args[0].(*object.Array)
		if len(arr.Elements) > 0 {
			return arr.Elements[0]
		}
		return NULL
	}},
	"last": {func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got =%d, want =1", len(args))
		}
		if args[0].Type() != object.ARRAY_OBJ {
			return newError("argument to `first` must be Array, got %s", args[0].Type())
		}

		arr := args[0].(*object.Array)
		if len(arr.Elements) > 0 {
			return arr.Elements[len(arr.Elements)-1]
		}
		return NULL
	}},
	"rest": {func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got =%d, want =1", len(args))
		}
		if args[0].Type() != object.ARRAY_OBJ {
			return newError("argument to `first` must be Array, got %s", args[0].Type())
		}
		arr := args[0].(*object.Array)

		length := len(arr.Elements)
		if length > 0 {
			newElement := make([]object.Object, length-1)
			copy(newElement, arr.Elements[1:length])
			return &object.Array{Elements: newElement}
		}
		return NULL
	}},
	"push": {func(args ...object.Object) object.Object {
		if len(args) != 2 {
			return newError("wrong number of arguments. got =%d, want =2", len(args))
		}
		if args[0].Type() != object.ARRAY_OBJ {
			return newError("argument to `first` must be Array, got %s", args[0].Type())
		}
		arr := args[0].(*object.Array)

		length := len(arr.Elements)
		if length > 0 {
			newElement := make([]object.Object, length)
			copy(newElement, arr.Elements)
			newElement = append(newElement, args[1])
			return &object.Array{Elements: newElement}
		}
		return NULL
	}},
}
