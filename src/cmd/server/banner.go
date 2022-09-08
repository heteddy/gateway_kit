// @Author : detaohe
// @File   : banner
// @Description:
// @Date   : 2022/9/7 20:30

package main

import (
	"fmt"
	"github.com/lukesampson/figlet/figletlib"
	"os"
	"path/filepath"
)

func printBanner() {
	cwd, _ := os.Getwd()
	fontsDir := filepath.Join(cwd, "fonts")
	f, err := figletlib.GetFontByName(fontsDir, "standard")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not find that font!")
		return
	}
	figletlib.FPrintMsg(os.Stdout, "gateway-kit", f, 100, f.Settings(), "center")
}
func printBanner2() {
	fmt.Println("                             _                                 _    _ _   ")
	fmt.Println("                  __ _  __ _| |_ _____      ____ _ _   _      | | _(_) |_ ")
	fmt.Println("                 / _` |/ _` | __/ _ \\ \\ /\\ / / _` | | | |_____| |/ / | __|")
	fmt.Println("                | (_| | (_| | ||  __/\\ V  V / (_| | |_| |_____|   <| | |_ ")
	fmt.Println("                 \\__, |\\__,_|\\__\\___| \\_/\\_/ \\__,_|\\__, |     |_|\\_\\_|\\__|")
	fmt.Println("                 |___/                             |___/                  ")
}
