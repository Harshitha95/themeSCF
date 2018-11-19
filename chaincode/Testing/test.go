package test
import (
	"fmt"
	"math/rand"
	"github.com/satori/go.uuid"

)

func main() {
	fmt.Println("Hello, playground")
	for i := 0; i < 10; i++ {
	genUUIDv4()
	}
}
func genUUIDv4() {
    id := uuid.NewV4()
    fmt.Printf("github.com/satori/go.uuid:   %s\n", id)
}