package interface_adapter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testReq struct {
	Title string `validate:"required"`
	Body  string `validate:"required"`
}

func TestValidRequest_Success(t *testing.T) {
	req := testReq{Title: "hello", Body: "world"}
	err := ValidRequest(req)
	assert.NoError(t, err)
}

func TestValidRequest_MissingRequired(t *testing.T) {
	req := testReq{Title: "", Body: ""}
	err := ValidRequest(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to validate request")
}

func TestValidRequest_PartialMissing(t *testing.T) {
	req := testReq{Title: "hello", Body: ""}
	err := ValidRequest(req)
	assert.Error(t, err)
}
