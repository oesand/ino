package openapi

import (
	"fmt"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"testing"

	"github.com/oesand/ino"
)

type TestUser struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
}

func TestGenerateSchema(t *testing.T) {
	routes := ino.Routes(
		Tag("Users"),
		ino.Get("/users/{id}", ino.ParamHandler1(
			ino.PathParam[int64]("id"),
			func(id int64, w http.ResponseWriter) {

			},
		)),
		ino.Post("/users", ino.ParamHandler1(ino.JsonParam[TestUser](), func(body *TestUser, w http.ResponseWriter) {

		})),
		ino.Post("/upload", ino.ParamHandler1(ino.FileParam("file"), func(file *multipart.FileHeader, w http.ResponseWriter) {

		})),
	)

	mux := ino.New(routes...)

	mux.Middleware(func(writer http.ResponseWriter, request *http.Request, next http.Handler) {
		log.Printf("[%s]%s", request.Method, request.RequestURI)

		next.ServeHTTP(writer, request)
	})

	schema, err := GenerateSchema(mux, nil)
	if err != nil {
		t.Fatal(err)
	}

	err = Register(mux, schema, RegisterOptions{})
	if err != nil {
		panic(err)
	}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	fmt.Println("Listening on", listener.Addr())

	err = http.Serve(listener, mux)
	if err != nil {
		panic(err)
	}
}
