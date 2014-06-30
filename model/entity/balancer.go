package entity

type Balancer struct {
	Id   string
	Args map[string]interface{}
}

func NewBalancer(id string, data map[string]interface{}) (Balancer, error) {
	return Balancer{
		Id:   id,
		Args: data,
	}, nil
}
