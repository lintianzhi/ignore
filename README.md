
## 符合.gitignore规则的包
[http://git-scm.com/docs/gitignore](http://git-scm.com/docs/gitignore)

规则1：

- 含有`/` : 
    - `r`  -> `\A/r(/**)?`
- 不含`/` : 
    - `r`  -> `**/r(/**)?`
    - `r/` -> `**/r/**`

规则2：

- `*`  -> `[^/]*`
- `**` -> `.*`

EXAMPLE:

```go
package main

import (
	"fmt"
	"github.com/lintianzhi/ignore"
)

func main() {
	gitIgn, _ := ignore.NewGitIgn(".gitignore")
	gitIgn.Start(".")

	ignored := gitIgn.IgnoreList()
	for _, v := range ignored {
		fmt.Println(v)
	}
}
```
