package instruction_cycle

func Noop(Tiempo int) int {
	return Tiempo
}

func GOTO(pcInstr int, valor int) int {
	return pcInstr + valor
}
