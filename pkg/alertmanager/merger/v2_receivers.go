package merger

import (
	"github.com/go-openapi/swag/jsonutils"
	v2_models "github.com/prometheus/alertmanager/api/v2/models"
)

// V2Receivers implements the Merger interface for GET /v2/receivers. It returns the union of receivers
// over all the responses. When a receiver with the same name exists in multiple responses, any one
// of them is returned (they should be identical across replicas).
type V2Receivers struct{}

func (V2Receivers) MergeResponses(in [][]byte) ([]byte, error) {
	seen := make(map[string]bool)
	result := make([]*v2_models.Receiver, 0)

	for _, body := range in {
		parsed := make([]*v2_models.Receiver, 0)
		if err := jsonutils.ReadJSON(body, &parsed); err != nil {
			return nil, err
		}

		// Deduplicate receivers by name as we process each response
		for _, receiver := range parsed {
			if receiver.Name == nil {
				continue
			}

			name := *receiver.Name
			if !seen[name] {
				seen[name] = true
				result = append(result, receiver)
			}
		}
	}

	return jsonutils.WriteJSON(result)
}
