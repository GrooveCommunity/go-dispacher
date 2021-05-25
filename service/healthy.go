package service

type Healthy struct {
	Status string
}

func ValidateHealthy() Healthy {
	return Healthy{Status: "Success!"}
}
