package nlm

import (
	"embed"
	"fmt"
	"io"
	"math/rand"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed articles/*
var articleFiles embed.FS

type Article struct {
	Title      string `yaml:"title"`
	Subtitle   string `yaml:"subtitle"`
	Paragraphs []string
}

func GetArticle() *Article {
	ret := &Article{
		Title:    "Lorem Ipsum",
		Subtitle: "Sometimes",
		Paragraphs: []string{
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec aliquam eros at venenatis vulputate. Vestibulum ullamcorper est et euismod tristique. Morbi sit amet facilisis sem, eu dignissim massa. Quisque gravida, dolor id maximus pharetra, tortor massa mattis arcu, sit amet porttitor sem libero ut est. Nunc eleifend ex id est finibus, eu blandit erat faucibus. Donec dui enim, vulputate eu lorem at, dictum gravida est. Vivamus ac ligula magna.",
			"Nullam feugiat justo ac bibendum elementum. Vestibulum volutpat urna vel felis venenatis ultricies. Nulla mauris ante, volutpat ac eros sit amet, volutpat imperdiet ex. Integer finibus tellus ut metus luctus, pellentesque interdum quam blandit. Aenean luctus ipsum at quam sagittis euismod. Pellentesque vitae felis eget lorem vestibulum maximus at non arcu. Ut ultricies mi a pharetra laoreet. Curabitur venenatis mollis massa, vel tincidunt dui. Etiam vel dolor volutpat, commodo leo a, cursus velit. Quisque efficitur diam id sem lacinia, nec bibendum purus rutrum. Vestibulum mattis sed lacus ut fermentum. In id interdum felis. Morbi quis ligula neque. Suspendisse sed fermentum justo, quis dapibus ligula. Vestibulum sed metus sit amet tortor suscipit interdum ac at eros. Etiam quis nibh erat.",
			"Proin dignissim accumsan ipsum ut rutrum. Vestibulum sit amet consectetur dui. Nunc sagittis rutrum arcu, quis vehicula libero volutpat non. Fusce consectetur ligula lacus, a viverra dolor tincidunt ut. Curabitur eget tincidunt dolor. Nullam porta odio in ex ultricies aliquet. Vivamus nec vehicula tortor. Quisque convallis justo vel dui vestibulum euismod. Maecenas maximus mattis quam vel luctus. Proin gravida velit vel nisi tincidunt viverra. Vivamus vitae odio eget enim tincidunt euismod.",
			"Quisque lorem sapien, luctus non aliquet et, dictum non purus. Praesent ante felis, tempor a risus et, malesuada tristique ante. Suspendisse non urna viverra, consequat eros in, rhoncus risus. Pellentesque quis diam ut sapien molestie molestie. Aliquam pretium nibh nec rhoncus auctor. Morbi hendrerit dolor vitae dui bibendum laoreet. Etiam eu tempus orci. Proin aliquet eget ipsum a pretium. Sed scelerisque diam bibendum dolor hendrerit, consectetur dignissim nibh mollis. Aenean faucibus, urna vel fermentum scelerisque, nisl mi congue velit, in auctor tortor dolor eget ligula. Donec finibus quis orci a tempor. Sed eleifend laoreet nibh, quis pellentesque elit. Praesent fringilla vestibulum consectetur. Suspendisse convallis ultrices sem, sit amet tincidunt leo.",
			"Etiam porta ut ex vel euismod. In hac habitasse platea dictumst. Donec mauris ligula, varius nec felis in, dapibus accumsan diam. Aliquam sed lorem blandit, vestibulum massa consectetur, gravida lacus. Vestibulum lectus augue, consectetur vel hendrerit et, rhoncus quis quam. Integer eu rutrum massa, in tempus dui. Aenean convallis ex vitae velit molestie, ut accumsan felis dictum. Sed rutrum quis ex in ornare. Morbi dictum mi vitae nulla dictum, id scelerisque mauris semper. Praesent maximus metus tortor, nec ultricies felis mattis id. In purus justo, euismod eu lorem nec, efficitur vestibulum libero. Donec ultricies eleifend lorem sed volutpat. Suspendisse tempor cursus eros, non egestas orci fermentum at."},
	}

	if rand.Float64() < 0.5 {
		ret.Subtitle = ""
	}

	rand.Shuffle(len(ret.Paragraphs), func(i, j int) { ret.Paragraphs[i], ret.Paragraphs[j] = ret.Paragraphs[j], ret.Paragraphs[i] })

	return ret
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
