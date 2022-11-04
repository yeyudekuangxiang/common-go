package initialize

func Close() {
	closeQueueProducer()
	closeLogger()
}
