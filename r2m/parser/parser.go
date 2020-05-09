package parser

import "github.com/wejick/phriction2Markdown/r2m/ast"

const maxHeadingLevel = 2

// Parser building AST from text input
type Parser struct {
	doc ast.Node
}

// Parse run parsing process and return AST
func (p *Parser) Parse(text []byte) (ast.Node, error) {
	err := p.parseBlock(text)

	return p.doc, err
}

// parseBlock walks through the text and build the block
// this is the main building block of the parser
func (p *Parser) parseBlock(text []byte) (err error) {

	for len(text) > 0 {

		// prefixed heading
		// = or ==
		if isPrefixHeading(text) {
			block, offset := processPrefixHeading(text)
			text = text[offset:]

			p.addBlock(block)
		}

		// horizontal line
		if p.isHorizontalRule(text) {
			block, offset := processHorizontalRule(text)
			text = text[offset:]

			p.addBlock(block)
		}
	}

	return
}

func (p *Parser) addBlock(block ast.Node) (err error) {
	ast.AppendChild(p.doc, block)
	return
}

// check whether it's HR
// ---
// ___
// ***
func (p *Parser) isHorizontalRule(text []byte) bool {
	// minimum 4 char to make it works
	// ***\n
	if len(text) < 4 {
		return false
	}

	if text[0] != '-' && text[0] != '_' && text[0] != '*' {
		return false
	}

	// when there's 3 or more consecutive cHR ended by \n, it's horizontal rule
	for i, c := range text {
		switch {
		case c == '\n':
			return i >= 3
		case c != text[0] && c != ' ':
			return false
		default:
		}
	}

	return false
}

func processHorizontalRule(text []byte) (block ast.Node, offset int) {
	block = &ast.HorizontalRule{}
	offset = skipCharUntil(text, '\n')

	return
}

// check whether it has header prefix
func isPrefixHeading(text []byte) bool {
	if text[0] == '=' {
		return true
	}
	return false
}

func processPrefixHeading(text []byte) (block ast.Node, offset int) {
	level := skipCharN(text, '=', maxHeadingLevel)
	end := skipCharUntil(text, '\n')

	if end > level && level != 0 {
		headingBlock := &ast.Heading{
			Level: level,
		}
		headingBlock.Content = text[level:end]

		return headingBlock, end
	}

	return nil, end
}

// advanced as long as char is found
func skipChar(text []byte, char byte) (n int) {
	for i := range text {
		if text[i] != char {
			return i
		}
	}

	return 0
}

// like skipChar but with maximum limit
func skipCharN(text []byte, char byte, maxN int) (n int) {
	for i := range text {
		if text[i] != char || i >= maxN {
			return i
		}
	}

	return 0
}

// advanced as long until char is found
func skipCharUntil(text []byte, char byte) (n int) {
	for i := range text {
		if text[i] == char {
			return i
		}
	}

	return len(text)
}
