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
	"regexp"
	"encoding/json"
)

// AST -> .gfl text (the inverse of ParseSyntax).
// FormatSyntax re-parses its own output and fails if it does not give back the same AST,
// so a successful result is always an exact translation.

//-------------------------------------------------

var gfSynIdentRe = regexp.MustCompile(`^[A-Za-z_$][A-Za-z0-9_$-]*$`)
var gfSynCallRe  = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// keys that are written as `key=value` transforms
var gfSynTransformKeys = func() map[string]bool {
	m := map[string]bool{}
	for _, k := range getSymbolsAndConstants().PredefinedPropertiesLst {
		m[k] = true
	}
	return m
}()

type gfSynPrinter struct {
	sb  strings.Builder
	err error
}

//-------------------------------------------------

func FormatSyntax(pProgramAST GFexpr) (string, error) {

	// normalize all nested list types to []interface{}
	b, err := json.Marshal(pProgramAST)
	if err != nil {
		return "", err
	}
	var prog []interface{}
	if err := json.Unmarshal(b, &prog); err != nil {
		return "", err
	}

	p := &gfSynPrinter{}
	p.stmts(prog, 0)
	if p.err != nil {
		return "", p.err
	}
	out := p.sb.String()

	// VERIFY
	back, err := ParseSyntax(out)
	if err != nil {
		return "", fmt.Errorf("formatted output does not parse: %v", err)
	}
	b2, _ := json.Marshal(back)
	if string(b) != string(b2) {
		return "", fmt.Errorf("formatted output does not round-trip:\n  want %s\n  got  %s", b, b2)
	}
	return out, nil
}

//-------------------------------------------------

func (p *gfSynPrinter) fail(pFormatStr string, pArgs ...interface{}) {
	if p.err == nil {
		p.err = fmt.Errorf(pFormatStr, pArgs...)
	}
}

func (p *gfSynPrinter) line(pIndent int, pStr string) {
	p.sb.WriteString(strings.Repeat("    ", pIndent))
	p.sb.WriteString(pStr)
	p.sb.WriteString("\n")
}

func (p *gfSynPrinter) stmts(pLst []interface{}, pIndent int) {
	for _, e := range pLst {
		p.stmt(e, pIndent)
	}
}

//-------------------------------------------------
// STATEMENTS

func (p *gfSynPrinter) stmt(pE interface{}, pIndent int) {

	lst, isLst := pE.([]interface{})
	if !isLst {
		// atom statement: strings are always quoted (a bare word is a list `["word"]`)
		if s, isStr := pE.(string); isStr {
			p.line(pIndent, strconv.Quote(s))
		} else {
			p.line(pIndent, p.tok(pE))
		}
		return
	}
	if p.stmtSpecial(lst, pIndent) {
		return
	}
	if len(lst) == 0 {
		p.fail("empty list statement")
		return
	}

	// [list, ...]: a chain of segments "a | b | c"
	if _, ok := lst[0].([]interface{}); ok {
		p.chain(lst, pIndent)
		return
	}

	// flat list
	toks := []string{}
	for _, v := range lst {
		toks = append(toks, p.tok(v))
	}
	p.line(pIndent, strings.Join(toks, " "))
}

// statement made of segments (transforms+geometry, setters, lists) joined with "|".
// a long chain continues on the next lines (a line ending with "|" continues).
func (p *gfSynPrinter) chain(pLst []interface{}, pIndent int) {

	segs := []string{}
	for i := 0; i < len(pLst); {
		var s string
		s, i = p.segElem(pLst, i)
		segs = append(segs, s)
	}
	// one segment that is not transforms would be read as a plain statement
	if _, ok := gfSynModsLike(pLst[0]); !ok && len(segs) == 1 {
		segs[0] = "| " + segs[0]
	}

	oneLine := strings.Join(segs, " | ")
	if pIndent*4+len(oneLine) <= 110 || len(segs) == 1 {
		p.line(pIndent, oneLine)
		return
	}
	for i, s := range segs {
		if i == 0 {
			p.line(pIndent, s+" |")
		} else if i < len(segs)-1 {
			p.line(pIndent+1, s+" |")
		} else {
			p.line(pIndent+1, s)
		}
	}
}

// one element (or transforms + the atoms after them) as a segment string. returns the next index.
func (p *gfSynPrinter) segElem(pLst []interface{}, pI int) (string, int) {

	e := pLst[pI]
	lst, isLst := e.([]interface{})
	if !isLst {
		return p.tok(e), pI + 1
	}

	// transforms, followed by the atoms (geometry)
	if mods, ok := gfSynModsLike(lst); ok {
		toks := p.kvToks(mods)
		i := pI + 1
		for ; i < len(pLst); i++ {
			if _, isL := pLst[i].([]interface{}); isL {
				break
			}
			toks = append(toks, p.tok(pLst[i]))
		}
		return strings.Join(toks, " "), i
	}

	if name, val, ok := p.setterShape(lst); ok {
		return name + " = " + val, pI + 1
	}

	// multiple words -> a flat segment (read back as a list), otherwise one group/call token
	if len(lst) >= 2 {
		if s, isStr := lst[0].(string); isStr && gfSynSafeWord(s) && !gfSynIsKeyword(s) {
			toks := []string{}
			for _, v := range lst {
				toks = append(toks, p.tok(v))
			}
			return strings.Join(toks, " "), pI + 1
		}
	}
	return p.tok(lst), pI + 1
}

