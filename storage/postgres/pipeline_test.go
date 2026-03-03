package postgres_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/transcodeuz/transcode-rest/genproto/transcoder_service"
)

func TestFindPipelines(t *testing.T) {
	projects, err := strg.Postgres().PipelinesFind(context.Background(), &transcoder_service.GetListPipelineRequest{
		Page:     1,
		Limit:    10,
		FromDate: "2025-08-11",
		ToDate:   "2025-08-12",
	})

	asdfaf, _ := json.Marshal(projects)

	fmt.Println("projects", string(asdfaf))

	assert.NoError(t, err)
	assert.NotNil(t, projects)
}
