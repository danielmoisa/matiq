package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/danielmoisa/matiq/internal/actionruntime/common"

	es "github.com/elastic/go-elasticsearch/v8"
)

type OperationRunner struct {
	client    *es.Client
	operation Action
}

func (o *OperationRunner) search() (common.RuntimeResult, error) {
	// Build the request body.
	var buf bytes.Buffer
	var searchQuery map[string]interface{}
	if err := json.Unmarshal([]byte(o.operation.Query), &searchQuery); err != nil {
		return common.RuntimeResult{Success: false}, err
	}
	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return common.RuntimeResult{Success: false}, err
	}

	// Perform the search request.
	res, err := o.client.Search(
		o.client.Search.WithContext(context.Background()),
		o.client.Search.WithIndex(o.operation.Index),
		o.client.Search.WithBody(&buf),
		o.client.Search.WithTrackTotalHits(true),
		o.client.Search.WithPretty(),
	)
	defer res.Body.Close()
	if err != nil {
		return common.RuntimeResult{Success: false}, err
	}

	// Format the response body.
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return common.RuntimeResult{Success: false}, err
	}

	return common.RuntimeResult{Success: true, Rows: []map[string]interface{}{result}}, nil
}

func (o *OperationRunner) insert() (common.RuntimeResult, error) {
	// Build the request body.
	var buf bytes.Buffer
	var insertBody map[string]interface{}
	if err := json.Unmarshal([]byte(o.operation.Body), &insertBody); err != nil {
		return common.RuntimeResult{Success: false}, err
	}
	if err := json.NewEncoder(&buf).Encode(insertBody); err != nil {
		return common.RuntimeResult{Success: false}, err
	}

	// Perform the insert a document request.
	res, err := o.client.Create(
		o.operation.Index,
		"",
		&buf,
		o.client.Create.WithContext(context.Background()),
		o.client.Create.WithPretty(),
	)
	defer res.Body.Close()
	if err != nil {
		return common.RuntimeResult{Success: false}, err
	}

	// Format the response body.
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return common.RuntimeResult{Success: false}, err
	}

	return common.RuntimeResult{Success: true, Rows: []map[string]interface{}{result}}, nil
}

func (o *OperationRunner) get() (common.RuntimeResult, error) {

	// Perform the get request.
	res, err := o.client.Get(
		o.operation.Index,
		o.operation.ID,
		o.client.Get.WithContext(context.Background()),
		o.client.Get.WithPretty(),
	)
	defer res.Body.Close()
	if err != nil {
		return common.RuntimeResult{Success: false}, err
	}

	// Format the response body.
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return common.RuntimeResult{Success: false}, err
	}

	return common.RuntimeResult{Success: true, Rows: []map[string]interface{}{result}}, nil
}

func (o *OperationRunner) update() (common.RuntimeResult, error) {
	// Build the request body.
	var buf bytes.Buffer
	var updateBody map[string]interface{}
	if err := json.Unmarshal([]byte(o.operation.Body), &updateBody); err != nil {
		return common.RuntimeResult{Success: false}, err
	}
	if err := json.NewEncoder(&buf).Encode(updateBody); err != nil {
		return common.RuntimeResult{Success: false}, err
	}

	// Perform the update request.
	res, err := o.client.Update(
		o.operation.Index,
		o.operation.ID,
		&buf,
		o.client.Update.WithContext(context.Background()),
		o.client.Update.WithPretty(),
	)
	defer res.Body.Close()
	if err != nil {
		return common.RuntimeResult{Success: false}, err
	}

	// Format the response body.
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return common.RuntimeResult{Success: false}, err
	}

	return common.RuntimeResult{Success: true, Rows: []map[string]interface{}{result}}, nil
}

func (o *OperationRunner) delete() (common.RuntimeResult, error) {

	// Perform the delete request.
	res, err := o.client.Delete(
		o.operation.Index,
		o.operation.ID,
		o.client.Delete.WithContext(context.Background()),
		o.client.Delete.WithPretty(),
	)
	defer res.Body.Close()
	if err != nil {
		return common.RuntimeResult{Success: false}, err
	}

	// Format the response body.
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return common.RuntimeResult{Success: false}, err
	}

	return common.RuntimeResult{Success: true, Rows: []map[string]interface{}{result}}, nil
}
