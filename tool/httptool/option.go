package httptool

type Option func(client *HttpClient)

func WithLogger(logger Logger) Option {
	return func(client *HttpClient) {
		client.logger = logger
	}
}
