package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type DeviceModel struct {
	Units int `json:"units" validate:"required,lte=52,gte=1"`
}

func Validate() {
	v := validator.New()

	dm := DeviceModel{
		Units: -1,
	}
	err := v.Struct(&dm)
	if err != nil {
		panic(err)
	}
	fmt.Println("OK")
}
