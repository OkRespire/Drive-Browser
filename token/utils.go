package token

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

func openURL(url string) {
	var err error

	fmt.Println("Opening in your browser...")
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		fmt.Println("Please open the following URL manually:", url)
	}
	if err != nil {
		log.Fatalf("Unable to open browser: %v", err)
	}

}