// ["set", name, value] that can be written as `name = value`
func (p *gfSynPrinter) setterShape(pLst []interface{}) (string, string, bool) {
	if len(pLst) != 3 || pLst[0] != "set" {
		return "", "", false
	}
	name, ok := pLst[1].(string)
	if !ok || !gfSynSafeWord(name) || gfSynIsKeyword(name) || gfSynInfixOps[name] || strings.Contains(name, "=") {
		return "", "", false
	}
	return name, p.tok(pLst[2]), true
}

// keyword blocks and setters. returns false if the list is not of those shapes.
func (p *gfSynPrinter) stmtSpecial(pLst []interface{}, pIndent int) bool {

	if len(pLst) < 3 {
		return false
	}
	head, ok := pLst[0].(string)
	if !ok {
		return false
	}

	switch head {

	case "*":
		scope, ok := pLst[2].([]interface{})
		if len(pLst) != 3 || !ok {
			return false
		}
		p.line(pIndent, "repeat "+p.tok(pLst[1])+":")
		p.scope(scope, pIndent+1)
		return true

	case "rule":
		name, ok := pLst[1].(string)
		if !ok || !gfSynSafeWord(name) {
			return false
		}
		header := "rule " + name
		var body []interface{}
		if len(pLst) == 3 {
			body, ok = pLst[2].([]interface{})
			if !ok {
				return false
			}
		} else if len(pLst) == 4 {
			mods, ok1 := pLst[2].([]interface{})
			body, ok = pLst[3].([]interface{})
			if !ok1 || !ok || len(mods) == 0 {
				return false
			}
			for _, m := range mods {
				pair, ok := m.([]interface{})
				if !ok || len(pair) != 2 {
					return false
				}
				k, ok := pair[0].(string)
				if !ok || !gfSynIdentRe.MatchString(k) || k == "max" {
					return false
				}
				if k == "iters_max" {
					k = "max"
				}
				header += " " + k + "=" + p.tok(pair[1])
			}
		} else {
			return false
		}
		p.line(pIndent, header+":")
		p.stmts(body, pIndent+1)
		return true

	case "if":
		body, ok := pLst[2].([]interface{})
		if len(pLst) != 3 || !ok {
			return false
		}
		cond := p.tok(pLst[1])
		if c, ok := pLst[1].([]interface{}); ok && len(c) == 3 {
			if op, ok := c[0].(string); ok && gfSynInfixOps[op] {
				cond = p.tok(c[1]) + " " + op + " " + p.tok(c[2])
			}
		}
		p.line(pIndent, "if "+cond+":")
		p.stmts(body, pIndent+1)
		return true

	case "shader":
		name, ok := pLst[1].(string)
		if !ok || !gfSynSafeWord(name) {
			return false
		}
		uniformsCount := 0
		for _, part := range pLst[2:] {
			if gfSynIsUniforms(part) {
				uniformsCount++
			}
		}
		if uniformsCount > 1 {
			return false
		}
		p.line(pIndent, "shader "+name+":")
		for _, part := range pLst[2:] {
			if gfSynIsUniforms(part) {
				for _, u := range part.([]interface{})[1].([]interface{}) {
					toks := []string{"uniform"}
					for _, v := range u.([]interface{}) {
						toks = append(toks, p.tok(v))
					}
					p.line(pIndent+1, strings.Join(toks, " "))
				}
				continue
			}
			partLst, ok := part.([]interface{})
			if !ok {
				p.fail("shader part is not a list")
				return true
			}
			toks := []string{}
			for _, v := range partLst {
				if s, ok := v.(string); ok && strings.Contains(s, "\n") {
					toks = append(toks, gfSynTripleQuote(s, pIndent+1))
				} else {
					toks = append(toks, p.tok(v))
				}
			}
			p.line(pIndent+1, strings.Join(toks, " "))
		}
		return true

	case "set":
		name, val, ok := p.setterShape(pLst)
		if !ok {
			return false
		}
		p.line(pIndent, name+" = "+val)
		return true
	}
	return false
}

func gfSynIsUniforms(pPart interface{}) bool {
	l, ok := pPart.([]interface{})
	if !ok || len(l) != 2 || l[0] != "uniforms" {
		return false
	}
	us, ok := l[1].([]interface{})
	if !ok {
		return false
	}
	for _, u := range us {
		if _, ok := u.([]interface{}); !ok {
			return false
		}
	}
	return true
}

