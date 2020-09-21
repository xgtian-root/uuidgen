# UUIDGen

全局唯一ID生成器(雪花算法)


## Example

```

package main

import (
	"fmt"

	"git.youxuetong.com/Micro/uuidgen"
)

func main() {
    // 如果你是在k8s内部使用，请使用
    // workID, err := uuidgen.GetK8sWorkID()
    
    // 如果是其他应用场景，请传入不大于65535的数字，确保不同实例的workID不一样
    sf, err := uuidgen.New(1)
    
    if err != nil {
        panic(err)
    }


    // Generate a snowflake ID.
    uuid := sf.Generate()

    // Print
    fmt.Println(uuid)
}

```