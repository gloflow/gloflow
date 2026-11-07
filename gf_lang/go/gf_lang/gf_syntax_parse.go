/*
GloFlow application and media management/publishing platform
Copyright (C) 2023 Ivan Trajkovic

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
*/

package gf_lang

import (
	"fmt"
	"strconv"
	"strings"
)

// Indentation-based surface syntax (.gfl) that parses into the same AST
// that the JSON-like (.gf) programs produce. See docs/syntax.md.

//-------------------------------------------------

type gfSynNode struct {
	text     string
	lineInt  int
	children []*gfSynNode
}

type gfSynItem struct {
	bare  bool   // unquoted word
	isEq  bool   // standalone "="
	isPipe bool  // standalone "|"
	isKV  bool   // key=value
	key   string
	val   interface{}
}

var gfSynInfixOps = map[string]bool{
	"+": true, "-": true, "*": true, "/": true, "%": true,
	"==": true, "!=": true, "<": true, ">": true, "<=": true, ">=": true,
	"&&": true, "||": true,
}

//-------------------------------------------------

// ParseSyntax parses .gfl source into a program AST (list of top-level expressions).
func ParseSyntax(pSrcStr string) (GFexpr, error) {

	roots, err := gfSynBuildTree(pSrcStr)
	if err != nil {
		return nil, err
	}
	stmts, err := gfSynBody(roots)
	if err != nil {
		return nil, err
	}
	return CastToExpr(stmts), nil
}

//-------------------------------------------------
// TREE (indentation)

func gfSynBuildTree(pSrcStr string) ([]*gfSynNode, error) {

	linesLst := strings.Split(strings.ReplaceAll(pSrcStr, "\r\n", "\n"), "\n")

	root := &gfSynNode{}
	type frame struct {
		indentInt int
		node      *gfSynNode
	}
	stack := []frame{{-1, root}}
	var pending *gfSynNode // a line that ends with "|" continues on the next line

	for i := 0; i < len(linesLst); i++ {
		raw := linesLst[i]
		lineNo := i + 1

		// MULTILINE_STRINGS """...""" -> single quoted string
		if strings.Count(raw, `"""`) == 1 {
			var sb strings.Builder
			sb.WriteString(raw[:strings.Index(raw, `"""`)])
			rest := raw[strings.Index(raw, `"""`)+3:]
			bodyLst := []string{rest}
			closed := false
			for i++; i < len(linesLst); i++ {
				l := linesLst[i]
				if idx := strings.Index(l, `"""`); idx >= 0 {
					bodyLst = append(bodyLst, l[:idx])
					sb.WriteString(strconv.Quote(gfSynDedent(bodyLst)))
					sb.WriteString(l[idx+3:])
					closed = true
					break
				}
				bodyLst = append(bodyLst, l)
			}
			if !closed {
				return nil, fmt.Errorf("line %d: unterminated \"\"\" string", lineNo)
			}
			raw = sb.String()
		}

		text := gfSynStripComment(raw)
		if strings.TrimSpace(text) == "" {
			continue
		}

		if pending != nil {
			pending.text += " " + strings.TrimSpace(text)
			if !strings.HasSuffix(pending.text, "|") {
				pending = nil
			}
			continue
		}

		indent := 0
		for indent < len(text) && text[indent] == ' ' {
			indent++
		}
		if indent < len(text) && text[indent] == '\t' {
			return nil, fmt.Errorf("line %d: tabs are not allowed for indentation", lineNo)
		}

		for stack[len(stack)-1].indentInt >= indent {
			stack = stack[:len(stack)-1]
		}
		n := &gfSynNode{text: strings.TrimSpace(text), lineInt: lineNo}
		parent := stack[len(stack)-1].node
		parent.children = append(parent.children, n)
		stack = append(stack, frame{indent, n})
		if strings.HasSuffix(n.text, "|") {
			pending = n
		}
	}
	return root.children, nil
}

// removes common leading whitespace and the empty first/last lines.
func gfSynDedent(pLinesLst []string) string {
	if len(pLinesLst) > 0 && strings.TrimSpace(pLinesLst[0]) == "" {
		pLinesLst = pLinesLst[1:]
	}
	if len(pLinesLst) > 0 && strings.TrimSpace(pLinesLst[len(pLinesLst)-1]) == "" {
		pLinesLst = pLinesLst[:len(pLinesLst)-1]
	}
	min := -1
	for _, l := range pLinesLst {
		if strings.TrimSpace(l) == "" {
			continue
		}
		ind := len(l) - len(strings.TrimLeft(l, " "))
		if min < 0 || ind < min {
			min = ind
		}
	}
	outLst := []string{}
	for _, l := range pLinesLst {
		if len(l) >= min && min > 0 {
			l = l[min:]
		}
		outLst = append(outLst, l)
	}
	return strings.Join(outLst, "\n")
}

