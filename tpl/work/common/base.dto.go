package common

type BaseRes struct {
	ApiCode    int    `swaggertype:"integer" validate:"required"`
	ApiMessage string `swaggertype:"string" validate:"required"`
}
