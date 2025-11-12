package smartolt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/enxservices/smartolt/internal/types"
)

type ConnectionDetails struct {
	OltID            string
	PonType          string
	GponChannel      string
	Board            string
	Port             string
	SN               string
	VLAN             string
	OnuType          string
	Zone             string
	ODB              string
	Name             string
	AddressOrComment string
	OnuMode          string
	OnuExternalID    string
}

type Client interface {
	GetOnuDetails(ID string) (*types.OnuDetails, error)
	GetOnuSignal(ID string) (*types.StatusSignal, error)
	GetSpeedProfiles() ([]types.SpeedProfile, error)
	GetAllOnusDetails() ([]types.OnuListItem, error)
	GetOdbs() ([]types.ODB, error)
	GetODBAvailability() ([]types.ODBAvailability, error)
	UpdateSpeedProfile(ID, downloadProfile, uploadProfile string) error
	RebootOnu(ID string) error
	DisableOnu(ID string) error
	EnableOnu(ID string) error
	DiscoverOnuNeededReboot() ([]string, error)
	AuthorizeConnection(connectionDetails ConnectionDetails) error
	UnconfiguredOnusForOlt(oltID string) ([]types.UnconfiguredOnu, error)
}

type client struct {
	client  *http.Client
	baseURL string
}

func NewSmartOLTClient(token, baseURL string) Client {
	return &client{
		client: &http.Client{
			Transport: &TransportWithToken{
				Token: token,
			},
		},
		baseURL: baseURL,
	}
}

func (c *client) doRequest(req *http.Request, respBody interface{}) error {
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

func (c *client) GetOnuDetails(ID string) (*types.OnuDetails, error) {
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

func (c *client) GetOnuSignal(ID string) (*types.StatusSignal, error) {
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

func (c *client) GetSpeedProfiles() ([]types.SpeedProfile, error) {
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

func (c *client) GetAllOnusDetails() ([]types.OnuListItem, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, types.GETALLONUSDETAILS)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	var resp types.OnuListResponse
	if err := c.doRequest(req, &resp); err != nil {
		return nil, err
	}

	return resp.Onus, nil
}

func (c *client) GetOdbs() ([]types.ODB, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, types.GETODBS)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	var resp types.ODBListResponse
	if err := c.doRequest(req, &resp); err != nil {
		return nil, err
	}
	return resp.Response, nil
}

func (c *client) GetODBAvailability() ([]types.ODBAvailability, error) {
	obds, err := c.GetOdbs()
	if err != nil {
		return nil, err
	}

	onus, err := c.GetAllOnusDetails()
	if err != nil {
		return nil, err
	}

	availability := CalculateODBAvailability(obds, onus)
	return availability, nil
}

func (c *client) UpdateSpeedProfile(ID, downloadProfile, uploadProfile string) error {
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

func (c *client) RebootOnu(ID string) error {
	url := fmt.Sprintf("%s%s/%s", c.baseURL, types.REBOOTONU, ID)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return c.doRequest(req, nil)
}

func (c *client) DisableOnu(ID string) error {
	url := fmt.Sprintf("%s%s/%s", c.baseURL, types.DISABLEONU, ID)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return c.doRequest(req, nil)
}

func (c *client) EnableOnu(ID string) error {
	url := fmt.Sprintf("%s%s%s", c.baseURL, types.ENABLEONU, ID)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return c.doRequest(req, nil)
}

func (c *client) DiscoverOnuNeededReboot() ([]string, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, types.STATUSESONUS)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	var resp types.Response[types.OnuStatus]
	if err := c.doRequest(req, &resp); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	var onus []string
	twoWeeksAgo := time.Now().Add(-14 * 24 * time.Hour)

	for _, onu := range resp.Response {
		if onu.Status != "Online" {
			continue
		}

		lastReboot, err := time.Parse("2006-01-02 15:04:05", onu.LastStatusChange)
		if err != nil {
			fmt.Println("%w", err)
			continue
		}

		if lastReboot.Before(twoWeeksAgo) {
			onus = append(onus, onu.ID)
		}
	}

	return onus, nil
}

func (c *client) AuthorizeConnection(connectionDetails ConnectionDetails) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if err := writer.WriteField("olt_id", connectionDetails.OltID); err != nil {
		return fmt.Errorf("%w", err)
	}
	if err := writer.WriteField("pon_type", connectionDetails.PonType); err != nil {
		return fmt.Errorf("%w", err)
	}
	if err := writer.WriteField("gpon_channel", connectionDetails.GponChannel); err != nil {
		return fmt.Errorf("%w", err)
	}
	if err := writer.WriteField("board", connectionDetails.Board); err != nil {
		return fmt.Errorf("%w", err)
	}
	if err := writer.WriteField("port", connectionDetails.Port); err != nil {
		return fmt.Errorf("%w", err)
	}
	if err := writer.WriteField("sn", connectionDetails.SN); err != nil {
		return fmt.Errorf("%w", err)
	}
	if err := writer.WriteField("vlan", connectionDetails.VLAN); err != nil {
		return fmt.Errorf("%w", err)
	}
	if err := writer.WriteField("onu_type", connectionDetails.OnuType); err != nil {
		return fmt.Errorf("%w", err)
	}
	if err := writer.WriteField("zone", connectionDetails.Zone); err != nil {
		return fmt.Errorf("%w", err)
	}
	if err := writer.WriteField("odb", connectionDetails.ODB); err != nil {
		return fmt.Errorf("%w", err)
	}
	if err := writer.WriteField("name", connectionDetails.Name); err != nil {
		return fmt.Errorf("%w", err)
	}
	if err := writer.WriteField("address_or_comment", connectionDetails.AddressOrComment); err != nil {
		return fmt.Errorf("%w", err)
	}
	if err := writer.WriteField("onu_mode", connectionDetails.OnuMode); err != nil {
		return fmt.Errorf("%w", err)
	}
	if err := writer.WriteField("onu_external_id", connectionDetails.OnuExternalID); err != nil {
		return fmt.Errorf("%w", err)
	}
	writer.Close()

	url := fmt.Sprintf("%s%s", c.baseURL, types.AUTHORIZECONNECTION)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return c.doRequest(req, nil)
}

func (c *client) UnconfiguredOnusForOlt(oltID string) ([]types.UnconfiguredOnu, error) {
	url := fmt.Sprintf("%s%s%s", c.baseURL, types.UNCONFIGUREDONU, oltID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	var resp types.Response[types.UnconfiguredOnu]
	if err := c.doRequest(req, &resp); err != nil {
		return nil, err
	}

	return resp.Response, nil
}
