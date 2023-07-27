package processor

import (
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

type TextProcessor struct {
	text string
}

func NewTextProcessor(inTxt string) *TextProcessor {
	return &TextProcessor{inTxt}
}

func (processor *TextProcessor) SanitizeText() *TextProcessor {
	var sb strings.Builder
	cr := regexp.MustCompile(`\n+`)
	tab := regexp.MustCompile(`\t+`)

	processor.text = cr.ReplaceAllString(processor.text, "\n")
	processor.text = tab.ReplaceAllString(processor.text, "\t")
	txtBlocks := strings.Split(processor.text, "\n")

	for _, blk := range txtBlocks {

		blk = strings.TrimSpace(blk)
		if len(blk) > 0 {
			sb.WriteString(blk + "\n")
		}
	}
	processor.text = sb.String()
	return processor
}
func (processor *TextProcessor) SanitizeHtml() *TextProcessor {

	p := bluemonday.UGCPolicy()

	processor.text = p.Sanitize(processor.text)

	return processor
}

func (processor *TextProcessor) String() string {

	return processor.text
}
