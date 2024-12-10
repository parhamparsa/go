package v1

import (
	"app/internal/domain/entity"
	"app/pkg/util"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func (restApiV1 *RestApiV1) createAccountHandler(w http.ResponseWriter, r *http.Request) {
	var account entity.Account
	body, err := io.ReadAll(r.Body)
	if err != nil {
		_ = util.WriteJSON(w, http.StatusBadRequest, err)
		return
	}

	if err := json.Unmarshal(body, &account); err != nil {
		_ = util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if _, err := restApiV1.app.GetAccountService().Create(r.Context(), &account); err != nil {
		_ = util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	responseErr := util.WriteJSON(w, http.StatusCreated, map[string]string{"data": "account created successfully"})
	if responseErr != nil {
		_ = util.WriteJSON(w, http.StatusBadRequest, responseErr)
		return
	}
}

func (restApiV1 *RestApiV1) updateAccountHandler(w http.ResponseWriter, r *http.Request) {
	var account entity.Account
	body, err := io.ReadAll(r.Body)
	if err != nil {
		_ = util.WriteJSON(w, http.StatusBadRequest, err)
		return
	}

	if err := json.Unmarshal(body, &account); err != nil {
		_ = util.WriteJSON(w, http.StatusBadRequest, err)
		return
	}

	if err := restApiV1.app.GetAccountService().UpdateAccountBalance(r.Context(), uint32(account.UserID), account.Balance); err != nil {
		_ = util.WriteJSON(w, http.StatusBadRequest, err)
		return
	}

	responseErr := util.WriteJSON(w, http.StatusCreated, map[string]string{"data": "account update successfully"})
	if responseErr != nil {
		_ = util.WriteJSON(w, http.StatusBadRequest, responseErr)
		return
	}
}

func (restApiV1 *RestApiV1) findAccountHandler(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.URL.Query().Get("user_id")
	userId, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		_ = util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	account, err := restApiV1.app.GetAccountService().Find(r.Context(), uint32(userId))
	if err != nil {
		_ = util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	responseErr := util.WriteJSON(w, http.StatusCreated, account)
	if responseErr != nil {
		_ = util.WriteJSON(w, http.StatusBadRequest, responseErr)
		return
	}
}
