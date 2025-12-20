package handlers

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"

	"github.com/gin-gonic/gin"
)

const constructionsBaseURL = "http://constructions_service:8080"

func (h *Handler) DeleteReviewProxy(c *gin.Context) {
	id := c.Param("id")

	req, err := http.NewRequest(http.MethodDelete, constructionsBaseURL+"/admin/reviews/"+id, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to contact constructions service"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

func (h *Handler) BulkUpdateReviewsProxy(c *gin.Context) {
	// читаем json как bytes
	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot read body"})
		return
	}

	req, err := http.NewRequest(http.MethodPut, constructionsBaseURL+"/admin/reviews/bulk", bytes.NewReader(raw))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to contact constructions service"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

func (h *Handler) CreateReviewProxy(c *gin.Context) {
	file, _ := c.FormFile("photo") // optional

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// поля
	_ = writer.WriteField("name", c.PostForm("name"))
	_ = writer.WriteField("position", c.PostForm("position"))
	_ = writer.WriteField("text", c.PostForm("text"))
	_ = writer.WriteField("rating", c.PostForm("rating"))
	_ = writer.WriteField("consent", c.PostForm("consent"))
	_ = writer.WriteField("canPublish", c.PostForm("canPublish"))

	if file != nil {
		src, err := file.Open()
		if err != nil {
			_ = writer.Close()
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to open file"})
			return
		}
		defer src.Close()

		// ✅ ВАЖНО: сохраняем Content-Type
		ct := file.Header.Get("Content-Type")
		if ct == "" {
			ct = "application/octet-stream"
		}

		hdr := textproto.MIMEHeader{}
		hdr.Set("Content-Disposition", `form-data; name="photo"; filename="`+file.Filename+`"`)
		hdr.Set("Content-Type", ct)

		part, err := writer.CreatePart(hdr)
		if err != nil {
			_ = writer.Close()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create form part"})
			return
		}

		if _, err := io.Copy(part, src); err != nil {
			_ = writer.Close()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to copy file"})
			return
		}
	}

	_ = writer.Close()

	req, err := http.NewRequest(http.MethodPost, constructionsBaseURL+"/admin/reviews", &buf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to contact constructions service"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}
