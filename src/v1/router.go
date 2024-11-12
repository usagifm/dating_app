package v1

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"github.com/usagifm/dating-app/src/middleware/auth"
	"github.com/usagifm/dating-app/src/v1/handler"

	_ "github.com/usagifm/dating-app/swagger/v1"
	// _ "github.com/usagifm/dating-app/swagger/v1/swagger.json" // Import your generated Swagger docs
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func post(r chi.Router, svc, path string, handlerFunc http.HandlerFunc) {
	r.With(otelhttp.NewMiddleware(fmt.Sprintf("%s %s%s", "POST", svc, path))).
		Post(path, handlerFunc)
}

func get(r chi.Router, svc, path string, handlerFunc http.HandlerFunc) {
	r.With(otelhttp.NewMiddleware(fmt.Sprintf("%s %s%s", "GET", svc, path))).
		Get(path, handlerFunc)
}

func Router(r *chi.Mux, deps *Dependency) {
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Mount("/swagger", httpSwagger.WrapHandler)

	r.Route("/api/v1/auth", func(authRoutes chi.Router) {
		authenticateRoute := authRoutes.With(auth.New(auth.DefaultAuthenticatedRequestValidator))
		authPath := "auth"
		// auth
		post(authRoutes, authPath, "/signup", handler.SignUp(&deps.Handlers.DatingApp))
		post(authRoutes, authPath, "/login", handler.Login((&deps.Handlers.DatingApp)))
		get(authenticateRoute, authPath, "/profile", handler.GetProfile((&deps.Handlers.DatingApp)))
	})

	r.Route("/api/v1/dating", func(datingRoutes chi.Router) {
		authenticateRoute := datingRoutes.With(auth.New(auth.DefaultAuthenticatedRequestValidator))
		datingPath := "dating"
		// dating

		get(authenticateRoute, datingPath, "/preference", handler.GetUserPreference((&deps.Handlers.DatingApp)))
		post(authenticateRoute, datingPath, "/preference", handler.UpdateUserPreference(&deps.Handlers.DatingApp))
		get(authenticateRoute, datingPath, "/", handler.GetProfilesByPreference((&deps.Handlers.DatingApp)))
		post(authenticateRoute, datingPath, "/swipe", handler.Swipe(&deps.Handlers.DatingApp))
		get(authenticateRoute, datingPath, "/matches", handler.GetUserMatches((&deps.Handlers.DatingApp)))

		get(authenticateRoute, datingPath, "/package", handler.GetPackages((&deps.Handlers.DatingApp)))
		post(authenticateRoute, datingPath, "/package/buy", handler.BuyPackage(&deps.Handlers.DatingApp))

	})
}
