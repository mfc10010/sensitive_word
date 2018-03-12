package sensitive

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	sensitive_key    = "isEnd"
	sensitive_flag_t = "1"
	sensitive_flag_f = "0"
)

type SensitiveWords struct {
	sensitiveWordMap map[string]interface{}
	key_word_set  map[string]bool
}


//var sensitiveWordMap map[string]interface{}

func  (g *SensitiveWords)  readLine(filePth string, key_word_set map[string]bool) error {
	f, err := os.Open(filePth)
	if err != nil {
		return err
	}
	defer f.Close()
	bfRd := bufio.NewReader(f)
	for {
		line_, err := bfRd.ReadBytes('\n')
		if err != nil { //遇到任何错误立即返回，并忽略 EOF 错误信息
			if err == io.EOF {
				return nil
			}
			return err
		}
		line := string(line_)
		hookfn(key_word_set, line) //放在错误处理前面，即使发生错误，也会处理已经读取到的数据。
	}
	return nil
}

func hookfn(key_word_set map[string]bool, line string) {
	line = strings.Replace(line, "\n", "", -1)
	line = strings.Replace(line, "\r", "", -1)
	key_word_set[line] = true
}

//"./config/sensitive.txt"
func (g *SensitiveWords) InitkeyWord(path string) {
	g.key_word_set = make(map[string]bool, 0)
	g.readLine(path, g.key_word_set)
	g.sensitiveWordMap = make(map[string]interface{}, len(g.key_word_set))
	g.loadTree(g.key_word_set)

}

func  (g *SensitiveWords)isChinese(str string) bool {
	var hzRegexp = regexp.MustCompile("^[\u4e00-\u9fa5]$")
	return hzRegexp.MatchString(str)
}


func (g *SensitiveWords) loadTree(key_word_set map[string]bool) {
	var nowMap map[string]interface{}
	var newWorMap map[string]interface{}

	for k, _ := range key_word_set {
		nowMap = g.sensitiveWordMap

		for i, j := range k {
			s := string(j)
			if g.isChinese(s) {
				c, ok := nowMap[s]
				if ok {
					nowMap = c.(map[string]interface{})
				} else {
					newWorMap = make(map[string]interface{})
					newWorMap[sensitive_key] = sensitive_flag_f
					nowMap[s] = newWorMap
					nowMap = newWorMap
				}
				if i == len(k)-3 {
					newWorMap[sensitive_key] = sensitive_flag_t
				}

			} else {
				c, ok := nowMap[s]
				if ok {
					nowMap = c.(map[string]interface{})
				} else {
					newWorMap = make(map[string]interface{})
					newWorMap[sensitive_key] = sensitive_flag_f
					nowMap[s] = newWorMap
					nowMap = newWorMap
				}
				if i == len(k)-1 {
					newWorMap[sensitive_key] = sensitive_flag_t
				}
			}

		}
	}
}

//过滤敏感词
func  (g *SensitiveWords) GetSensitiveWord(txt string) string {
	ru, _ := utf8.DecodeRuneInString("*")
	rs := []rune(txt)
	n := len(rs)
	for j := 0; j < n; j++ {
		l := g.checkSensitiveWord(txt, j)
		if l > 0 {
			for i := 0; i < l; i++ {
				rs[j+i] = ru
			}
		}
	}
	return string(rs)
}

//检查是否存在敏感词
func  (g *SensitiveWords)CheckExistSensitive(txt string) bool {
	rs := []rune(txt)
	n := len(rs)
	for j := 0; j < n; j++ {
		l := g.checkSensitiveWord(txt, j)
		if l > 0 {
			return true
		}
	}
	return false
}

func  (g *SensitiveWords)checkSensitiveWord(txt string, beginIndex int) int {
	flag := false
	matchFlag := 0

	var nowMap map[string]interface{}

	nowMap = g.sensitiveWordMap

	rs := []rune(txt)
	n := len(rs)
	for i := beginIndex; i < n; i++ {
		word := rs[i]
		str_word := string(word)

		s, ok := nowMap[str_word]
		if ok {
			nowMap = s.(map[string]interface{})
			matchFlag++
			if c, ok := nowMap[sensitive_key]; ok && c == sensitive_flag_t {
				flag = true
				break
			}
		} else {
			break
		}

	}
	if matchFlag < 2 && !flag {
		matchFlag = 0
	}

	if flag {
		return matchFlag
	}
	return 0
}
