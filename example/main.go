package main

import (
	"github.com/sensitive_word/key"
	"fmt"
)

func main()  {
	s:=new(sensitive.SensitiveWords)
	s.InitkeyWord("./sensitive.txt")
	fmt.Println(s.GetSensitiveWord("办证ss11翻八九民communistgggtttthe Communist Party"))
}