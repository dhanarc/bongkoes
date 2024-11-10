package deployment

import (
	"bufio"
	"fmt"
	"github.com/samber/lo"
	"os"
	"strings"
)

func readStdIn(prompt string) (*string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return nil, err
	}

	// Trim the newline character from the input
	input = strings.TrimSpace(input)
	return lo.ToPtr(input), nil
}
