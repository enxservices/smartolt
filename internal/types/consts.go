package types

const (
	GETSPEEDPROFILE    = "smartolt.com/api/system/get_speed_profiles"
	GETONUDETAILS      = "smartolt.com/api/onu/get_onu_details"
	GETONUSIGNAL       = "smartolt.com/api/onu/get_onu_signal"
	UPDATESPEEDPROFILE = "smartolt.com/api/onu/update_onu_speed_profiles"
	REBOOTONU          = "smartolt.com/api/onu/reboot"
	DISABLEONU         = "smartolt.com/api/onu/disable/"
	ENABLEONU          = "smartolt.com/api/onu/enable"
	STATUSESONUS       = "smartolt.com/api/onu/get_onu_statuses?olt_id=1" //temos apenas uma OLT...

	DefaultAPIKey  = ""
	DefaultBaseURL = "https://enx."
)
