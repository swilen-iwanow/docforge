package reactor

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func TestMe(t *testing.T) {
	hrefAttrMatchRegex := regexp.MustCompile(`href=["\']?([^"\'>]+)["\']?`)
	srcAttrMatchRegex := regexp.MustCompile(`src=["\']?([^"\'>]+)["\']?`)
	p := parser.NewParser(parser.WithBlockParsers(parser.DefaultBlockParsers()...),
		parser.WithInlineParsers(parser.DefaultInlineParsers()...),
		parser.WithParagraphTransformers(parser.DefaultParagraphTransformers()...),
	)
	source := []byte("This is `inline element`\n\n*empahsized text*\n## Heading 2\n[GitHub](\"https://github.com\") ![ImgTitle](\"https://somehwere.org/someurl\")\n ## Another heading 2\n <p><a href=\"https://github.com\">alabala</a> <img src=\"../images/logo.png\"></p>")

	reader := text.NewReader(source)
	doc := p.Parse(reader)
	modifiedDoc := make([]byte, len(source))
	if err := ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if node.Kind() == ast.KindLink {
				n := node.(*ast.Link)
				fmt.Printf("Link [%s](%s)\n", n.Text(source), n.Destination)
				n.Destination = []byte("changed")
				fmt.Printf("-->Link [%s](%s)\n", n.Text(source), n.Destination)
			}
			if node.Kind() == ast.KindImage {
				n := node.(*ast.Image)
				fmt.Printf("Image ![](%s)\n", n.Destination)
			}
			if node.Kind() == ast.KindRawHTML {
				n := node.(*ast.RawHTML)
				l := n.Segments.Len()
				for i := 0; i < l; i++ {
					segment := n.Segments.At(i)
					segmentStr := segment.Value(source)
					match := hrefAttrMatchRegex.Find([]byte(segmentStr))
					if len(match) > 0 {
						url := strings.Split(string(match), "=")[1]
						fmt.Printf("HREF attribute value: %v\n", url)
						continue
					}
					match = srcAttrMatchRegex.Find([]byte(segmentStr))
					if len(match) > 0 {
						url := strings.Split(string(match), "=")[1]
						fmt.Printf("SRC attribute value: %v\n", url)
						continue
					}
				}
				// return ast.WalkSkipChildren, nil
			}
			if node.Type() == ast.TypeBlock {
				lines := node.Lines()
				l := lines.Len()
				for i := 0; i < l; i++ {
					segment := lines.At(i)
					modifiedDoc = append(modifiedDoc, segment.Value(source)...)
					//fmt.Printf("%s\n", segment.Value(source))
				}
			}
		}
		//fmt.Print(string(node.Text(source)))
		return ast.WalkContinue, nil
	}); err != nil {
		t.Errorf("%v", err)
	}

	fmt.Printf("%v\n", string(modifiedDoc))

	// fmt.Printf("%v\n", string(source))
	// if err := ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
	// 	if !entering {
	// 		return ast.WalkContinue, nil
	// 	}
	// 	if node.Type() == ast.TypeBlock {
	// 		lines := node.Lines()
	// 		l := lines.Len()
	// 		for i := 0; i < l; i++ {
	// 			segment := lines.At(i)
	// 			fmt.Printf("%s\n", segment.Value(source))
	// 		}
	// 	}
	// 	return ast.WalkContinue, nil
	// }); err != nil {
	// 	t.Errorf("%v", err)
	// }
	// fmt.Println()
}
