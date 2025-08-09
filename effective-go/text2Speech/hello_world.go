package main

import (
	"os/exec"
)

func main() {
	// 场景1：测试清晰的短句
	sentence1 := "To be or not to be, that is the question."
	exec.Command("espeak", "-ven+f3", "-s130", sentence1).Run()

	// 场景2：测试长句和连读
	sentence2 := "The quick brown fox jumps over the lazy dog."
	exec.Command("espeak", "-ven+m2", "-s150", "-k8", sentence2).Run()

	// 场景3：测试情感语调调整（故意放慢）
	sentence3 := "I have a dream that one day this nation will rise up..."
	exec.Command("espeak", "-ven-us", "-s100", "-a120", sentence3).Run()
}
