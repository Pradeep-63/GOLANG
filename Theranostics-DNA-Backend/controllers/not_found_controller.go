// controllers/not_found_controller.go
package controllers

import (
    "theransticslabs/m/utils"
    "net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
    utils.JSONResponse(w, http.StatusNotFound, false, utils.MsgEndpointNotFound, nil)
}