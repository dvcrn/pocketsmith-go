package pocketsmith

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ContentTypeMeta struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Extension   string `json:"extension"`
}

type AttachmentVariants struct {
	ThumbURL string `json:"thumb_url"`
	LargeURL string `json:"large_url"`
}

type Attachment struct {
	ID              int64              `json:"id"`
	Title           string             `json:"title"`
	FileName        string             `json:"file_name"`
	Type            string             `json:"type"`
	ContentType     string             `json:"content_type"`
	ContentTypeMeta ContentTypeMeta    `json:"content_type_meta"`
	OriginalURL     string             `json:"original_url"`
	Variants        AttachmentVariants `json:"variants"`
	CreatedAt       string             `json:"created_at"`
	UpdatedAt       string             `json:"updated_at"`
}

// ListAttachments retrieves all attachments for a given user
func (c *Client) ListAttachments(userID int, unassigned bool) ([]*Attachment, error) {
	baseURL := fmt.Sprintf("https://api.pocketsmith.com/v2/users/%d/attachments", userID)
	url := baseURL
	if unassigned {
		url = fmt.Sprintf("%s?unassigned=1", baseURL)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var attachments []*Attachment
	if err := c.doAndDecode(req, &attachments); err != nil {
		return nil, err
	}

	return attachments, nil
}

type CreateAttachment struct {
	Title    string `json:"title"`
	FileName string `json:"file_name"`
	FileData string `json:"file_data"`
}

// CreateAttachment creates a new attachment for the specified user
func (c *Client) CreateAttachment(userID int, attachment *CreateAttachment) (*Attachment, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/users/%d/attachments", userID)

	payload, err := json.Marshal(attachment)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/json")

	var createdAttachment Attachment
	if err := c.doAndDecode(req, &createdAttachment); err != nil {
		return nil, err
	}

	return &createdAttachment, nil
}

// ListTransactionAttachments retrieves all attachments for a specific transaction
func (c *Client) ListTransactionAttachments(transactionID int64) ([]*Attachment, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/transactions/%d/attachments", transactionID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var attachments []*Attachment
	if err := c.doAndDecode(req, &attachments); err != nil {
		return nil, err
	}

	return attachments, nil
}

type AttachToTransaction struct {
	AttachmentID int64 `json:"attachment_id"`
}

// AttachToTransaction attaches an existing attachment to a transaction
func (c *Client) AttachToTransaction(transactionID int64, attachmentID int64) error {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/transactions/%d/attachments", transactionID)

	payload := AttachToTransaction{
		AttachmentID: attachmentID,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Add("content-type", "application/json")

	return c.doAndDecode(req, nil)
}

// UnassignAttachment removes an attachment from a transaction
func (c *Client) UnassignAttachment(transactionID int64, attachmentID int64) error {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/transactions/%d/attachments/%d", transactionID, attachmentID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	return c.doAndDecode(req, nil)
}

// GetAttachment retrieves a single attachment by its ID
func (c *Client) GetAttachment(attachmentID int64) (*Attachment, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/attachments/%d", attachmentID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var attachment Attachment
	if err := c.doAndDecode(req, &attachment); err != nil {
		return nil, err
	}

	return &attachment, nil
}

type UpdateAttachment struct {
	Title string `json:"title"`
}

// UpdateAttachment updates an existing attachment's title
func (c *Client) UpdateAttachment(attachmentID int64, update *UpdateAttachment) (*Attachment, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/attachments/%d", attachmentID)

	payload, err := json.Marshal(update)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/json")

	var updatedAttachment Attachment
	if err := c.doAndDecode(req, &updatedAttachment); err != nil {
		return nil, err
	}

	return &updatedAttachment, nil
}
