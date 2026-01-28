package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

type adminContactsSnapshotResp struct {
	Email     string                  `json:"email"`
	Phones    []entity.ContactNumber  `json:"phones"`
	Addresses []entity.ContactAddress `json:"addresses"`
}

type adminSetEmailReq struct {
	Email string `json:"email"`
}

type adminUpsertPhoneReq struct {
	ID        string `json:"id"`
	Phone     string `json:"phone"`
	Label     string `json:"label"`
	SortOrder int    `json:"sortOrder"`
}

type adminUpsertAddressReq struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Address   string  `json:"address"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	SortOrder int     `json:"sortOrder"`
}

type adminReorderReq struct {
	IDs []string `json:"ids"`
}

// GET /admin/contacts
func (h *Handler) AdminGetContacts(c *gin.Context) {
	email, err := h.adminContactsService.GetContactsEmail(c.Request.Context())
	if err != nil {
		h.logger.Error("AdminGetContacts: email", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get email"})
		return
	}

	phones, err := h.adminContactsService.ListContactNumbers(c.Request.Context())
	if err != nil {
		h.logger.Error("AdminGetContacts: phones", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get phones"})
		return
	}

	addrs, err := h.adminContactsService.ListContactAddresses(c.Request.Context())
	if err != nil {
		h.logger.Error("AdminGetContacts: addresses", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get addresses"})
		return
	}

	c.JSON(http.StatusOK, adminContactsSnapshotResp{
		Email:     email,
		Phones:    phones,
		Addresses: addrs,
	})
}

// PUT /admin/contacts/email
func (h *Handler) AdminSetContactsEmail(c *gin.Context) {
	var req adminSetEmailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	if err := h.adminContactsService.SetContactsEmail(c.Request.Context(), req.Email); err != nil {
		h.logger.Error("AdminSetContactsEmail", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// GET /admin/contacts/phones
func (h *Handler) AdminListContactPhones(c *gin.Context) {
	items, err := h.adminContactsService.ListContactNumbers(c.Request.Context())
	if err != nil {
		h.logger.Error("AdminListContactPhones", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list phones"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// PUT /admin/contacts/phones
func (h *Handler) AdminUpsertContactPhone(c *gin.Context) {
	var req adminUpsertPhoneReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	req.ID = strings.TrimSpace(req.ID)
	if req.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	err := h.adminContactsService.UpsertContactNumber(c.Request.Context(), entity.ContactNumber{
		ID:        req.ID,
		Phone:     req.Phone,
		Label:     req.Label,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		h.logger.Error("AdminUpsertContactPhone", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upsert phone"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// DELETE /admin/contacts/phones/:id
func (h *Handler) AdminDeleteContactPhone(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	if err := h.adminContactsService.DeleteContactNumber(c.Request.Context(), id); err != nil {
		h.logger.Error("AdminDeleteContactPhone", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete phone"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// PUT /admin/contacts/phones/reorder
func (h *Handler) AdminReorderContactPhones(c *gin.Context) {
	var req adminReorderReq
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json or empty ids"})
		return
	}

	if err := h.adminContactsService.SetContactNumbersOrder(c.Request.Context(), req.IDs); err != nil {
		h.logger.Error("AdminReorderContactPhones", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reorder phones"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// GET /admin/contacts/addresses
func (h *Handler) AdminListContactAddresses(c *gin.Context) {
	items, err := h.adminContactsService.ListContactAddresses(c.Request.Context())
	if err != nil {
		h.logger.Error("AdminListContactAddresses", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list addresses"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// PUT /admin/contacts/addresses
func (h *Handler) AdminUpsertContactAddress(c *gin.Context) {
	var req adminUpsertAddressReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	req.ID = strings.TrimSpace(req.ID)
	if req.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	err := h.adminContactsService.UpsertContactAddress(c.Request.Context(), entity.ContactAddress{
		ID:        req.ID,
		Title:     req.Title,
		Address:   req.Address,
		Lat:       req.Lat,
		Lon:       req.Lon,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		h.logger.Error("AdminUpsertContactAddress", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upsert address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// DELETE /admin/contacts/addresses/:id
func (h *Handler) AdminDeleteContactAddress(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	if err := h.adminContactsService.DeleteContactAddress(c.Request.Context(), id); err != nil {
		h.logger.Error("AdminDeleteContactAddress", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// PUT /admin/contacts/addresses/reorder
func (h *Handler) AdminReorderContactAddresses(c *gin.Context) {
	var req adminReorderReq
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json or empty ids"})
		return
	}

	if err := h.adminContactsService.SetContactAddressesOrder(c.Request.Context(), req.IDs); err != nil {
		h.logger.Error("AdminReorderContactAddresses", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reorder addresses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
