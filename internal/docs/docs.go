package docs

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/elliotchance/orderedmap/v2"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/anchor"
)

const (
	DocsBaseDir = "./docs"
)

type metadata struct {
	title   string
	summary string
	tags    []string
}

type file struct {
	metadata metadata
	data     bytes.Buffer
}

type docs struct {
	*orderedmap.OrderedMap[string, file]
}

var (
	onceDocs sync.Once

	docsFs *docs
)

func GetDocs() *docs {

	onceDocs.Do(func() { // <-- atomic, does not allow repeating
		var err error
		docsFs, err = loadDocs()
		if err != nil {
			panic(err)
		}

		log.Println("Docs loaded")
	})

	return docsFs
}

func loadDocs() (*docs, error) {
	docsFS := orderedmap.NewOrderedMap[string, file]()

	markdown := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithBlockParsers(),
			parser.WithInlineParsers(),
			parser.WithParagraphTransformers(),
			parser.WithASTTransformers(),
			parser.WithAutoHeadingID(), // Required by anchor
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
			extension.Typographer,
			meta.Meta,
			highlighting.NewHighlighting(
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
			&anchor.Extender{
				Texter: anchor.Text("🔗"),
			},
		),
	)

	files, err := os.ReadDir(DocsBaseDir)
	if err != nil {
		return nil, err
	}

	for _, path := range files {
		// Read Markdown and convert it to HTML
		filePath := DocsBaseDir + "/" + path.Name()
		source, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		var buf bytes.Buffer

		context := parser.NewContext()
		err = markdown.Convert(source, &buf, parser.WithContext(context))
		if err != nil {
			panic(err)
		}
		metaData := meta.Get(context)

		// Write generated HTML into virtual filesystem
		fsFilepath := strings.TrimSuffix(path.Name(), ".md")
		docsFS.Set(fsFilepath, file{
			metadata: metadata{
				title:   metaData["Title"].(string),
				summary: metaData["Summary"].(string),
				tags:    ConvertInterfaceSliceToStringSlice(metaData["Tags"].([]interface{})),
			},
			data: buf,
		})

		log.Println("docs: loaded file", fsFilepath, "from", filePath)
	}

	return &docs{docsFS}, nil
}

func ConvertInterfaceSliceToStringSlice(slice []interface{}) []string {
	var strSlice []string
	for _, v := range slice {
		str, ok := v.(string)
		if !ok {
			panic(fmt.Errorf("element is not a string: %v", v))
		}
		strSlice = append(strSlice, str)
	}
	return strSlice
}
