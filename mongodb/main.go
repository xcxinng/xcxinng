package main

// docker 构建单实例环境即可

func main() {
	// 插入模拟数据
	// RunInsert()
	// RunInsertDeviceConfig()

	// 测试动态连表pipeline查询
	// RunDynamicJoin()

	// watch
	RunSyncMongoToEs()
}
