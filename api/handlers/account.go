package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"

	"basic-go-project/src/entities/dtos"
)

// TODO: move validator to PKG
//var validate *validator.Validate

type AccountUseCase interface {
	Create(account *dtos.Account) error
	GetAll(limit, offset int) ([]dtos.Account, error)
	GetByID(id string) (dtos.Account, error)
	Update(updateAccountData *dtos.UpdateAccountRequest) error
	Delete(id string) error
}

type AccountHandler struct {
	accountUseCase AccountUseCase
	log            *zerolog.Logger
}

func NewAccountHandler(accountUseCase AccountUseCase, log *zerolog.Logger) AccountHandler {
	return AccountHandler{
		accountUseCase: accountUseCase,
		log:            log,
	}
}

// TODO: add login pass validation
type CreateAccountRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (h AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	requestData := new(CreateAccountRequest)
	if err := DecodeBody(r.Body, requestData); err != nil {
		RespondErr(w, h.log, fmt.Errorf("DecodeBody() failed: %w", err), http.StatusBadRequest)

		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(requestData); err != nil {
		RespondErr(w, h.log, fmt.Errorf("validate.Struct() failed: %w", err), http.StatusBadRequest)

		return
	}

	account := &dtos.Account{
		Login:    requestData.Login,
		Password: requestData.Password,
		IsActive: true,
	}

	if err := h.accountUseCase.Create(account); err != nil {
		RespondErr(w, h.log, fmt.Errorf("accountUseCase.Create(): %w", err), http.StatusInternalServerError)

		return
	}

	RespondStatusOk(w, h.log)
}

type GetAllAccountsRequest struct {
	Limit  int `json:"limit" validate:"gte=1,lte=500"`
	Offset int `json:"offset" validate:"gte=0"`
}

func (h AccountHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	requestData := new(GetAllAccountsRequest)
	if err := DecodeBody(r.Body, requestData); err != nil {
		RespondErr(w, h.log, fmt.Errorf("DecodeBody() failed: %w", err), http.StatusBadRequest)

		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(requestData); err != nil {
		RespondErr(w, h.log, fmt.Errorf("validate.Struct() failed: %w", err), http.StatusBadRequest)

		return
	}

	accounts, err := h.accountUseCase.GetAll(requestData.Limit, requestData.Offset)
	if err != nil {
		RespondErr(w, h.log, fmt.Errorf("accountUseCase.GetAll: %w", err), http.StatusInternalServerError)

		return
	}

	Respond(w, h.log, accounts)
}

func (h AccountHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		RespondErr(w, h.log, fmt.Errorf("id is required"), http.StatusBadRequest)

		return
	}

	account, err := h.accountUseCase.GetByID(id)
	if err != nil {
		RespondErr(w, h.log, fmt.Errorf("accountUseCase.GetByID(): %w", err), http.StatusInternalServerError)

		return
	}

	Respond(w, h.log, account)
}

// TODO: add login pass validation
type UpdateAccountRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	IsActive bool   `json:"isActive" validate:"required"`
}

func (h AccountHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		RespondErr(w, h.log, fmt.Errorf("id is required"), http.StatusBadRequest)

		return
	}

	requestData := new(dtos.UpdateAccountRequest)
	if err := DecodeBody(r.Body, requestData); err != nil {
		RespondErr(w, h.log, fmt.Errorf("DecodeBody() failed: %w", err), http.StatusBadRequest)

		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(requestData); err != nil {
		RespondErr(w, h.log, fmt.Errorf("validate.Struct() failed: %w", err), http.StatusBadRequest)

		return
	}

	requestData.ID = id

	if err := h.accountUseCase.Update(requestData); err != nil {
		RespondErr(w, h.log, fmt.Errorf("accountUseCase.Update(): %w", err), http.StatusInternalServerError)

		return
	}

	RespondStatusOk(w, h.log)
}

func (h AccountHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		RespondErr(w, h.log, fmt.Errorf("id is required"), http.StatusBadRequest)

		return
	}

	if err := h.accountUseCase.Delete(id); err != nil {
		RespondErr(w, h.log, fmt.Errorf("accountUseCase.Delete(): %w", err), http.StatusInternalServerError)

		return
	}

	RespondStatusOk(w, h.log)
}
