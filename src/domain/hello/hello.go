package hello

type HelloService struct{}

func NewHelloService() *HelloService {
	return &HelloService{}
}

func (s *HelloService) Hello() string {
	return "Hello World"
}
