// Most part of the code here are copied from https://github.com/gomarkdown/markdown
// with little modification as it fit

// Markdown is distributed under the Simplified BSD License:

// Copyright © 2011 Russ Ross
// Copyright © 2018 Krzysztof Kowalczyk
// Copyright © 2018 Authors
// All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:

// 1.  Redistributions of source code must retain the above copyright
//     notice, this list of conditions and the following disclaimer.

// 2.  Redistributions in binary form must reproduce the above
//     copyright notice, this list of conditions and the following
//     disclaimer in the documentation and/or other materials provided with
//     the distribution.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
// FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE
// COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN
// ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package ast

// Node AST node
type Node interface {
	GetParent() Node
	SetParent(newParent Node)
	GetChildren() []Node
	SetChildren(newChildren []Node)
}

// Container node main data structure
type Container struct {
	Parent   Node
	Children []Node

	Content []byte
}

// GetParent returns parent node
func (c *Container) GetParent() Node {
	return c.Parent
}

// SetParent sets the parent node
func (c *Container) SetParent(newParent Node) {
	c.Parent = newParent
}

// GetChildren returns children nodes
func (c *Container) GetChildren() []Node {
	return c.Children
}

// SetChildren sets children node
func (c *Container) SetChildren(newChildren []Node) {
	c.Children = newChildren
}

// Leaf is a type of node that cannot have children
type Leaf struct {
	Parent Node

	Literal []byte // Text contents of the leaf nodes
	Content []byte // Markdown content of the block nodes
}

// AsContainer returns nil
func (l *Leaf) AsContainer() *Container {
	return nil
}

// AsLeaf returns itself as *Leaf
func (l *Leaf) AsLeaf() *Leaf {
	return l
}

// GetParent returns parent node
func (l *Leaf) GetParent() Node {
	return l.Parent
}

// SetParent sets the parent nodd
func (l *Leaf) SetParent(newParent Node) {
	l.Parent = newParent
}

// GetChildren returns nil because Leaf cannot have children
func (l *Leaf) GetChildren() []Node {
	return nil
}

// SetChildren will panic becuase Leaf cannot have children
func (l *Leaf) SetChildren(newChildren []Node) {
	panic("leaf node cannot have children")
}

// Document representation of a remarkup document
// the root of the AST
type Document struct {
	Container
}

// Heading representation of heading node in ast
type Heading struct {
	Container

	Level int
}

// HorizontalRule represents markdown horizontal rule node
type HorizontalRule struct {
	Leaf
}

// AppendChild appends child to children of parent
// It panics if either node is nil.
func AppendChild(parent Node, child Node) {
	RemoveFromTree(child)
	child.SetParent(parent)
	newChildren := append(parent.GetChildren(), child)
	parent.SetChildren(newChildren)
}

// RemoveFromTree removes this node from tree
func RemoveFromTree(n Node) {
	if n.GetParent() == nil {
		return
	}
	// important: don't clear n.Children if n has no parent
	// we're called from AppendChild and that might happen on a node
	// that accumulated Children but hasn't been inserted into the tree
	n.SetChildren(nil)
	p := n.GetParent()
	newChildren := removeNodeFromArray(p.GetChildren(), n)
	if newChildren != nil {
		p.SetChildren(newChildren)
	}
}

func removeNodeFromArray(a []Node, node Node) []Node {
	n := len(a)
	for i := 0; i < n; i++ {
		if a[i] == node {
			return append(a[:i], a[i+1:]...)
		}
	}
	return nil
}
