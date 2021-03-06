package processors

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	gostatic "github.com/passingcuriosity/gostatic/lib"
)

type TagsProcessor struct {
}

func NewTagsProcessor() *TagsProcessor {
	return &TagsProcessor{}
}

func (p *TagsProcessor) Process(page *gostatic.Page, args []string) error {
	return ProcessTags(page, args)
}

func (p *TagsProcessor) Description() string {
	return "generate taxonomy pages for tags mentioned in page header " +
		"(arguments: [field] template)"
}

func (p *TagsProcessor) Mode() int {
	return gostatic.Pre
}

func SplitTags(str string) ([]string, error) {
	tags := []string{}
	for _, t := range strings.Split(str, ",") {
		tags = append(tags, strings.TrimSpace(t))
	}
	return tags, nil
}

func ProcessTags(page *gostatic.Page, args []string) error {
	var fieldName string
	var pathPattern string

	if len(args) == 1 {
		fieldName = "Tags"
		pathPattern = args[0]
	} else if len(args) == 2 {
		fieldName = gostatic.Capitalize(args[0])
		pathPattern = args[1]
	} else {
		return errors.New("'tags' rule needs one or arguments")
	}

	debug("Looking for tags in field %s\n", fieldName)

	var tags []string = nil
	if fieldName == "Tags" {
		tags = page.Tags
	} else {
		t, exists := page.Other[fieldName]
		if exists {
			ts, err := SplitTags(t)
			if err == nil {
				tags = ts
			}
		} else {
			debug("Looking didn't find tags in field %s\n", fieldName)
			debug("%s", page.Other)
		}
	}

	if tags == nil {
		return nil
	}

	site := page.Site

	for _, tag := range tags {
		tagpath := strings.Replace(pathPattern, "*", tag, 1)

		debug("Found page for '%s' at '%s'\n", tag, tagpath)

		tagpage := site.Pages.BySource(tagpath)
		if tagpage == nil {
			pattern, rules := site.Rules.MatchedRules(tagpath)
			if rules == nil {
				return fmt.Errorf("Tag path '%s' does not match any rule", tagpath)
			}
			if len(rules) > 1 {
				return fmt.Errorf("Tags are not supported with multiple rules. Tag in question: '%s'", tagpath)
			}

			tagpage := &gostatic.Page{
				PageHeader: gostatic.PageHeader{Title: tag},
				Site:       site,
				Pattern:    pattern,
				Deps:       append(make(gostatic.PageSlice, 0), page),
				Rule:       rules[0],
				Source:     tagpath,
				Path:       tagpath,
				// tags are never new, because they only depend on pages and
				// have not a bit of original content
				ModTime: time.Unix(0, 0),
			}
			tagpage.SetWasRead(true)
			page.Site.Pages = append(page.Site.Pages, tagpage)
			tagpage.Peek()
			debug("Added '%s' to '%s'\n", page.Source, tagpage.Source)
		} else {
			tagpage.Deps = append(tagpage.Deps, page)
			debug("Added '%s' to '%s'\n", page.Source, tagpage.Source)
		}
	}

	return nil
}

func debug(format string, args ...interface{}) {
	fmt.Printf("tags: "+format, args...)
	os.Stdout.Sync()
}
