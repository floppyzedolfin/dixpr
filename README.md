## Requirements
- You need to manually download the dictionary you want from [the wikimedia site](https://dumps.wikimedia.org/backup-index.html).
  - I used frwiktionary and then "All pages, current version only"


## Usage
Run `> go run . ${wikimedia}.tar.bz2`

This loads the file, and suggests a prompt for, first, a list of words to match, and, then, a pronunciation to match.
Pronunciations use the IPA. 
