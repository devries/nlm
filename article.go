package nlm

import (
	"embed"
	"fmt"
	"io"
	"math/rand"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed articles/*.txt
var articleFiles embed.FS

type Article struct {
	Title      string `yaml:"title"`
	Subtitle   string `yaml:"subtitle"`
	Paragraphs []string
}

func LoadArticle(r io.Reader) (*Article, error) {
	var buf strings.Builder

	_, err := io.Copy(&buf, r)
	if err != nil {
		return nil, err
	}

	var article Article
	sections := strings.Split(buf.String(), "---\n")

	if len(sections) != 3 {
		return nil, fmt.Errorf("article format is incorrect. Article needs to start with three dashes and have three dashes separating the preamble from the article")
	}

	err = yaml.Unmarshal([]byte(sections[1]), &article)
	if err != nil {
		return nil, fmt.Errorf("unable to parse preamble: %s", err)
	}

	article.Paragraphs = strings.Split(sections[2], "\n\n")

	return &article, nil
}

func GetArticleList() ([]string, error) {
	dirEntries, err := articleFiles.ReadDir("articles")
	if err != nil {
		return nil, fmt.Errorf("unable to load articles: %s", err)
	}

	filenames := make([]string, len(dirEntries))
	for i, v := range dirEntries {
		filenames[i] = v.Name()
	}

	return filenames, nil
}

func GetNamedArticle(name string) (*Article, error) {
	r, err := articleFiles.Open(fmt.Sprintf("articles/%s", name))
	if err != nil {
		return nil, fmt.Errorf("error opening article: %s", err)
	}
	defer r.Close()

	article, err := LoadArticle(r)
	if err != nil {
		return nil, fmt.Errorf("error reading article: %s", err)
	}

	return article, nil
}

type ArticleBuilder struct {
	TitleMarkov     MarkovSource `json:"title_markov"`
	ContentMarkov   MarkovSource `json:"content_markov"`
	TitleStarters   []string     `json:"title_starters"`
	ContentStarters []string     `json:"content_starters"`
}

// NewArticleBuilder creates a Markov article builder with a chain
// size of titlesize for titles and size for content.
func NewArticleBuilder(titlesize, size int) (*ArticleBuilder, error) {
	articles, err := GetArticleList()
	if err != nil {
		return nil, fmt.Errorf("unable to get article list: %s", err)
	}

	titlemb := NewMarkovBuilder()
	contentmb := NewMarkovBuilder()
	titleStarters := []string{}
	contentStarters := []string{}

	for _, articleName := range articles {
		article, err := GetNamedArticle(articleName)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve article %s: %s", articleName, err)
		}

		titleStarters = append(titleStarters, getPrefix(article.Title, titlesize))
		titlemb.AddText(article.Title, titlesize, EndOfDocument)

		for i, para := range article.Paragraphs {
			var endRune rune
			if i >= len(article.Paragraphs)-1 {
				endRune = EndOfDocument
			} else {
				endRune = EndOfParagraph
			}

			contentStarters = append(contentStarters, getPrefix(para, size))
			contentmb.AddText(para, size, endRune)
		}
	}

	ret := ArticleBuilder{titlemb.ConvertToSource(), contentmb.ConvertToSource(), titleStarters, contentStarters}
	return &ret, nil
}

func getPrefix(s string, size int) string {
	r := []rune(s)
	return string(r[:size])
}

func (ab *ArticleBuilder) GenerateArticle(titleSize, articleSize int) *Article {
	titleStart := ab.TitleStarters[rand.Int()%len(ab.TitleStarters)]
	title := ab.TitleMarkov.GenerateText(titleStart, titleSize)[0]

	contentStart := ab.ContentStarters[rand.Int()%len(ab.ContentStarters)]
	content := ab.ContentMarkov.GenerateText(contentStart, articleSize)

	ret := Article{title, "", content}

	return &ret
}
