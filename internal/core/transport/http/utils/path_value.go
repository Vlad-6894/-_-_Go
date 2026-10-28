package core_http_utils

import (
	"fmt"
	"net/http"
	core_errors "security/internal/core/errors"
)

func GetStringPathValue(r *http.Request, key string) (string, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return "", fmt.Errorf("no key=%s in path values: %w", key, core_errors.ErrInvalidArgument)
	}

	return pathValue, nil
}
