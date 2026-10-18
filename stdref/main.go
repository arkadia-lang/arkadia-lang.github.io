// Command stdref writes the standard library reference of the website: one
// page per module of lib/, from what `ark doc` prints for it, and an index
// of the modules. The pages are HTML fragments wrapped in the site template
// by pandoc, like the other pages.
//
//	go run ./stdref -ark ../arkadia/bin/ark -lib ../arkadia/lib -out _site/std
package main

import (
	"bytes"
	"flag"
	"fmt"
	"html"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func main() {
	ark := flag.String("ark", "bin/ark", "the ark command")
	lib := flag.String("lib", "lib", "the standard library")
	out := flag.String("out", "_site/std", "the output directory")
	template := flag.String("template", "template.html", "the pandoc template")
	flag.Parse()
	log.SetFlags(0)
	log.SetPrefix("stdref: ")

	names, err := moduleNames(*lib)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(*out, 0o755); err != nil {
		log.Fatal(err)
	}
	var mods []module
	for _, name := range names {
		text, err := exec.Command(*ark, "doc", name).Output()
		if err != nil {
			log.Fatalf("ark doc %s: %v", name, err)
		}
		mods = append(mods, parse(name, string(text)))
	}
	page := func(file, title string, body []byte) {
		if err := render(*template, filepath.Join(*out, file), title, body); err != nil {
			log.Fatalf("%s: %v", file, err)
		}
	}
	page("index.html", "標準ライブラリ", indexPage(mods))
	for _, m := range mods {
		page(m.name+".html", m.name, modulePage(m, mods))
	}
}

// moduleNames lists the modules of the library: Prelude, then Std.* in
// alphabetical order.
func moduleNames(lib string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(lib, "Std", "*.ark"))
	if err != nil {
		return nil, err
	}
	var names []string
	for _, f := range files {
		names = append(names, "Std."+strings.TrimSuffix(filepath.Base(f), ".ark"))
	}
	sort.Strings(names)
	return append([]string{"Prelude"}, names...), nil
}

// render wraps an HTML fragment in the site template with pandoc.
func render(template, file, title string, body []byte) error {
	frag, err := os.CreateTemp("", "stdref-*.html")
	if err != nil {
		return err
	}
	defer os.Remove(frag.Name())
	if _, err := frag.Write(body); err != nil {
		return err
	}
	if err := frag.Close(); err != nil {
		return err
	}
	cmd := exec.Command("pandoc", "-f", "markdown", "-t", "html5", "-s",
		"--template", template, "-B", frag.Name(),
		"-M", "pagetitle="+title, "-V", "root=../", "-V", "nav-std=true",
		"-o", file)
	cmd.Stdin = strings.NewReader("")
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ------------------------------------------------------------ parsing

// The output of `ark doc`: the module line and its documentation indented
// by 4, then the declarations at column 0, each followed by its
// documentation indented by 4 and, for classes and effects, by its members
// at column 2 with their documentation indented by 6.

type module struct {
	name  string
	doc   []string
	decls []decl
}

type decl struct {
	sig     string
	doc     []string
	members []decl
}

func parse(name, text string) module {
	m := module{name: name}
	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
	var cur *[]string // the documentation lines being collected
	cur = &m.doc
	for _, line := range lines[1:] {
		switch indent := len(line) - len(strings.TrimLeft(line, " ")); {
		case strings.TrimSpace(line) == "":
			*cur = append(*cur, "")
		case indent == 0:
			m.decls = append(m.decls, decl{sig: line})
			cur = &m.decls[len(m.decls)-1].doc
		case indent == 2 && len(m.decls) > 0:
			d := &m.decls[len(m.decls)-1]
			d.members = append(d.members, decl{sig: line[2:]})
			cur = &d.members[len(d.members)-1].doc
		default:
			*cur = append(*cur, line)
		}
	}
	return m
}

// ------------------------------------------------------------ rendering

var (
	keywordRe = regexp.MustCompile(`^(type|class|effect) (\S+)`)
	codeRe    = regexp.MustCompile("`([^`]+)`")
)

// declName is the name a declaration defines, used as its anchor.
func declName(sig string) string {
	if m := keywordRe.FindStringSubmatch(sig); m != nil {
		return m[2]
	}
	name, _, _ := strings.Cut(sig, " : ")
	return name
}

// summary is the first sentence of the module documentation.
func summary(m module) string {
	var words []string
	for _, l := range m.doc {
		if strings.TrimSpace(l) == "" {
			break
		}
		words = append(words, strings.TrimSpace(l))
	}
	s := strings.Join(words, " ")
	if i := strings.Index(s, ". "); i >= 0 {
		s = s[:i+1]
	}
	return s
}

func inline(s string) string {
	return codeRe.ReplaceAllString(html.EscapeString(s), "<code>$1</code>")
}

// docHTML renders documentation: paragraphs of text, and code where the
// lines are indented further than the text.
func docHTML(b *bytes.Buffer, lines []string) {
	base := -1
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			n := len(l) - len(strings.TrimLeft(l, " "))
			if base < 0 || n < base {
				base = n
			}
		}
	}
	isCode := func(l string) bool { return len(l) > base && l[base] == ' ' }
	for i := 0; i < len(lines); {
		switch l := lines[i]; {
		case strings.TrimSpace(l) == "":
			i++
		case isCode(l):
			// code runs to the next text line, without trailing blank lines
			j := i
			for j < len(lines) && (strings.TrimSpace(lines[j]) == "" || isCode(lines[j])) {
				j++
			}
			for strings.TrimSpace(lines[j-1]) == "" {
				j--
			}
			var code []string
			for _, c := range lines[i:j] {
				if len(c) > base+2 {
					code = append(code, c[base+2:])
				} else {
					code = append(code, "")
				}
			}
			fmt.Fprintf(b, "<pre><code>%s</code></pre>\n", html.EscapeString(strings.Join(code, "\n")))
			i = j
		default:
			var text []string
			for i < len(lines) && strings.TrimSpace(lines[i]) != "" && !isCode(lines[i]) {
				text = append(text, strings.TrimSpace(lines[i]))
				i++
			}
			fmt.Fprintf(b, "<p>%s</p>\n", inline(strings.Join(text, " ")))
		}
	}
}

