// controllers/welcome_controller.go
package controllers

import (
    "theransticslabs/m/utils"
    "net/http"
)

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
    utils.JSONResponse(w, http.StatusOK, true, utils.MsgWelcome, nil)
}
