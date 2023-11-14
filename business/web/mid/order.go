package mid

import (
	"errors"
	"github.com/yourusername/basic-a/business/data/order"
	"net/http"
	"strings"
)

// GetOrderBy constructs a order.By value by parsing a string in the form
// of "field,direction".
func GetOrderBy(r *http.Request, defaultOrder order.By) (order.By, error) {
	v := r.URL.Query().Get("orderBy")

	if v == "" {
		return defaultOrder, nil
	}

	orderParts := strings.Split(v, ",")

	var by order.By
	switch len(orderParts) {
	case 1:
		by = order.NewBy(strings.Trim(orderParts[0], " "), order.ASC)
	case 2:
		by = order.NewBy(strings.Trim(orderParts[0], " "), strings.Trim(orderParts[1], " "))
	default:
		return order.By{}, errors.New("invalid ordering information")
	}

	return by, nil
}