func gfSynStripComment(pStr string) string {
	inStr := false
	for i := 0; i < len(pStr); i++ {
		c := pStr[i]
		if inStr {
			if c == '\\' {
				i++
			} else if c == '"' {
				inStr = false
			}
			continue
		}
		if c == '"' {
			inStr = true
		} else if c == '/' && i+1 < len(pStr) && pStr[i+1] == '/' {
			return pStr[:i]
		}
	}
	return pStr
}

//-------------------------------------------------
// LINE ITEMS

type gfSynScanner struct {
	s string
	i int
}

func (sc *gfSynScanner) skipWS() {
	for sc.i < len(sc.s) && (sc.s[sc.i] == ' ' || sc.s[sc.i] == ',') {
		sc.i++
	}
}

// parses items until end of input or a closing ")".
func (sc *gfSynScanner) items(pCloser byte) ([]gfSynItem, error) {
	out := []gfSynItem{}
	for {
		sc.skipWS()
		if sc.i >= len(sc.s) {
			if pCloser != 0 {
				return nil, fmt.Errorf("missing '%c'", pCloser)
			}
			return out, nil
		}
		if sc.s[sc.i] == ')' || sc.s[sc.i] == '}' {
			if sc.s[sc.i] != pCloser {
				return nil, fmt.Errorf("unexpected '%c'", sc.s[sc.i])
			}
			sc.i++
			return out, nil
		}
		it, err := sc.item()
		if err != nil {
			return nil, err
		}
		out = append(out, it)
	}
}

func (sc *gfSynScanner) group() (interface{}, error) {
	sc.i++ // "("
	its, err := sc.items(')')
	if err != nil {
		return nil, err
	}
	// infix: (a op b)
	if len(its) == 3 && its[1].bare {
		if op, ok := its[1].val.(string); ok && gfSynInfixOps[op] {
			return []interface{}{op, its[0].val, its[2].val}, nil
		}
	}
	lst := []interface{}{}
	for _, it := range its {
		if it.isKV {
			lst = append(lst, it.key, it.val)
		} else {
			lst = append(lst, it.val)
		}
	}
	return lst, nil
}

// {a | b}: a nested statement (list of segments, see gfSynChain)
func (sc *gfSynScanner) brace() (interface{}, error) {
	sc.i++ // "{"
	its, err := sc.items('}')
	if err != nil {
		return nil, err
	}
	return gfSynChain(its, 0)
}

func (sc *gfSynScanner) quoted() (string, error) {
	start := sc.i
	sc.i++
	for sc.i < len(sc.s) {
		if sc.s[sc.i] == '\\' {
			sc.i += 2
			continue
		}
		if sc.s[sc.i] == '"' {
			sc.i++
			return strconv.Unquote(sc.s[start:sc.i])
		}
		sc.i++
	}
	return "", fmt.Errorf("unterminated string")
}

func (sc *gfSynScanner) value() (interface{}, bool, error) { // returns (val, bare)
	switch sc.s[sc.i] {
	case '(':
		v, err := sc.group()
		return v, false, err
	case '{':
		v, err := sc.brace()
		return v, false, err
	case '"':
		v, err := sc.quoted()
		return v, false, err
	}
	return nil, false, nil
}

func (sc *gfSynScanner) item() (gfSynItem, error) {

	if sc.s[sc.i] == '(' || sc.s[sc.i] == '{' || sc.s[sc.i] == '"' {
		v, _, err := sc.value()
		return gfSynItem{val: v}, err
	}

	start := sc.i
	for sc.i < len(sc.s) && !strings.ContainsRune(" ,(){}\"", rune(sc.s[sc.i])) {
		sc.i++
	}
	word := sc.s[start:sc.i]

	if word == "=" {
		return gfSynItem{isEq: true}, nil
	}
	if word == "|" {
		return gfSynItem{isPipe: true}, nil
	}

	// key=value (but not operators like == != <= >=)
	if idx := strings.Index(word, "="); idx > 0 && !gfSynInfixOps[word] {
		key, rest := word[:idx], word[idx+1:]
		if rest == "" {
			if sc.i < len(sc.s) && (sc.s[sc.i] == '(' || sc.s[sc.i] == '{' || sc.s[sc.i] == '"') {
				v, _, err := sc.value()
				return gfSynItem{isKV: true, key: key, val: v}, err
			}
			return gfSynItem{}, fmt.Errorf("missing value for '%s='", key)
		}
		v, err := sc.callOrAtom(rest)
		return gfSynItem{isKV: true, key: key, val: v}, err
	}

	v, err := sc.callOrAtom(word)
	_, isStr := v.(string)
	return gfSynItem{bare: isStr, val: v}, err
}

