package models

type Network struct {
	Name string //The name of the network
}

func (network *Network) GetName() *string {
	if network == nil {
		return nil
	}
	return &network.Name
}
