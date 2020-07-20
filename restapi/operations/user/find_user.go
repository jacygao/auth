// Code generated by go-swagger; DO NOT EDIT.

package user

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// FindUserHandlerFunc turns a function with the right signature into a find user handler
type FindUserHandlerFunc func(FindUserParams) middleware.Responder

// Handle executing the request and returning a response
func (fn FindUserHandlerFunc) Handle(params FindUserParams) middleware.Responder {
	return fn(params)
}

// FindUserHandler interface for that can handle valid find user params
type FindUserHandler interface {
	Handle(FindUserParams) middleware.Responder
}

// NewFindUser creates a new http.Handler for the find user operation
func NewFindUser(ctx *middleware.Context, handler FindUserHandler) *FindUser {
	return &FindUser{Context: ctx, Handler: handler}
}

/*FindUser swagger:route GET /user/find/{id} user findUser

Find user

Find a user by User ID

*/
type FindUser struct {
	Context *middleware.Context
	Handler FindUserHandler
}

func (o *FindUser) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewFindUserParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}