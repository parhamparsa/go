package v1

import (
	"app/internal/domain/entity"
	"app/pkg/util"
	"encoding/json"
	"io"
	"net/http"
)

func (restApiV1 *RestApiV1) transferHandler(w http.ResponseWriter, r *http.Request) {
	var transfer entity.Transfer
	body, err := io.ReadAll(r.Body)
	if err != nil {
		_ = util.WriteJSON(w, http.StatusInternalServerError, map[string]string{"data": "Failed to read body"})
		return
	}

	if err := json.Unmarshal(body, &transfer); err != nil {
		_ = util.WriteJSON(w, http.StatusBadRequest, err)
		return
	}

	if err := restApiV1.app.GetTransferService().Transfer(r.Context(), &transfer); err != nil {
		_ = util.WriteJSON(w, http.StatusBadRequest, err)
		return
	}

	responseErr := util.WriteJSON(w, http.StatusCreated, map[string]string{"data": "transfer successfully"})
	if responseErr != nil {
		_ = util.WriteJSON(w, http.StatusBadRequest, responseErr)
		return
	}
}
