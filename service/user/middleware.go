package user

import (
	"context"
	"log"
	"net/http"

	ory "github.com/ory/client-go"
)

// save the cookies for any upstream calls to the Ory apis
func withCookies(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, "req.cookies", v)
}

func getCookies(ctx context.Context) string {
	return ctx.Value("req.cookies").(string)
}

// save the session to display it on the dashboard
func withSession(ctx context.Context, v *ory.Session) context.Context {
	return context.WithValue(ctx, "req.session", v)
}

func getSession(ctx context.Context) *ory.Session {
	return ctx.Value("req.session").(*ory.Session)
}

func (s *Service) ValidateSession(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("handling middleware request\n")

		// set the cookies on the ory client
		var cookies string

		// this example passes all request.Cookies
		// to `ToSession` function
		//
		// However, you can pass only the value of
		// ory_session_projectid cookie to the endpoint
		cookies = r.Header.Get("Cookie")

		// check if we have a session
		session, _, err := s.ory.FrontendApi.ToSession(r.Context()).Cookie(cookies).Execute()
		if (err != nil && session == nil) || (err == nil && !*session.Active) {
			// this will redirect the user to the managed Ory Login UI
			http.Redirect(w, r, "/.ory/self-service/login/browser", http.StatusSeeOther)
			return
		}

		ctx := withCookies(r.Context(), cookies)
		ctx = withSession(ctx, session)

		// continue to the requested page (in our case the Dashboard)
		next.ServeHTTP(w, r.WithContext(ctx))
		return
	}
}