func declHTML(b *bytes.Buffer, d decl, anchor string) {
	idAttr := ""
	if anchor != "" {
		idAttr = fmt.Sprintf(` id="%s"`, html.EscapeString(anchor))
	}
	fmt.Fprintf(b, "<div class=\"decl\"%s>\n<pre class=\"sig\"><code>%s</code></pre>\n", idAttr, html.EscapeString(d.sig))
	if len(d.doc) > 0 {
		b.WriteString("<div class=\"decl-doc\">\n")
		docHTML(b, d.doc)
		b.WriteString("</div>\n")
	}
	if len(d.members) > 0 {
		b.WriteString("<div class=\"members\">\n")
		for _, mem := range d.members {
			declHTML(b, mem, "")
		}
		b.WriteString("</div>\n")
	}
	b.WriteString("</div>\n")
}

func sidebar(b *bytes.Buffer, current string, mods []module) {
	b.WriteString("<nav class=\"sidebar\" aria-label=\"モジュール\">\n<ul>\n")
	for _, m := range mods {
		cur := ""
		if m.name == current {
			cur = ` aria-current="page"`
		}
		fmt.Fprintf(b, "<li><a href=\"%s.html\"%s>%s</a></li>\n", m.name, cur, m.name)
	}
	b.WriteString("</ul>\n</nav>\n")
}

func indexPage(mods []module) []byte {
	var b bytes.Buffer
	b.WriteString("<div class=\"with-sidebar\">\n")
	sidebar(&b, "", mods)
	b.WriteString("<article class=\"doc\">\n<h1>標準ライブラリ</h1>\n")
	b.WriteString("<p><code>Prelude</code> はすべてのモジュールに暗黙にインポートされます。<code>Std.*</code> のモジュールは <code>import</code> して使います。" +
		"同じ内容はコマンドラインの <code>ark doc &lt;Module&gt;</code> でも表示できます。</p>\n")
	b.WriteString("<dl class=\"modules\">\n")
	for _, m := range mods {
		fmt.Fprintf(&b, "<dt><a href=\"%s.html\">%s</a></dt>\n<dd>%s</dd>\n", m.name, m.name, inline(summary(m)))
	}
	b.WriteString("</dl>\n</article>\n</div>\n")
	return b.Bytes()
}

func modulePage(m module, mods []module) []byte {
	var b bytes.Buffer
	b.WriteString("<div class=\"with-sidebar\">\n")
	sidebar(&b, m.name, mods)
	fmt.Fprintf(&b, "<article class=\"doc\">\n<h1>%s</h1>\n", m.name)
	docHTML(&b, m.doc)
	if len(m.decls) > 0 {
		b.WriteString("<h2>一覧</h2>\n<ul class=\"decl-index\">\n")
		for _, d := range m.decls {
			n := declName(d.sig)
			fmt.Fprintf(&b, "<li><a href=\"#%s\"><code>%s</code></a></li>\n", html.EscapeString(n), html.EscapeString(n))
		}
		b.WriteString("</ul>\n<h2>定義</h2>\n")
	}
	for _, d := range m.decls {
		declHTML(&b, d, declName(d.sig))
	}
	b.WriteString("</article>\n</div>\n")
	return b.Bytes()
}
