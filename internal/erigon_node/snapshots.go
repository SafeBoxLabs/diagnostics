package erigon_node

import (
	"context"
	"encoding/json"
)

type DownloadStatistics map[string]interface{}

type ShanshotFilesList interface{}

func (c *NodeClient) ShanphotSync(ctx context.Context) (DownloadStatistics, error) {
	var syncStats DownloadStatistics

	request, err := c.fetch(ctx, "snapshot-sync", nil)

	if err != nil {
		return syncStats, err
	}

	_, result, err := request.nextResult(ctx)

	if err != nil {
		return DownloadStatistics{}, err
	}

	if err := json.Unmarshal(result, &syncStats); err != nil {
		return DownloadStatistics{}, err
	}

	return syncStats, nil
}

func (c *NodeClient) ShanphotFiles(ctx context.Context) (ShanshotFilesList, error) {
	var filesList ShanshotFilesList

	request, err := c.fetch(ctx, "snapshot-files-list", nil)

	if err != nil {
		return filesList, err
	}

	_, result, err := request.nextResult(ctx)

	if err != nil {
		return filesList, err
	}

	if err := json.Unmarshal(result, &filesList); err != nil {
		return filesList, err
	}

	return filesList, nil
}
