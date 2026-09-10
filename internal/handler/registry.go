package handler

import healthhandler "cleaning/internal/handler/health"

type Register struct {
	Health *healthhandler.Health
}

func NewRegistry() *Register {
	return &Register{
		Health: healthhandler.NewHandler(),
	}
}
