package camino

import (
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	idPrefix    = "#"
	classPrefix = "."
)

type Matcher interface {
	Match(node *html.Node) bool
}

func Match(doc *html.Node, delegate Matcher) *html.Node {
	matches := AllMatches(doc, delegate, 1)
	if len(matches) > 0 {
		return matches[0]
	}
	return nil
}

func AllMatches(doc *html.Node, delegate Matcher, limit int) []*html.Node {
	matches := make([]*html.Node, 0)

	var f func(*html.Node)
	f = func(n *html.Node) {

		if delegate.Match(n) {
			matches = append(matches, n)
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
			if len(matches) == limit {
				return
			}
		}
	}

	f(doc)

	return matches
}

type elementAtom struct {
	atom atom.Atom
}

func AtomMatcher(a atom.Atom) Matcher {
	return &elementAtom{
		atom: a,
	}
}

func (ea *elementAtom) Match(node *html.Node) bool {
	if node.DataAtom == ea.atom {
		return true
	}
	return false
}

type atomId struct {
	atom atom.Atom
	id   string
}

func AtomIdMatcher(a atom.Atom, id string) Matcher {
	return &atomId{
		atom: a,
		id:   id,
	}
}

func (eti *atomId) Match(node *html.Node) bool {
	if node.DataAtom != eti.atom ||
		(eti.id != "" && len(node.Attr) == 0) {
		return false
	}

	for _, attr := range node.Attr {
		if attr.Key == "id" {
			return attr.Val == eti.id
		}
	}

	return false
}

type atomClass struct {
	atom   atom.Atom
	class  string
	equals bool
}

func AtomClassMatcher(a atom.Atom, class string, equals bool) Matcher {
	return &atomClass{
		atom:   a,
		class:  class,
		equals: equals,
	}
}

func (etc *atomClass) Match(node *html.Node) bool {
	if node.DataAtom != etc.atom ||
		(etc.class != "" && len(node.Attr) == 0) {
		return false
	}

	if etc.class == "" {
		return true
	}

	for _, attr := range node.Attr {
		if attr.Key == "class" {
			if etc.equals {
				return attr.Val == etc.class
			}

			return strings.Contains(attr.Val, etc.class)
		}
	}

	return false
}

type selector struct {
	atom  atom.Atom
	id    string
	class string
}

func SelectorMatcher(query string) Matcher {
	tagName, id, class := "", "", ""

	if strings.Contains(query, idPrefix) {
		rest := ""
		tagName, rest, _ = strings.Cut(query, idPrefix)
		if strings.Contains(rest, classPrefix) {
			id, class, _ = strings.Cut(rest, classPrefix)
		} else {
			id = rest
		}
	} else if strings.Contains(query, classPrefix) {
		tagName, class, _ = strings.Cut(query, classPrefix)
	} else {
		tagName = query
	}

	a := atom.Lookup([]byte(tagName))

	return &selector{
		atom:  a,
		id:    id,
		class: class,
	}
}

func (s *selector) Match(node *html.Node) bool {

	if (s.atom != 0 && node.DataAtom != s.atom) ||
		(s.id != "" && len(node.Attr) == 0) ||
		(s.class != "" && len(node.Attr) == 0) {
		return false
	}

	if s.id != "" {
		for _, attr := range node.Attr {
			if attr.Key == "id" && attr.Val == s.id {
				return true
			}
		}
		return false
	}

	if s.class != "" {
		for _, attr := range node.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, s.class) {
				return true
			}
		}
		return false
	}

	return true
}
