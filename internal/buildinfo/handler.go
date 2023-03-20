package buildinfo

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// Handler returns an HTTP handler for version information.
func Handler(bi BuildInfo) gin.HandlerFunc {
	var body []byte

	return func(ctx *gin.Context) {
		if body == nil {
			var err error

			body, err = json.Marshal(bi)
			if err != nil {
				panic(errors.Wrap(err, "failed to render version information"))
			}
		}

		ctx.Data(http.StatusOK, "application/json", body)
	}
}
