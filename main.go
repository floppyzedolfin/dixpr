package main

import (
	"fmt"
	"os"
	"slices"

	"golang.org/x/sync/errgroup"

	"github.com/floppyzedolfin/dixpr/wikidata"
)

func main() {
	eg := errgroup.Group{}

	wordChan := make(chan *wikidata.Word)
	eg.Go(func() error {
		return wikidata.Load(os.Args[1], wordChan)
	})

	eg.Go(func() error {
		forbiddenPronunciations := []rune{'r', 'ʁ', 'ɹ', 'e', 'ɛ'}
		for word := range wordChan {
			sp := []rune(word.Spelling)
			ipa := []rune(word.IPA)
			if sp[len(sp)-1] == 'r' && !slices.Contains(forbiddenPronunciations, ipa[len(ipa)-1]) {
				fmt.Printf("%s -- /%s/\n", word.Spelling, word.IPA)
			}
		}
		return nil
	})

	err := eg.Wait()
	if err != nil {
		panic(err)
	}
}
