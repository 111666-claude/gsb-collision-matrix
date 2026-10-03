// Command collision-matrix 跑碰撞矩阵样例。
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"example.com/collision"
)

// Run 执行一次命令行调用，返回退出码。
func Run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("collision-matrix", flag.ContinueOnError)
	flags.SetOutput(stderr)
	sample := flags.String("sample", "oneway", "oneway / trigger / dup / order / scan")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	switch *sample {
	case "oneway":
		matrix := collision.New()
		matrix.SetLayer("player", 2)
		matrix.SetLayer("mob", 0)
		matrix.AddBody("p", "player", false)
		matrix.AddBody("m", "mob", false)
		fmt.Fprintf(stdout, "hit=%v\n", matrix.ShouldHit("p", "m"))
	case "trigger":
		matrix := collision.New()
		matrix.SetLayer("player", 2)
		matrix.SetLayer("mob", 1)
		matrix.AddBody("p", "player", false)
		matrix.AddBody("m", "mob", true)
		fmt.Fprintf(stdout, "pairs=%d\n", len(matrix.Pairs()))
	case "dup":
		matrix := collision.New()
		matrix.SetLayer("player", 2)
		matrix.SetLayer("mob", 1)
		matrix.AddBody("p", "player", false)
		matrix.AddBody("m", "mob", false)
		fmt.Fprintf(stdout, "pairs=%d\n", len(matrix.Pairs()))
	case "order":
		matrix := collision.New()
		matrix.SetLayer("player", 3)
		matrix.SetLayer("mob", 3)
		matrix.AddBody("p", "player", false)
		matrix.AddBody("m2", "mob", false)
		matrix.AddBody("m1", "mob", false)
		pairs := matrix.Pairs()
		first := "-"
		if len(pairs) > 0 {
			first = pairs[0].A + "," + pairs[0].B
		}
		fmt.Fprintf(stdout, "pairs=%d first=%s\n", len(pairs), first)
	case "scan":
		matrix := collision.New()
		for index := 0; index < 3000; index++ {
			matrix.SetLayer(fmt.Sprintf("l-%d", index), 1)
		}
		matrix.AddBody("a", "l-0", false)
		matrix.AddBody("b", "l-1", false)
		for index := 0; index < 3000; index++ {
			matrix.ShouldHit("a", "b")
		}
		fmt.Fprintf(stdout, "scanned=%d\n", matrix.Scanned())
	default:
		fmt.Fprintln(stderr, "需要 --sample oneway|trigger|dup|order|scan")
		return 2
	}
	return 0
}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}
