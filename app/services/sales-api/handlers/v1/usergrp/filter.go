package usergrp

import (
	"github.com/yourusername/basic-a/business/core/user"
	"net/http"
)

func getFilter(r *http.Request) (user.QueryFilter, error) {
	values := r.URL.Query()

	var filter user.QueryFilter
	filter.ByID(values.Get("id"))
	filter.ByName(values.Get("name"))
	filter.ByEmail(values.Get("email"))

	return filter, nil
}
