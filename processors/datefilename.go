package processors

import (
	"regexp"
	"time"
	"path/filepath"

	gostatic "github.com/passingcuriosity/gostatic/lib"
)

type DatefilenameProcessor struct {
}

func NewDatefilenameProcessor() *DatefilenameProcessor {
	return &DatefilenameProcessor{}
}

func (p *DatefilenameProcessor) Process(page *gostatic.Page, args []string) error {
	return ProcessDatefilename(page, args)
}

func (p *DatefilenameProcessor) Description() string {
	return "process filename 2014-05-06-name.md to path /name.html and set page.Date to 2014-05-06"
}

func (p *DatefilenameProcessor) Mode() int {
	return gostatic.Pre
}

func ProcessDatefilename(page *gostatic.Page, args []string) error {
	name := page.Name()
	dir := filepath.Dir(page.Path)

	validName := regexp.MustCompile(`(\d{4}-\d{2}-\d{2})-(.*)`)
	if validName.MatchString(name) {
		fnamecomponents := validName.FindStringSubmatch(name)
		t, err := time.Parse("2006-01-02", fnamecomponents[1])
		if err == nil {
			page.Date = t
			page.Path = dir + "/" + fnamecomponents[2]
		}
	}
	return nil
}
