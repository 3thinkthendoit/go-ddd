package main

import (
	"github.com/go-spring/spring-core/gs"
	_ "think.com/go-ddd/infrastructure/component/server"
)

func init() {
	gs.Setenv("GS_SPRING_CONFIG_LOCATIONS", "infrastructure/common/config/")
}

func main() {
	gs.Run()
}
