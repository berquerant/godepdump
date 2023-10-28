package write

import (
	"encoding/json"
	"fmt"
	"os"
)

func Format(format string, v ...any) {
	fmt.Fprintf(os.Stdout, format, v...)
}

func JSON(v any) {
	b, _ := json.Marshal(v)
	fmt.Fprintln(os.Stdout, string(b))
}
