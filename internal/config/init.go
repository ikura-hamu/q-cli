package config

type Init interface {
	GetForce() (bool, error)
}
