package smartolt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/enxservices/smartolt/internal/types"
)

type SmartOLTClient struct {
	client  *http.Client
	baseURL string
}

func NewSmartOLTClient(token, baseURL string) *SmartOLTClient {
	return &SmartOLTClient{
		client: &http.Client{
			Transport: &TransportWithToken{
				Token: token,
			},
		},
		baseURL: baseURL,
	}
}

func (c *SmartOLTClient) doRequest(req *http.Request, respBody interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var respErr types.ResponseError
		data, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(data, &respErr); err != nil {
			return fmt.Errorf("%s", respErr.Error)
		}

		return fmt.Errorf("%v", resp.StatusCode)
	}

	if respBody == nil {
		return nil
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := json.Unmarshal(data, respBody); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (c *SmartOLTClient) GetOnuDetails(ID string) (*types.OnuDetails, error) {
	url := fmt.Sprintf("%s%s/%s", c.baseURL, types.GETONUDETAILS, ID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	var details types.OnuDetails
	if err := c.doRequest(req, &details); err != nil {
		return nil, err
	}
	return &details, nil
}

func (c *SmartOLTClient) GetOnuSignal(ID string) (*types.StatusSignal, error) {
	url := fmt.Sprintf("%s%s/%s", c.baseURL, types.GETONUSIGNAL, ID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	var signal types.StatusSignal
	if err := c.doRequest(req, &signal); err != nil {
		return nil, err
	}
	return &signal, nil
}

func (c *SmartOLTClient) GetSpeedProfiles() ([]types.SpeedProfile, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, types.GETSPEEDPROFILE)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	var resp types.Response[types.SpeedProfile]
	if err := c.doRequest(req, &resp); err != nil {
		return nil, err
	}

	return resp.Response, nil
}

func (c *SmartOLTClient) UpdateSpeedProfile(ID, downloadProfile, uploadProfile string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if err := writer.WriteField("upload_speed_profile_name", uploadProfile); err != nil {
		return fmt.Errorf("%w", err)
	}
	if err := writer.WriteField("download_speed_profile_name", downloadProfile); err != nil {
		return fmt.Errorf("%w", err)
	}
	writer.Close()

	url := fmt.Sprintf("%s%s/%s", c.baseURL, types.UPDATESPEEDPROFILE, ID)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return c.doRequest(req, nil)
}

func (c *SmartOLTClient) RebootOnu(ID string) error {
	url := fmt.Sprintf("%s%s/%s", c.baseURL, types.REBOOTONU, ID)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return c.doRequest(req, nil)
}

func (c *SmartOLTClient) DisableOnu(ID string) error {
	url := fmt.Sprintf("%s%s/%s", c.baseURL, types.DISABLEONU, ID)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return c.doRequest(req, nil)
}

func (c *SmartOLTClient) EnableOnu(ID string) error {
	url := fmt.Sprintf("%s%s/%s", c.baseURL, types.ENABLEONU, ID)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return c.doRequest(req, nil)
}
