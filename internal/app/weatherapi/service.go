package weatherapi



type Service struct {
	client Client
}

func NewService(client Client) Service {
	return Service{
		client: client,
	}
}