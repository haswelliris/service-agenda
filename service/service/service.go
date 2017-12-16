package service

import (
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"

	"github.com/haswelliris/service-agenda/service/model"
)

// NewServer returns a negroni server that has already initialized
// routes
func NewServer() *negroni.Negroni {

	model.init()

	formatter := render.New(render.Options{IndentJSON: true})

	n := negroni.Classic()
	mx := mux.NewRouter()

	initRoutes(mx, formatter)

	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
	// # User
	// ### user register
	mx.HandleFunc("/v1/user/register", userRegisterHandler(formatter)).Methods("POST")
	// ### user login
	mx.HandleFunc("user/login", userLoginHandler(formatter)).Methods("GET")
	// ### user logout
	mx.HandleFunc("/v1/user/logout", userLogoutHandler(formatter)).Methods("POST")

	// # Users
	// ### List all Users
	mx.HandleFunc("/v1/users", listAllUsersHandler(formatter)).Methods("GET")

	// # Meeting
	// ### Create a new meeting
	mx.HandleFunc("/v1/meeting", createNewMeetingHandler(formatter)).Methods("POST")

	// ### List all meetings
	mx.HandleFunc("/v1/meetings", listAllMeetingsHandler(formatter)).Methods("GET")

}
