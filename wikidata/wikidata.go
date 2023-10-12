package wikidata

import (
	"compress/bzip2"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func Load(path string, wordChan chan<- *Word) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("unable to open file: %w", err)
	}
	defer f.Close()

	decoder := xml.NewDecoder(bzip2.NewReader(f))

	for {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("unable to decode token: %w", err)
		}

		switch se := token.(type) {
		case xml.StartElement:
			// if it's a page:
			if se.Name.Local == "page" {
				var page Page
				err := decoder.DecodeElement(&page, &se)
				if err != nil {
					return fmt.Errorf("unable to decode element: %w", err)
				}

				pronunciations := getPronunciations(page.Revisions[0].Text)
				for _, p := range pronunciations {
					wordChan <- &Word{Spelling: page.Title, IPA: p}
				}
			}
		}
	}

	close(wordChan)
	return nil
}

func getPronunciations(contents string) []string {
	// in each section, look for the pronunciations.
	var pronunciationBlock = regexp.MustCompile("''' {{pron\\|([^\\|]+?)\\|fr}}")

	frenchPronunciations := pronunciationBlock.FindAllString(contents, -1)
	pronunciations := make([]string, 0)
	mapPronunciations := make(map[string]bool)
	for _, fPron := range frenchPronunciations {
		ipa := strings.TrimPrefix(fPron, "''' {{pron|")
		ipa = strings.TrimSuffix(ipa, "|fr}}")

		ipa = strings.TrimSpace(ipa)
		if mapPronunciations[ipa] {
			continue
		}

		mapPronunciations[ipa] = true
		pronunciations = append(pronunciations, ipa)
	}
	return pronunciations
}

type Word struct {
	Spelling string
	IPA      string
}

type Revision struct {
	XMLName xml.Name `xml:"revision"`
	Text    string   `xml:"text"`
}

type Page struct {
	XMLName   xml.Name   `xml:"page"`
	Title     string     `xml:"title"`
	Revisions []Revision `xml:"revision"`
}

type MediaWiki struct {
	XMLName xml.Name `xml:"mediawiki"`
	Pages   []Page   `xml:"page"`
}
