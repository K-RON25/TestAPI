package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

// Options for the CLI.
type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"8888"`
}

// GreetingInput represents the greeting operation request.
type GreetingInput struct {
	Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
}

// GreetingOutput represents the greeting operation response.
type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

// ResetPasswordInput represents the reset password operation request.
type ResetPasswordInput struct {
	Password string `path:"password" minLength:"12" example:"@12jk.,['l]" doc:"Password to reset"`
}

// ResetPasswordOutput represents the greeting operation response.
type ResetPasswordOutput struct {
	Body struct {
		Message string `json:"message" example:"This is a password reset" doc:"Pasword reset"`
	}
}

func main() {
	// Create a CLI app which takes a port option.
	cli := huma.NewCLI(func(hooks huma.Hooks, options *Options) {
		// Create a new router & API
		router := chi.NewMux()
		api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))

		// Register GET /greeting/{name}
		huma.Register(api, huma.Operation{
			OperationID: "get-greeting",
			Summary:     "Get a greeting",
			Method:      http.MethodGet,
			Path:        "/greeting/{name}",
		}, func(ctx context.Context, input *GreetingInput) (*GreetingOutput, error) {
			resp := &GreetingOutput{}
			resp.Body.Message = fmt.Sprintf("Hello, %s!", input.Name)
			return resp, nil
		})

		// Register GET /greeting/{name}
		huma.Register(api, huma.Operation{
			OperationID: "reset-password",
			Summary:     "Reset a password",
			Method:      http.MethodPost,
			Path:        "/reset-password/{password}",
		}, func(ctx context.Context, input *ResetPasswordInput) (*ResetPasswordOutput, error) {
			resp := &ResetPasswordOutput{}
			resp.Body.Message = fmt.Sprintf("Congrats, you reset a password %s", input.Password)
			return resp, nil
		})

		// Tell the CLI how to start your router.
		hooks.OnStart(func() {
			http.ListenAndServe(fmt.Sprintf(":%d", options.Port), router)
		})
	})

	// Run the CLI. When passed no commands, it starts the server.
	cli.Run()
}
