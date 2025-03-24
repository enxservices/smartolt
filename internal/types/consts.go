package types

const (
	GETSPEEDPROFILE    = "/system/get_speed_profiles"
	GETONUDETAILS      = "/onu/get_onu_details"
	GETONUSIGNAL       = "/onu/get_onu_signal"
	UPDATESPEEDPROFILE = "/onu/update_onu_speed_profiles"
	REBOOTONU          = "/onu/reboot"
	DISABLEONU         = "/onu/disable/"
	ENABLEONU          = "/onu/enable"
	STATUSESONUS       = "/onu/get_onu_statuses?olt_id=1" //temos apenas uma OLT...

	DefaultAPIKey  = ""
	DefaultBaseURL = "https://enx.smartolt.com/api"
)
