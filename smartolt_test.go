package smartolt

import (
	"testing"

	"github.com/enxservices/smartolt/internal/types"
)

func TestSmartOLT(t *testing.T) {
	smartOLTClient := NewSmartOLTClient(types.DefaultAPIKey, types.DefaultBaseURL)
	ONUExternalID := "HWTCAD08F2AA"

	t.Run("Should reboot ONU successfully", func(t *testing.T) {
		err := smartOLTClient.RebootOnu(ONUExternalID)
		if err != nil {
			t.Fatalf("Expected success, but got error: %v", err)
		}
	})

	t.Run("Should list speed profiles successfully", func(t *testing.T) {
		profiles, err := smartOLTClient.GetSpeedProfiles()
		if err != nil {
			t.Fatalf("Expected success, but got error: %v", err)
		}

		if len(profiles) == 0 {
			t.Fatalf("Expected at least one speed profile, but got none")
		}

		t.Logf("profiles: [%v]", profiles)
	})

	t.Run("Should update speed profile successfully", func(t *testing.T) {
		err := smartOLTClient.UpdateSpeedProfile(ONUExternalID, "1G", "1G")
		if err != nil {
			t.Fatalf("Expected success, but got error: [%v]", err)
		}
	})

	t.Run("Should not update speed profile because Download speed profile doesnt exists", func(t *testing.T) {
		downSpeed := "11G"
		err := smartOLTClient.UpdateSpeedProfile(ONUExternalID, downSpeed, "1G")
		if err == nil {
			t.Fatal("Expected an error, but got none")
		}

		expectedErrMsg := "Invalid parameters: No such Download speed profile exists"
		if err.Error() != expectedErrMsg {
			t.Fatalf("Expected error message [%s], but got [%s]", expectedErrMsg, err.Error())
		}
	})

	t.Run("Should not update speed profile because Upload speed profile doenst exists", func(t *testing.T) {
		upSpeed := "11GB"
		err := smartOLTClient.UpdateSpeedProfile(ONUExternalID, "1G", upSpeed)
		if err == nil {
			t.Fatal("Expected an error, but got none")
		}

		expectErrMsg := "Invalid parameters: No such Upload speed profile exists"
		if err.Error() != expectErrMsg {
			t.Fatalf("Expected error message [%s], but got [%s]", expectErrMsg, err.Error())
		}

	})

	t.Run("Should return onu signal sucessfully", func(t *testing.T) {
		resp, err := smartOLTClient.GetOnuSignal(ONUExternalID)
		if err != nil {
			t.Fatalf("Expected sucess, but got error: [%v]", err)
		}

		t.Logf("ONU Signal: [%v]", resp)
	})

	t.Run("Should return onu details sucessfully", func(t *testing.T) {
		resp, err := smartOLTClient.GetOnuDetails(ONUExternalID)
		if err != nil {
			t.Fatalf("Expected sucess, but got error: [%v]", err)
		}

		t.Logf("ONU Details: [%v]", resp)
	})
}

func TestCalculateODBAvailability(t *testing.T) {
	obds := []types.ODB{
		{ID: "2", Name: "test odb", ZoneID: "1", Ports: 8},
		{ID: "3", Name: "test odb2", ZoneID: "18", Ports: 16},
	}
	onus := []types.OnuListItem{
		{UniqueExternalID: "test1", ODBName: "test odb"},
		{UniqueExternalID: "test2", ODBName: "test odb2"},
		{UniqueExternalID: "test3", ODBName: "test odb2"},
	}

	got := CalculateODBAvailability(obds, onus)

	if len(got) != 2 {
		t.Fatalf("expected 2 ODB availability entries, got %d", len(got))
	}

	// Map by ODB ID for assertions
	byID := map[string]types.ODBAvailability{}
	for _, g := range got {
		byID[g.OdbID] = g
	}

	a := byID["2"]
	if a.TotalPorts != 8 || a.UsedPorts != 1 || a.AvailablePorts != 7 {
		t.Fatalf("unexpected availability for ODB 2: %+v", a)
	}

	b := byID["3"]
	if b.TotalPorts != 16 || b.UsedPorts != 2 || b.AvailablePorts != 14 {
		t.Fatalf("unexpected availability for ODB 3: %+v", b)
	}
}
