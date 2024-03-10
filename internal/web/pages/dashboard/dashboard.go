package dashboard

import (
	"fmt"
	"net/http"
	"tasklify/internal/handlers"
	"tasklify/internal/web/pages"
)

func Dashboard(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	c := pages.Index(fmt.Sprint(params.UserID))
	return pages.Layout(c, "Tasklify").Render(r.Context(), w)
}
