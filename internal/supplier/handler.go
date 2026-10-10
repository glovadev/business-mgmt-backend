package supplier

import (
	"encoding/json"
	"net/http"
	"strconv"

	"business-mgmt-backend/internal/models"
	"business-mgmt-backend/pkg/utils"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateSupplier(w http.ResponseWriter, r *http.Request) {
	var req models.CreateSupplierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	supplier, err := h.service.CreateSupplier(r.Context(), req)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, supplier)
}

func (h *Handler) ListSuppliers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	
	page, _ := strconv.Atoi(query.Get("page"))
	pageSize, _ := strconv.Atoi(query.Get("page_size"))

	filters := make(map[string]string)
	if status := query.Get("status"); status != "" {
		filters["status"] = status
	}
	if q := query.Get("q"); q != "" {
		filters["q"] = q
	}

	resp, err := h.service.ListSuppliers(r.Context(), page, pageSize, filters)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch suppliers")
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetSupplier(w http.ResponseWriter, r *http.Request) {
	supplierIDStr := r.PathValue("supplier_id")
	supplierID, err := uuid.Parse(supplierIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid supplier ID format")
		return
	}

	supplier, err := h.service.GetSupplier(r.Context(), supplierID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, supplier)
}

func (h *Handler) UpdateSupplier(w http.ResponseWriter, r *http.Request) {
	supplierIDStr := r.PathValue("supplier_id")
	supplierID, err := uuid.Parse(supplierIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid supplier ID format")
		return
	}

	var req models.UpdateSupplierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.UpdateSupplier(r.Context(), supplierID, req); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to update supplier")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ActivateSupplier(w http.ResponseWriter, r *http.Request) {
	supplierIDStr := r.PathValue("supplier_id")
	supplierID, err := uuid.Parse(supplierIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid supplier ID format")
		return
	}

	if err := h.service.ActivateSupplier(r.Context(), supplierID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to activate supplier")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeactivateSupplier(w http.ResponseWriter, r *http.Request) {
	supplierIDStr := r.PathValue("supplier_id")
	supplierID, err := uuid.Parse(supplierIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid supplier ID format")
		return
	}

	if err := h.service.DeactivateSupplier(r.Context(), supplierID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to deactivate supplier")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
