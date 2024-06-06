package myhttp

import (
	"errors"
	"strings"
)

var (
	ErrRouteNotFound      = errors.New("Could not find the asked route.")
	ErrDuplicatePathParam = errors.New("Cannot register two different path params under same prefix")
	ErrDuplicatePath      = errors.New("Cannot register identical paths")
)

type node struct {
	isPathParam bool
	segment     string
	handlerFunc HandlerFunc
	children    map[string]*node
}

func NewNode(segment string) *node {
	return &node{
		segment:     segment,
		handlerFunc: nil,
		children:    make(map[string]*node),
		isPathParam: false,
	}
}

type RouterTree struct {
	root *node
}

func NewRouterTree() RouterTree {
	return RouterTree{
		root: NewNode("NOOP"),
	}
}

func (rt RouterTree) Add(path string, f HandlerFunc) error {
	splitPath := strings.Split(path, "/")
	splitPath[0] = "/"
	length := len(splitPath)
	if splitPath[length-1] == "" {
		splitPath = splitPath[:length-1]
	}

	trav := rt.root
	for _, segment := range splitPath {
		isPathParam := segment[0] == ':'
		if isPathParam {
			segment = segment[1:]
		}
		node, found := trav.children[segment]
		if !found {
			//doing a sanity check to see if there are any other
			//path params registered under the same prefix

			for _, child := range trav.children {
				if child.isPathParam {
					return ErrDuplicatePathParam
				}
			}

			node = NewNode(segment)
			node.isPathParam = isPathParam
			trav.children[segment] = node
		}
		trav = node
	}

	if trav.handlerFunc != nil {
		return ErrDuplicatePath
	}

	trav.handlerFunc = f

	return nil
}

// here we are passing the real target string obtained from request
func (rt RouterTree) GetHandler(ctx *HttpContext) (HandlerFunc, error) {
	splitPath := strings.Split(ctx.request.RequestLine.Target, "/")
	splitPath[0] = "/"
	length := len(splitPath)
	if splitPath[length-1] == "" {
		splitPath = splitPath[:length-1]
	}

	trav := rt.root
	for _, segment := range splitPath {
		node, found := trav.children[segment]

		if !found {
			//it might be a path parameter, so let's see if children
			//has any path params
			hasChildWithPathParam := false
			for _, child := range trav.children {
				if child.isPathParam {
					//if we get a path param, we will update the
					//value received from the request in the context
					//for further access

					ctx.pathParams[child.segment] = segment
					trav = child
					hasChildWithPathParam = true
					break
				}
			}

			if !hasChildWithPathParam {
				return nil, ErrRouteNotFound
			}
			continue
		}
		trav = node
	}

	if trav.handlerFunc == nil {
		return nil, ErrRouteNotFound
	}

	return trav.handlerFunc, nil
}
