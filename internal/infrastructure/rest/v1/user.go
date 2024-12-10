package v1

import (
	"app/internal/domain/entity"
	"app/pkg/util"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func (restApiV1 *RestApiV1) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	body, err := io.ReadAll(r.Body)
	if err != nil {
		_ = util.WriteJSON(w, http.StatusBadRequest, err)
		return
	}

	if err := json.Unmarshal(body, &user); err != nil {
		_ = util.WriteJSON(w, http.StatusBadRequest, err)
		return
	}
	if err := restApiV1.app.GetUserService().Create(r.Context(), &user); err != nil {
		_ = util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	responseErr := util.WriteJSON(w, http.StatusCreated, map[string]string{"data": "user created successfully"})
	if responseErr != nil {
		_ = util.WriteJSON(w, http.StatusBadRequest, responseErr)
		return
	}
}

func (restApiV1 *RestApiV1) findUserHandler(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.URL.Query().Get("id")
	userId, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		_ = util.WriteJSON(w, http.StatusBadRequest, err)
		return
	}
	user, err := restApiV1.app.GetUserService().Find(r.Context(), uint32(userId))
	if err != nil {
		_ = util.WriteJSON(w, http.StatusBadRequest, err)
		return
	}

	responseErr := util.WriteJSON(w, http.StatusCreated, user)
	if responseErr != nil {
		_ = util.WriteJSON(w, http.StatusBadRequest, responseErr)
	}
}