// name(args) -> [name, args...]; otherwise an atom.
func (sc *gfSynScanner) callOrAtom(pWordStr string) (interface{}, error) {
	if sc.i < len(sc.s) && sc.s[sc.i] == '(' {
		g, err := sc.group()
		if err != nil {
			return nil, err
		}
		return append([]interface{}{pWordStr}, g.([]interface{})...), nil
	}
	return gfSynAtom(pWordStr), nil
}

func gfSynAtom(pStr string) interface{} {
	if pStr == "true" {
		return true
	}
	if pStr == "false" {
		return false
	}
	if f, err := strconv.ParseFloat(pStr, 64); err == nil {
		return f
	}
	return pStr
}

//-------------------------------------------------
// STATEMENTS

func gfSynKeyword(pIts []gfSynItem) string {
	if len(pIts) > 0 && pIts[0].bare {
		s := pIts[0].val.(string)
		switch s {
		case "repeat", "rule", "if", "shader":
			return s
		}
	}
	return ""
}

func gfSynLine(pN *gfSynNode) ([]gfSynItem, bool, error) {
	text := pN.text
	block := false
	if strings.HasSuffix(text, ":") {
		block = true
		text = strings.TrimSpace(text[:len(text)-1])
	}
	sc := &gfSynScanner{s: text}
	its, err := sc.items(0)
	if err != nil {
		return nil, false, fmt.Errorf("line %d: %v", pN.lineInt, err)
	}
	return its, block, nil
}

