package smartolt

import (
	"strings"

	"github.com/enxservices/smartolt/internal/types"
)

func CalculateODBAvailability(obds []types.ODB, onus []types.OnuListItem) []types.ODBAvailability {
	onuCountByODB := make(map[string]int)
	for _, onu := range onus {
		name := strings.TrimSpace(strings.ToLower(onu.ODBName))
		if name == "" {
			continue
		}
		onuCountByODB[name]++
	}

	result := make([]types.ODBAvailability, 0, len(obds))
	for _, odb := range obds {
		key := strings.TrimSpace(strings.ToLower(odb.Name))
		used := onuCountByODB[key]
		available := odb.Ports - used
		result = append(result, types.ODBAvailability{
			OdbID:          odb.ID,
			OdbName:        odb.Name,
			TotalPorts:     odb.Ports,
			UsedPorts:      used,
			AvailablePorts: available,
		})
	}
	return result
}
