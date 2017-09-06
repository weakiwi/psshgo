package main
import "fmt"
func main() {
    ch := make(chan int, 3)
    ch <- 1
    ch <- 2
    ch <- 3

	for v := range ch {
	    fmt.Println(v)
	    if len(ch) <= 0 { // 如果现有数据量为0，跳出循环
	        break
	    }
	}
}
