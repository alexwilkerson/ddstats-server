package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) clientMessage(w http.ResponseWriter, status int, message string) {
	data := struct {
		Message string `json:"message"`
	}{message}

	js, err := json.Marshal(data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func validVersion(version string) (bool, error) {
	var vMajor, vMinor, vPatch, ovMajor, ovMinor, ovPatch int
	_, err := fmt.Sscanf(version, "%d.%d.%d", &vMajor, &vMinor, &vPatch)
	if err != nil {
		return false, err
	}
	_, err = fmt.Sscanf(oldestValidClientVersion, "%d.%d.%d", &ovMajor, &ovMinor, &ovPatch)
	if err != nil {
		return false, err
	}
	if vMajor > ovMajor ||
		(vMajor == ovMajor && vMinor > ovMinor) ||
		(vMajor == ovMajor && vMinor == ovMinor && vPatch >= ovPatch) {
		return true, nil
	}
	return false, nil
}

func updateAvailable(version string) (bool, error) {
	var vMajor, vMinor, vPatch, cvMajor, cvMinor, cvPatch int
	_, err := fmt.Sscanf(version, "%d.%d.%d", &vMajor, &vMinor, &vPatch)
	if err != nil {
		return false, err
	}
	_, err = fmt.Sscanf(currentClientVersion, "%d.%d.%d", &cvMajor, &cvMinor, &cvPatch)
	if err != nil {
		return false, err
	}
	if cvMajor > vMajor ||
		(cvMajor == vMajor && cvMinor > vMinor) ||
		(cvMajor == vMajor && cvMinor == vMinor && cvPatch > vPatch) {
		return true, nil
	}
	return false, nil
}
