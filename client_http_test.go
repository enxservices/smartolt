package smartolt

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func fakeSmartOLT(t *testing.T, routes map[string]struct {
	code int
	body string
}) Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Token") != "token" {
			t.Errorf("missing X-Token on %s", r.URL.Path)
		}
		route, ok := routes[r.URL.Path]
		if !ok {
			t.Errorf("unexpected call %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(route.code)
		_, _ = w.Write([]byte(route.body))
	}))
	t.Cleanup(srv.Close)

	return NewSmartOLTClient("token", srv.URL)
}

func TestGetOnuSpeedProfiles(t *testing.T) {
	cli := fakeSmartOLT(t, map[string]struct {
		code int
		body string
	}{
		"/onu/get_onu_speed_profiles/HWTCB5D242B4": {200, `{"upload_speed_profile_name":"200M","download_speed_profile_name":"400M","status":true,"response_code":"success"}`},
		"/onu/get_onu_speed_profiles/UNKNOWN":      {400, `{"status":false,"error":"Invalid parameters: Please specify a valid ONU external ID","error_code":"specify_valid_field"}`},
	})

	profiles, err := cli.GetOnuSpeedProfiles("HWTCB5D242B4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if profiles.Download != "400M" || profiles.Upload != "200M" {
		t.Fatalf("got download %q upload %q, want 400M and 200M", profiles.Download, profiles.Upload)
	}

	_, err = cli.GetOnuSpeedProfiles("UNKNOWN")
	if err == nil || err.Error() != "Invalid parameters: Please specify a valid ONU external ID" {
		t.Fatalf("got error %v, want SmartOLT's message", err)
	}
}

func TestGetOnuDetailsReadsTheNestedDetails(t *testing.T) {
	cli := fakeSmartOLT(t, map[string]struct {
		code int
		body string
	}{
		"/onu/get_onu_details/HWTCB5D242B4": {200, `{"status":true,"response_code":"success","onu_details":{"unique_external_id":"HWTCB5D242B4","sn":"HWTCB5D242B4","zone_name":"Centro","odb_name":"CTO-12","service_ports":[{"service_port":"1018","vlan":"100","upload_speed":"200M","download_speed":"400M"}]}}`},
	})

	details, err := cli.GetOnuDetails("HWTCB5D242B4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if details.SN != "HWTCB5D242B4" || details.ZoneName != "Centro" || details.ODBName != "CTO-12" {
		t.Fatalf("details not read from onu_details: %+v", details)
	}
	if len(details.ServicePort) != 1 || details.ServicePort[0].DownloadSpeed != "400M" || details.ServicePort[0].UploadSpeed != "200M" {
		t.Fatalf("service ports not read: %+v", details.ServicePort)
	}
}

func TestErrorWithoutMessageReportsTheStatusCode(t *testing.T) {
	cli := fakeSmartOLT(t, map[string]struct {
		code int
		body string
	}{
		"/onu/get_onu_speed_profiles/X": {502, `<html>bad gateway</html>`},
	})

	_, err := cli.GetOnuSpeedProfiles("X")
	if err == nil || err.Error() != "502" {
		t.Fatalf("got error %v, want 502", err)
	}
}

func TestGetOdbsReadsThePortCount(t *testing.T) {
	cli := fakeSmartOLT(t, map[string]struct {
		code int
		body string
	}{
		"/system/get_odbs": {200, `{"status":true,"response":[{"id":"1","name":"BH001-01","latitude":null,"longitude":null,"zone_id":"1","zone_name":"1","nr_of_ports":"8","has_coords":"0"},{"id":"9","name":"BH002-03","latitude":null,"longitude":null,"zone_id":"2","zone_name":"2","nr_of_ports":"16","has_coords":"0"}]}`},
	})

	odbs, err := cli.GetOdbs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(odbs) != 2 || odbs[0].Ports != 8 || odbs[1].Ports != 16 || odbs[0].ZoneID != "1" {
		t.Fatalf("got %+v, want ports 8 and 16", odbs)
	}
}
