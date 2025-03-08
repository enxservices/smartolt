package types

type StatusSignal struct {
	Status         bool   `json:"status"`
	OnuSignal      string `json:"onu_signal"`
	OnuSignalValue string `json:"onu_signal_value"`
	OnuSignal1310  string `json:"onu_signal_1310"`
	OnuSignal1490  string `json:"onu_signal_1490"`
}

type OnuDetails struct {
	ID                   int           `json:"unique_external_id"`
	SN                   string        `json:"sn"`
	Name                 string        `json:"name"`
	OltID                string        `json:"olt_id"`
	OltName              string        `json:"olt_name"`
	Board                string        `json:"board"`
	Port                 string        `json:"port"`
	Onu                  string        `json:"onu"`
	OnuTypeID            string        `json:"onu_type_id"`
	OnuTypeName          string        `json:"onu_type_name"`
	ZoneID               string        `json:"zone_id"`
	ZoneName             string        `json:"zone_name"`
	Address              *string       `json:"address"`
	ODBName              string        `json:"odb_name"`
	Mode                 string        `json:"mode"`
	WanMode              string        `json:"wan_mode"`
	IpAddress            *string       `json:"ip_address"`
	SubnetMask           *string       `json:"subnet_mask"`
	DefaultGateway       *string       `json:"default_gateway"`
	DNS1                 *string       `json:"dns1"`
	DNS2                 *string       `json:"dns2"`
	Username             *string       `json:"username"`
	Password             *string       `json:"password"`
	CatV                 *string       `json:"catv"`
	AdministrativeStatus string        `json:"administrative_status"`
	ServicePort          []ServicePort `json:"service_ports"`
}

type ServicePort struct {
	Port             string `json:"service_port"`
	Vlan             string `json:"vlan"`
	CVlan            string `json:"c_vlan"`
	SVlan            string `json:"s_vlan"`
	TagTransformMode string `json:"tag_transform_mode"`
	UploadSpeed      string `json:"upload_speed"`
	DownloadSpeed    string `json:"download_speed"`
}

type SpeedProfile struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Speed     string `json:"speed"`
	Direction string `json:"direction"`
	Type      string `json:"type"`
}

type ResponseError struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type Response[T any] struct {
	Status   bool `json:"status"`
	Response []T  `json:"response"`
}

type SmartOLT interface {
	GetOnuDetails(ID string) (*OnuDetails, error)
	GetOnuSignal(ID string) (*StatusSignal, error)
	GetSpeedProfiles() ([]SpeedProfile, error)
	UpdateSpeedProfile(ID, downloadProfile, uploadProfile string) error
	RebootOnu(ID string) error
	DisableOnu(ID string) error
}