// list of statements
func gfSynBody(pNodes []*gfSynNode) ([]interface{}, error) {
	out := []interface{}{}
	for _, n := range pNodes {
		s, err := gfSynStmt(n)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}

// splits generic items into transform modifiers and plain values
func gfSynSplit(pIts []gfSynItem) (mods []interface{}, words []interface{}) {
	for _, it := range pIts {
		if it.isKV {
			mods = append(mods, it.key, it.val)
		} else {
			words = append(words, it.val)
		}
	}
	return
}

func gfSynStmt(pN *gfSynNode) (interface{}, error) {

	its, block, err := gfSynLine(pN)
	if err != nil {
		return nil, err
	}
	errf := func(f string, a ...interface{}) error {
		return fmt.Errorf("line %d: %s", pN.lineInt, fmt.Sprintf(f, a...))
	}

	for _, it := range its {
		if it.isPipe {
			if block {
				return nil, errf("unexpected block")
			}
			return gfSynChain(its, pN.lineInt)
		}
	}

	switch gfSynKeyword(its) {

	case "repeat":
		if !block || len(its) != 2 {
			return nil, errf("expected 'repeat <n>:' with an indented block")
		}
		scope, err := gfSynScope(pN.children)
		if err != nil {
			return nil, err
		}
		return []interface{}{"*", its[1].val, scope}, nil

	case "rule":
		if !block || len(its) < 2 {
			return nil, errf("expected 'rule <name> [max=N]:' with an indented block")
		}
		body, err := gfSynBody(pN.children)
		if err != nil {
			return nil, err
		}
		mods := []interface{}{}
		for _, it := range its[2:] {
			if !it.isKV {
				return nil, errf("rule modifiers must be key=value")
			}
			key := it.key
			if key == "max" {
				key = "iters_max"
			}
			mods = append(mods, []interface{}{key, it.val})
		}
		if len(mods) == 0 {
			return []interface{}{"rule", its[1].val, body}, nil
		}
		return []interface{}{"rule", its[1].val, mods, body}, nil

	case "if":
		if !block {
			return nil, errf("expected 'if <cond>:' with an indented block")
		}
		var cond interface{}
		rest := its[1:]
		if len(rest) == 3 && rest[1].bare && gfSynInfixOps[rest[1].val.(string)] {
			cond = []interface{}{rest[1].val, rest[0].val, rest[2].val}
		} else if len(rest) == 1 {
			cond = rest[0].val
		} else {
			return nil, errf("condition must be '<a> <op> <b>' or '(...)'")
		}
		body, err := gfSynBody(pN.children)
		if err != nil {
			return nil, err
		}
		return []interface{}{"if", cond, body}, nil

	case "shader":
		if !block || len(its) != 2 {
			return nil, errf("expected 'shader <name>:' with an indented block")
		}
		out := []interface{}{"shader", its[1].val}
		uniforms := []interface{}{}
		uniformsAt := -1
		for _, c := range pN.children {
			cits, cblock, err := gfSynLine(c)
			if err != nil {
				return nil, err
			}
			if cblock {
				return nil, fmt.Errorf("line %d: blocks are not allowed in shaders", c.lineInt)
			}
			_, words := gfSynSplit(cits)
			if len(cits) > 0 && cits[0].bare && cits[0].val == "uniform" {
				if uniformsAt < 0 {
					uniformsAt = len(out)
					out = append(out, nil)
				}
				uniforms = append(uniforms, words[1:])
			} else {
				out = append(out, words)
			}
		}
		if uniformsAt >= 0 {
			out[uniformsAt] = []interface{}{"uniforms", uniforms}
		}
		return out, nil
	}

	if block {
		return nil, errf("unexpected block")
	}

	// SETTER: name = value
	if len(its) >= 3 && its[0].bare && its[1].isEq {
		return gfSynSetter(its), nil
	}

	for _, it := range its {
		if it.isEq {
			return nil, errf("unexpected '='")
		}
	}

	// a line with a single quoted string, number or bool is that atom itself (not a list),
	// ex. `"R1"` -> "R1", while the bare word `R1` -> ["R1"]
	if len(its) == 1 && !its[0].bare && !its[0].isKV {
		if _, isLst := its[0].val.([]interface{}); !isLst {
			return its[0].val, nil
		}
	}

	mods, words := gfSynSplit(its)
	if mods == nil {
		return words, nil
	}
	return append([]interface{}{mods}, words...), nil
}

func gfSynSetter(pIts []gfSynItem) interface{} {
	var v interface{}
	if len(pIts) == 3 {
		v = pIts[2].val
	} else {
		vs := []interface{}{}
		for _, it := range pIts[2:] {
			vs = append(vs, it.val)
		}
		v = vs
	}
	return []interface{}{"set", pIts[0].val, v}
}

// `a | b | c`: a statement made of segments. each segment gives the same elements
// as a line in a repeat scope (transforms + geometry, setters, lists).
func gfSynChain(pIts []gfSynItem, pLineInt int) (interface{}, error) {

	out := []interface{}{}
	seg := []gfSynItem{}

	flush := func() error {
		if len(seg) == 0 {
			return nil
		}
		defer func() { seg = nil }()
		if gfSynKeyword(seg) != "" {
			return fmt.Errorf("line %d: keyword blocks can not be used in a '|' chain", pLineInt)
		}
		if len(seg) >= 3 && seg[0].bare && seg[1].isEq {
			out = append(out, gfSynSetter(seg))
			return nil
		}
		for _, it := range seg {
			if it.isEq {
				return fmt.Errorf("line %d: unexpected '='", pLineInt)
			}
		}
		mods, words := gfSynSplit(seg)
		if mods != nil {
			out = append(out, mods)
			out = append(out, words...)
		} else if len(words) == 1 {
			out = append(out, words[0])
		} else {
			out = append(out, words)
		}
		return nil
	}

	for _, it := range pIts {
		if it.isPipe {
			if err := flush(); err != nil {
				return nil, err
			}
		} else {
			seg = append(seg, it)
		}
	}
	if err := flush(); err != nil {
		return nil, err
	}
	return out, nil
}

// a block body used as a single transform scope: transforms and geometry are
// flattened into one list, e.g. ["y",2] "cube"
func gfSynScope(pNodes []*gfSynNode) ([]interface{}, error) {
	out := []interface{}{}
	for _, n := range pNodes {
		its, block, err := gfSynLine(n)
		if err != nil {
			return nil, err
		}
		isGeneric := gfSynKeyword(its) == "" && !block
		if isGeneric {
			for _, it := range its {
				if it.isEq || it.isPipe {
					isGeneric = false
				}
			}
		}
		if !isGeneric {
			s, err := gfSynStmt(n)
			if err != nil {
				return nil, err
			}
			out = append(out, s)
			continue
		}
		mods, words := gfSynSplit(its)
		if mods != nil {
			out = append(out, mods)
			out = append(out, words...)
		} else if len(words) == 1 {
			out = append(out, words[0])
		} else {
			out = append(out, words)
		}
	}
	return out, nil
}
