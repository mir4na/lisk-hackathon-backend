package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/receiv3/backend/internal/config"
	"github.com/receiv3/backend/internal/utils"
)

type PinataService struct {
	apiKey     string
	secretKey  string
	jwt        string
	gatewayURL string
}

type PinataResponse struct {
	IpfsHash    string `json:"IpfsHash"`
	PinSize     int    `json:"PinSize"`
	Timestamp   string `json:"Timestamp"`
	IsDuplicate bool   `json:"isDuplicate,omitempty"`
}

type PinataMetadata struct {
	Name      string            `json:"name"`
	KeyValues map[string]string `json:"keyvalues,omitempty"`
}

func NewPinataService(cfg *config.Config) *PinataService {
	return &PinataService{
		apiKey:     cfg.PinataAPIKey,
		secretKey:  cfg.PinataSecretKey,
		jwt:        cfg.PinataJWT,
		gatewayURL: cfg.PinataGatewayURL,
	}
}

func (s *PinataService) UploadFile(fileData []byte, fileName string, metadata map[string]string) (*PinataResponse, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, "", err
	}
	if _, err := part.Write(fileData); err != nil {
		return nil, "", err
	}

	// Add pinata metadata
	pinataMetadata := PinataMetadata{
		Name:      fileName,
		KeyValues: metadata,
	}
	metadataBytes, _ := json.Marshal(pinataMetadata)
	if err := writer.WriteField("pinataMetadata", string(metadataBytes)); err != nil {
		return nil, "", err
	}

	if err := writer.Close(); err != nil {
		return nil, "", err
	}

	req, err := http.NewRequest("POST", "https://api.pinata.cloud/pinning/pinFileToIPFS", body)
	if err != nil {
		return nil, "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	if s.jwt != "" {
		req.Header.Set("Authorization", "Bearer "+s.jwt)
	} else {
		req.Header.Set("pinata_api_key", s.apiKey)
		req.Header.Set("pinata_secret_api_key", s.secretKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("pinata upload failed: %s", string(bodyBytes))
	}

	var pinataResp PinataResponse
	if err := json.NewDecoder(resp.Body).Decode(&pinataResp); err != nil {
		return nil, "", err
	}

	fileHash := utils.SHA256Hash(fileData)

	return &pinataResp, fileHash, nil
}

func (s *PinataService) UploadJSON(data interface{}, name string) (*PinataResponse, error) {
	body := map[string]interface{}{
		"pinataContent": data,
		"pinataMetadata": map[string]string{
			"name": name,
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.pinata.cloud/pinning/pinJSONToIPFS", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if s.jwt != "" {
		req.Header.Set("Authorization", "Bearer "+s.jwt)
	} else {
		req.Header.Set("pinata_api_key", s.apiKey)
		req.Header.Set("pinata_secret_api_key", s.secretKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("pinata upload failed: %s", string(bodyBytes))
	}

	var pinataResp PinataResponse
	if err := json.NewDecoder(resp.Body).Decode(&pinataResp); err != nil {
		return nil, err
	}

	return &pinataResp, nil
}

func (s *PinataService) GetIPFSURL(hash string) string {
	return s.gatewayURL + hash
}

// NFTMetadata represents the metadata for an Invoice NFT
type NFTMetadata struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Image       string                 `json:"image,omitempty"`
	ExternalURL string                 `json:"external_url,omitempty"`
	Attributes  []NFTAttribute         `json:"attributes"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
}

type NFTAttribute struct {
	TraitType   string      `json:"trait_type"`
	Value       interface{} `json:"value"`
	DisplayType string      `json:"display_type,omitempty"`
}

func (s *PinataService) UploadNFTMetadata(metadata *NFTMetadata) (string, error) {
	resp, err := s.UploadJSON(metadata, metadata.Name)
	if err != nil {
		return "", err
	}
	return s.GetIPFSURL(resp.IpfsHash), nil
}
