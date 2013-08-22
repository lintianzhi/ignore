package ignore

import (
	"bufio"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

//含有/ : r  -> \A/r(/**)?
//不含/ : r  -> **/r(/**)?
//		  r/ -> **/r/**

//        *  -> [^/]*
//        ** -> .*

type Ignore struct {
	Ignore   []*regexp.Regexp
	Excluded []*regexp.Regexp
}

func NewIgnore() *Ignore {
	return &Ignore{make([]*regexp.Regexp, 0, 10), make([]*regexp.Regexp, 0, 10)}
}

func (ignore *Ignore) ParseLine(line string) {
	line = strings.TrimSpace(line)
	if len(line) == 0 || line[0] == '#' {
		return
	}
	if line[0] != '!' {
		if line[0] == '\\' {
			line = line[1:]
		}
		ignore.ParseIgnore(line)
	} else {
		ignore.ParseExcluded(line[1:])
	}
}

func (ignore *Ignore) Parse(line string) *regexp.Regexp {
	if strings.Index(line, "/") >= 0 {
		//	如果含有`/`则是当前目录开始的绝对路径
		line = "\\A/" + line

		if line[len(line)-1] != '/' {
			line = line + "(/**)?"
		} else {
			line = line + "**"
		}
	} else {
		line = "**/" + line + "(/**)?"
	}

	line = replace(line)

	return regexp.MustCompile(line)
}

var (
	STAR_TMP = "//.//"
	TWO_STAR = STAR_TMP + STAR_TMP
)

func replace(line string) string {
	line = path.Clean(line)
	line = strings.Replace(line, "*", STAR_TMP, -1)
	r := strings.NewReplacer(TWO_STAR, ".*", STAR_TMP, "[^/]*")
	return r.Replace(line)
}

func (ignore *Ignore) ParseIgnore(line string) {
	ignore.Ignore = append(ignore.Ignore, ignore.Parse(line))
}

func (ignore *Ignore) ParseExcluded(line string) {
	ignore.Excluded = append(ignore.Excluded, ignore.Parse(line))
}

func (ignore *Ignore) MatchIgnore(path string) bool {
	for _, v := range ignore.Ignore {
		if v.MatchString(path) {
			return true
		}
	}
	return false
}

func (ignore *Ignore) MatchExcluded(path string) bool {
	for _, v := range ignore.Excluded {
		if v.MatchString(path) {
			return true
		}
	}
	return false
}

type GitIgn struct {
	Ign      *Ignore
	IgnFiles []string
}

func NewGitIgn(path string) (*GitIgn, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	gitIgn := &GitIgn{NewIgnore(), make([]string, 0, 0)}

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		gitIgn.Ign.ParseLine(line)
	}

	return gitIgn, nil
}

func (gitIgn *GitIgn) TestIgnore(path string) bool {
	if len(path) == 0 {
		return false
	}
	if path[0] != '/' {
		path = "/" + path
	}

	return gitIgn.Ign.MatchIgnore(path) && !gitIgn.Ign.MatchExcluded(path)

}

func (gitIgn *GitIgn) Start(path string) {
	gitIgn.IgnFiles = make([]string, 0, 0)

	walkAction := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if gitIgn.TestIgnore(path) {
			gitIgn.IgnFiles = append(gitIgn.IgnFiles, path)
		}

		return nil
	}

	filepath.Walk(path, walkAction)
}

func (gitIgn *GitIgn) IgnoreList() []string {
	return gitIgn.IgnFiles
}

