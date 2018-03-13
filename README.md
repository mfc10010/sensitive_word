# sensitive_word

sensitive_word is a tool for filter word.you can use it to replace sensitive word to "*".
## installtion

### use to get
```
go get -u github.com/mfc10010/sensitive_word
```
## usage

```go
package main

import (
	"github.com/sensitive_word"
	"fmt"
)

func main()  {
	s:=new(sensitive.SensitiveWords)
	s.InitkeyWord("./sensitive.txt")
	fmt.Println(s.GetSensitiveWord("办证ss11翻八九民communistgggtttthe Communist Party"))
}
```

## print result 
```
**ss11翻************gggttt*******************	
```
