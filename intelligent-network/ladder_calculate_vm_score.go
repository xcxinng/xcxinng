package intelligentnetwork

// 定义计算得分的函数，按照实际情况调节三者比重
func calculateMachineScore(cpuUsage, memoryUsage, systemLoad float64) float64 {
	cpuWeight := 0.3
	memoryWeight := 0.3
	loadWeight := 0.4

	// 非线性模型：随着系统负载增加，得分加速下降
	normalizedLoadScore := 1 / (1 + (systemLoad * systemLoad))

	cpuScore := 1 - cpuUsage
	memoryScore := 1 - memoryUsage

	loadScore := normalizedLoadScore

	totalScore := cpuScore*cpuWeight + memoryScore*memoryWeight + loadScore*loadWeight
	return totalScore
}
