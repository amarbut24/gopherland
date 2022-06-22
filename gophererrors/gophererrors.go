package gophererrors

import (
	"fmt"

	msgraph_errors "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

func HandleODataErr(err error, context string) error {
	oderr := err.(*msgraph_errors.ODataError).GetError()
	c := *oderr.GetCode()
	m := *oderr.GetMessage()
	return fmt.Errorf("%v\nCode=%v\nmessage=%v", context, c, m)
}

func GetODataDetails(err error) (string, string) {
	oderr := err.(*msgraph_errors.ODataError).GetError()
	c := *oderr.GetCode()
	m := *oderr.GetMessage()
	return c, m
}
