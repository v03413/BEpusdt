package api

import (
	"encoding/json"
	"fmt"
)

type BoundaryResp struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

func (a Api) GetBoundaryHeights(start, end int64, network string) (int64, int64) {
	url := fmt.Sprintf("%s/api/block/get_boundary_heights?network=%s&start=%d&end=%d",
		a.api, network, start, end)

	body, err := a.get(url)
	if err != nil {
		a.error(err)
		return 0, 0
	}

	var response BoundaryResp
	if err := json.Unmarshal(body, &response); err != nil {
		a.error(err)
		return 0, 0
	}

	return response.Start, response.End
}
