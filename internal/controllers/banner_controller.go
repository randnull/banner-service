package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/randnull/banner-service/internal/errors"
	"github.com/randnull/banner-service/internal/service"
	"github.com/randnull/banner-service/pkg/models"
)


type BannerHandlers struct {
	Service *service.BannerService
}


func NewHandler(BannerServ *service.BannerService) *BannerHandlers {
	return &BannerHandlers{
		Service: BannerServ,
	}
}


func CheckIsAdmin(ctx context.Context) bool {
	is_admin := ctx.Value("is_admin")

	return (is_admin == false)
}


func (handler *BannerHandlers) CreateBanner(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r_ctx := r.Context()

	if CheckIsAdmin(r_ctx) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var banner models.Banner

	err := json.NewDecoder(r.Body).Decode(&banner)

	if err != nil {
		sendJSONError(w, http.StatusBadRequest, errors.InvalidRequestBody)
		return
	}

	id, err := handler.Service.CreateBanner(&banner)
	
	if err != nil {
		if err == errors.BannerAlreadyExist {
			sendJSONError(w, http.StatusBadRequest, err)
		} else {
			sendJSONError(w, http.StatusInternalServerError, errors.InternalError)
		}
		return
	}

	ResponseIdResponse := models.IdResponse{BannerId: id}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ResponseIdResponse)
}


func (handler *BannerHandlers) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r_ctx := r.Context()

	if CheckIsAdmin(r_ctx) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)

	BannerId, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(w, errors.BadBannerIdParam.Error(), http.StatusBadRequest)
		return
	}

	err = handler.Service.DeleteBanner(BannerId)

	if err != nil {
		if err == errors.BannerNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
			sendJSONError(w, http.StatusInternalServerError, errors.InternalError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


func (handler *BannerHandlers) GetAllBanners(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r_ctx := r.Context()

	if CheckIsAdmin(r_ctx) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	query_params := r.URL.Query()

	feature_id, err := strconv.Atoi(query_params.Get("feature_id"))

	if err != nil {
		feature_id = -1
	}

	tag_id, err := strconv.Atoi(query_params.Get("tag_id"))

	if err != nil {
		tag_id = -1
	}

	limit, err := strconv.Atoi(query_params.Get("limit"))

	if err != nil {
		limit = -1
	}

	offset, err := strconv.Atoi(query_params.Get("offset"))

	if err != nil {
		offset = -1
	}

	banners, err := handler.Service.GetAllBanners(tag_id, feature_id, limit, offset)

	if err != nil {
		sendJSONError(w, http.StatusInternalServerError, errors.InternalError)
		return
	}

	JsonAnswer, err := json.Marshal(banners)

	if err != nil {
		sendJSONError(w, http.StatusInternalServerError, errors.InternalError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(JsonAnswer)
}


func (handler *BannerHandlers) GetBanner(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	query_params := r.URL.Query()

	tag_id, err := strconv.Atoi(query_params.Get("tag_id"))

	if err != nil {
		http.Error(w, errors.BadTagIdParam.Error(), http.StatusBadRequest)
		return
	}

	feature_id, err := strconv.Atoi(query_params.Get("feature_id"))

	if err != nil {
		http.Error(w, errors.BadFeatureIdParam.Error(), http.StatusBadRequest)
		return
	}

	use_last_revision, err := strconv.ParseBool(query_params.Get("use_last_version"))

	if err != nil {
		use_last_revision = false
	}
	
	content, err := handler.Service.GetBanner(tag_id, feature_id, use_last_revision)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	JsonAnswer, err := json.Marshal(content)

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(JsonAnswer)
}


func (handler *BannerHandlers) UpdateBanner(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r_ctx := r.Context()

	if CheckIsAdmin(r_ctx) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)

	BannerId, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(w, errors.BadBannerIdParam.Error(), http.StatusBadRequest)
		return
	}

	var new_banner models.UpdateBanner

	err = json.NewDecoder(r.Body).Decode(&new_banner)
	
	if err != nil {
		http.Error(w, errors.BadUpdateData.Error(), http.StatusBadRequest)
		return
	}

	err = handler.Service.UpdateBanner(BannerId, &new_banner)

	if err != nil {
		if err == errors.BannerAlreadyExist {
			sendJSONError(w, http.StatusBadRequest, err)
		} else {
			sendJSONError(w, http.StatusInternalServerError, errors.InternalError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
