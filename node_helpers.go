package camino

import (
	"strings"

	"golang.org/x/net/html"
)

func TextContent(node *html.Node) string {
	sb := &strings.Builder{}

	var f func(*html.Node)
	f = func(n *html.Node) {

		if n.Type == html.TextNode {
			sb.WriteString(n.Data)
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(node)

	return sb.String()
}

func GetAttribute(node *html.Node, attribute string) string {
	if node == nil {
		return ""
	}
	for _, attr := range node.Attr {
		if attr.Key == attribute {
			return attr.Val
		}
	}
	return ""
}
