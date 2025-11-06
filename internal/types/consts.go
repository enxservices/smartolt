package types

const (
	GETSPEEDPROFILE     = "/system/get_speed_profiles"
	GETONUDETAILS       = "/onu/get_onu_details"
	GETONUSIGNAL        = "/onu/get_onu_signal"
	GETALLONUSDETAILS   = "/onu/get_all_onus_details"
	UPDATESPEEDPROFILE  = "/onu/update_onu_speed_profiles"
	REBOOTONU           = "/onu/reboot"
	DISABLEONU          = "/onu/disable/"
	ENABLEONU           = "/onu/enable/"
	GETODBS             = "/system/get_odbs"
	STATUSESONUS        = "/onu/get_onu_statuses?olt_id=1"
	AUTHORIZECONNECTION = "/onu/authorize_onu"
	UNCONFIGUREDONU     = "/onu/unconfigured_onus_for_olt/"
	DefaultAPIKey       = ""
	DefaultBaseURL      = "https://enx.smartolt.com/api"
)
