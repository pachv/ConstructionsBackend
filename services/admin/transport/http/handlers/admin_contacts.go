package handlers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Admin contacts proxy to constructions service.
// Base URL is the same as in sections.go (PublicAPIBaseURLC).

const (
	PublicAPIBaseURLC = "http://constructions_service:8080"
)

func (h *Handler) ProxyAdminContactsGet(c *gin.Context) {
	h.proxyJSON(c, http.MethodGet, PublicAPIBaseURLC+"/admin/contacts", nil)
}

func (h *Handler) ProxyAdminContactsSetEmail(c *gin.Context) {
	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cant read body"})
		return
	}
	h.proxyJSON(c, http.MethodPut, PublicAPIBaseURLC+"/admin/contacts/email", raw)
}

// ---- phones ----

func (h *Handler) ProxyAdminContactsListPhones(c *gin.Context) {
	h.proxyJSON(c, http.MethodGet, PublicAPIBaseURLC+"/admin/contacts/phones", nil)
}

func (h *Handler) ProxyAdminContactsUpsertPhone(c *gin.Context) {
	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cant read body"})
		return
	}
	h.proxyJSON(c, http.MethodPut, PublicAPIBaseURLC+"/admin/contacts/phones", raw)
}

func (h *Handler) ProxyAdminContactsDeletePhone(c *gin.Context) {
	id := c.Param("id")
	h.proxyJSON(c, http.MethodDelete, PublicAPIBaseURLC+"/admin/contacts/phones/"+id, nil)
}

func (h *Handler) ProxyAdminContactsReorderPhones(c *gin.Context) {
	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cant read body"})
		return
	}
	h.proxyJSON(c, http.MethodPut, PublicAPIBaseURLC+"/admin/contacts/phones/reorder", raw)
}

// ---- addresses ----

func (h *Handler) ProxyAdminContactsListAddresses(c *gin.Context) {
	h.proxyJSON(c, http.MethodGet, PublicAPIBaseURLC+"/admin/contacts/addresses", nil)
}

func (h *Handler) ProxyAdminContactsUpsertAddress(c *gin.Context) {
	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cant read body"})
		return
	}
	h.proxyJSON(c, http.MethodPut, PublicAPIBaseURLC+"/admin/contacts/addresses", raw)
}

func (h *Handler) ProxyAdminContactsDeleteAddress(c *gin.Context) {
	id := c.Param("id")
	h.proxyJSON(c, http.MethodDelete, PublicAPIBaseURLC+"/admin/contacts/addresses/"+id, nil)
}

func (h *Handler) ProxyAdminContactsReorderAddresses(c *gin.Context) {
	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cant read body"})
		return
	}
	h.proxyJSON(c, http.MethodPut, PublicAPIBaseURLC+"/admin/contacts/addresses/reorder", raw)
}
