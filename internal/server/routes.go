package server

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/markbates/goth/gothic"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		t, _ := template.New("foo").Parse(indexTemplate)
		t.Execute(writer, nil)
	})
	r.Get("/health", s.healthHandler)
	r.Get("/auth/{provider}/callback", s.getAuthCallBack)
	r.Get("/logout/{provider}", s.logOut)
	r.Get("/auth", gothic.BeginAuthHandler)

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}

func (s *Server) getAuthCallBack(res http.ResponseWriter, req *http.Request) {
	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		fmt.Println("error ", err)
		http.Error(res, "Authentication failed", http.StatusInternalServerError)
		return
	}

	// Fetch the user's role from the database based on their email
	role, err := s.getUserRole(user.Email)
	if err != nil {
		fmt.Println("error fetching user role: ", err)
		http.Error(res, "Failed to fetch user role", http.StatusInternalServerError)
		return
	}

	// Add the user's role to the context or session (optional)
	ctx := context.WithValue(req.Context(), "userRole", role)
	req = req.WithContext(ctx)

	// Check if the user has access to the page based on their role
	if role == "admin" {
		// Admins can access any page
		t, _ := template.New("foo").Parse(adminTemplate)
		t.Execute(res, user)
	} else if role == "user" {
		// Regular users can access a different page
		t, _ := template.New("foo").Parse(userTemplate)
		t.Execute(res, user)
	} else {
		// Handle unauthorized access for other roles
		http.Error(res, "Forbidden", http.StatusForbidden)
	}
}

var adminTemplate = `
<p><a href="/logout/{{.Provider}}">logout</a></p>
<h1>Admin Page</h1>
<p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
<p>ExpiresAt: {{.ExpiresAt}}</p>
<p>RefreshToken: {{.RefreshToken}}</p>
`

func (s *Server) getUserRole(email string) (string, error) {
	// Use the database service to fetch the user role
	role, err := s.db.GetUserRole(email)
	if err != nil {
		return "", fmt.Errorf("failed to get user role: %w", err)
	}
	return role, nil
}

func (s *Server) logOut(w http.ResponseWriter, r *http.Request) {
	gothic.Logout(w, r)
	r.Header.Set("Location", "/")
}

// func (s *Server) provider(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("000000000000------")

// 	if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
// 		t, _ := template.New("foo").Parse(userTemplate)
// 		t.Execute(w, gothUser)
// 	} else {
// 		fmt.Println("xxxxx------", err, gothUser)

// 		gothic.BeginAuthHandler(w, r)
// 	}
// }

var userTemplate = `
<p><a href="/logout/{{.Provider}}">logout</a></p>
<p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
<p>ExpiresAt: {{.ExpiresAt}}</p>
<p>RefreshToken: {{.RefreshToken}}</p>
`

var indexTemplate = `
<p><a href="/auth?provider=twitter">Log in with Twitter</a></p>
<p><a href="/auth?provider=facebook">Log in with Facebook</a></p>
<p><a href="/auth?provider=google">Log in with Google</a></p>
`