// elements of a `repeat` scope (flattened transforms and geometry)
func (p *gfSynPrinter) scope(pLst []interface{}, pIndent int) {

	for i := 0; i < len(pLst); {
		if lst, ok := pLst[i].([]interface{}); ok {
			if _, isMods := gfSynModsLike(lst); !isMods {
				if p.stmtSpecial(lst, pIndent) {
					i++
					continue
				}
				// a nested statement that starts with a list
				if len(lst) > 0 {
					if _, ok := lst[0].([]interface{}); ok {
						p.chain(lst, pIndent)
						i++
						continue
					}
				}
			}
		}
		var s string
		s, i = p.segElem(pLst, i)
		p.line(pIndent, s)
	}
}

//-------------------------------------------------
// TOKENS

// list of [key, value, key, value...] with identifier keys
func gfSynModsLike(pE interface{}) ([]interface{}, bool) {
	lst, ok := pE.([]interface{})
	if !ok || len(lst) < 2 || len(lst)%2 != 0 {
		return nil, false
	}
	for i := 0; i < len(lst); i += 2 {
		k, ok := lst[i].(string)
		if !ok || !gfSynTransformKeys[k] {
			return nil, false
		}
	}
	return lst, true
}

func (p *gfSynPrinter) kvToks(pMods []interface{}) []string {
	out := []string{}
	for i := 0; i < len(pMods); i += 2 {
		out = append(out, pMods[i].(string)+"="+p.tok(pMods[i+1]))
	}
	return out
}

func gfSynIsKeyword(pStr string) bool {
	switch pStr {
	case "repeat", "rule", "if", "shader", "uniform":
		return true
	}
	return false
}

// a string that can be written as a bare word
func gfSynSafeWord(pStr string) bool {
	if pStr == "" || strings.ContainsAny(pStr, " \t\r\n,(){}\"") || strings.Contains(pStr, "//") ||
		strings.HasSuffix(pStr, ":") || pStr == "true" || pStr == "false" {
		return false
	}
	if _, err := strconv.ParseFloat(pStr, 64); err == nil {
		return false
	}
	if strings.Contains(pStr, "=") {
		return gfSynInfixOps[pStr]
	}
	return true
}

func (p *gfSynPrinter) tok(pV interface{}) string {

	switch v := pV.(type) {
	case string:
		if gfSynSafeWord(v) && !gfSynIsKeyword(v) {
			return v
		}
		return strconv.Quote(v)

	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)

	case bool:
		return strconv.FormatBool(v)

	case []interface{}:

		// nested statement that starts with a list: {a | b}
		if len(v) > 0 {
			if _, ok := v[0].([]interface{}); ok {
				segs := []string{}
				for i := 0; i < len(v); {
					var s string
					s, i = p.segElem(v, i)
					segs = append(segs, s)
				}
				return "{" + strings.Join(segs, " | ") + "}"
			}
		}

		// infix: (a op b)
		if len(v) == 3 {
			if op, ok := v[0].(string); ok && gfSynInfixOps[op] {
				return "(" + p.tok(v[1]) + " " + op + " " + p.tok(v[2]) + ")"
			}
		}

		args := v
		prefix := ""
		if len(v) >= 1 {
			if name, ok := v[0].(string); ok && gfSynCallRe.MatchString(name) && name != "true" && name != "false" {
				prefix = name
				args = v[1:]
			}
		}
		toks := []string{}
		for _, a := range args {
			toks = append(toks, p.tok(a))
		}
		// "a op b" shaped args would be read as an infix expression - quote the middle one
		if len(args) == 3 {
			if op, ok := args[1].(string); ok && gfSynInfixOps[op] {
				toks[1] = strconv.Quote(op)
			}
		}
		return prefix + "(" + strings.Join(toks, " ") + ")"

	default:
		p.fail("unsupported value %v (%T)", pV, pV)
		return ""
	}
}

// """ block, content indented one level deeper than the line. falls back to a quoted string
// when the content would not survive the dedent.
func gfSynTripleQuote(pStr string, pIndent int) string {

	lines := strings.Split(pStr, "\n")
	ok := len(lines) >= 2 && !strings.Contains(pStr, `"""`) && !strings.ContainsAny(pStr, "\t\r") &&
		strings.TrimSpace(lines[0]) != "" && strings.TrimSpace(lines[len(lines)-1]) != ""
	if ok {
		minInd := -1
		for _, l := range lines {
			if strings.TrimSpace(l) == "" {
				continue
			}
			ind := len(l) - len(strings.TrimLeft(l, " "))
			if minInd < 0 || ind < minInd {
				minInd = ind
			}
		}
		ok = minInd == 0
	}
	if !ok {
		return strconv.Quote(pStr)
	}
	ind := strings.Repeat("    ", pIndent)
	var sb strings.Builder
	sb.WriteString("\"\"\"\n")
	for _, l := range lines {
		if strings.TrimSpace(l) == "" {
			sb.WriteString(l + "\n")
		} else {
			sb.WriteString(ind + "    " + l + "\n")
		}
	}
	sb.WriteString(ind + `"""`)
	return sb.String()
}
