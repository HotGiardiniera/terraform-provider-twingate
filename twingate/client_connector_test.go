package twingate

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestClientConnectorCreateOk(t *testing.T) {
	t.Run("Test Twingate Resource : Client Connector Create Ok", func(t *testing.T) {
		// response JSON
		createConnectorOkJson := `{
	  "data": {
		"connectorCreate": {
		  "entity": {
			"id": "test-id",
			"name" : "test-name"
		  },
		  "ok": true,
		  "error": null
		}
	  }
	}`

		client := newHTTPMockClient()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("POST", client.GraphqlServerURL,
			httpmock.NewStringResponder(200, createConnectorOkJson))

		connector, err := client.createConnector(context.Background(), "test")

		assert.Nil(t, err)
		assert.EqualValues(t, "test-id", connector.ID)
		assert.EqualValues(t, "test-name", connector.Name)
	})
}

func TestClientConnectorUpdateOk(t *testing.T) {
	t.Run("Test Twingate Resource : Client Connector Update Ok", func(t *testing.T) {
		// response JSON
		updateConnectorOkJson := `{
	  "data": {
		"connectorUpdate": {
		  "entity": {
			"id": "test-id",
			"name" : "test-name"
		  },
		  "ok": true,
		  "error": null
		}
	  }
	}`

		client := newHTTPMockClient()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("POST", client.GraphqlServerURL,
			httpmock.NewStringResponder(200, updateConnectorOkJson))

		err := client.updateConnector(context.Background(), "test-id", "test-name")

		assert.Nil(t, err)
	})
}

func TestClientConnectorDeleteOk(t *testing.T) {
	t.Run("Test Twingate Resource : Client Connector Delete Ok", func(t *testing.T) {
		// response JSON
		deleteConnectorOkJson := `{
		  "data": {
			"connectorDelete": {
			  "ok": true,
			  "error": null
			}
		  }
		}`

		client := newHTTPMockClient()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("POST", client.GraphqlServerURL,
			httpmock.NewStringResponder(200, deleteConnectorOkJson))

		err := client.deleteConnector(context.Background(), "test")

		assert.NoError(t, err)
	})
}

func TestClientConnectorCreateError(t *testing.T) {
	t.Run("Test Twingate Resource : Client Connector Create Error", func(t *testing.T) {

		// response JSON
		createNetworkErrorJson := `{
	  "data": {
		"connectorCreate": {
		  "ok": false,
		  "error": "error_1"
		}
	  }
	}`

		client := newHTTPMockClient()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("POST", client.GraphqlServerURL,
			httpmock.NewStringResponder(200, createNetworkErrorJson))

		remoteNetwork, err := client.createConnector(context.Background(), "test")

		assert.EqualError(t, err, "failed to create connector: error_1")
		assert.Nil(t, remoteNetwork)
	})
}

func TestClientConnectorUpdateError(t *testing.T) {
	t.Run("Test Twingate Resource : Client Connector Update Error", func(t *testing.T) {

		// response JSON
		createNetworkOkJson := `{
		  "data": {
			"connectorUpdate": {
			  "ok": false,
			  "error": "error_1"
			}
		  }
		}`

		client := newHTTPMockClient()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("POST", client.GraphqlServerURL,
			httpmock.NewStringResponder(200, createNetworkOkJson))
		connectorId := "test-id"

		err := client.updateConnector(context.Background(), connectorId, "test-name")

		assert.EqualError(t, err, fmt.Sprintf("failed to update connector with id %s: error_1", connectorId))
	})
}

func TestClientConnectorUpdateErrorWhenIdEmpty(t *testing.T) {
	t.Run("Test Twingate Resource : Client Connector Update Error on empty ID", func(t *testing.T) {

		// response JSON
		createNetworkOkJson := `{
		  "data": {
			"connectorUpdate": {
			  "ok": false,
			  "error": "error_1"
			}
		  }
		}`

		client := newHTTPMockClient()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("POST", client.GraphqlServerURL,
			httpmock.NewStringResponder(200, createNetworkOkJson))

		err := client.updateConnector(context.Background(), "", "")

		assert.EqualError(t, err, "failed to update connector: connector id is empty")
	})
}

func TestClientConnectorEmptyNetworkIDCreateError(t *testing.T) {
	t.Run("Test Twingate Resource : Client Connector Empty Network ID Create Error", func(t *testing.T) {

		// response JSON
		createNetworkOkJson := `{}`

		client := newHTTPMockClient()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("POST", client.GraphqlServerURL,
			httpmock.NewStringResponder(200, createNetworkOkJson))

		remoteNetwork, err := client.createConnector(context.Background(), "")

		assert.EqualError(t, err, "failed to create connector: network id is empty")
		assert.Nil(t, remoteNetwork)
	})
}

func TestClientConnectorDeleteError(t *testing.T) {
	t.Run("Test Twingate Resource : Client Connector Delete Error", func(t *testing.T) {

		// response JSON
		deleteConnectorOkJson := `{
		  "data": {
			"connectorDelete": {
			  "ok": false,
			  "error": "error_1"
			}
		  }
		}`

		client := newHTTPMockClient()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("POST", client.GraphqlServerURL,
			httpmock.NewStringResponder(200, deleteConnectorOkJson))
		connectorId := "test"

		err := client.deleteConnector(context.Background(), connectorId)

		assert.EqualError(t, err, fmt.Sprintf("failed to delete connector with id %s: error_1", connectorId))
	})
}

func TestClientConnectorReadError(t *testing.T) {
	t.Run("Test Twingate Resource : Client Connector Read Error", func(t *testing.T) {

		// response JSON
		readNetworkOkJson := `{
	  "data": {
		"connector": null
	  }
	}`

		client := newHTTPMockClient()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("POST", client.GraphqlServerURL,
			httpmock.NewStringResponder(200, readNetworkOkJson))
		connectorId := "test"

		connector, err := client.readConnector(context.Background(), connectorId)

		assert.Nil(t, connector)
		assert.EqualError(t, err, fmt.Sprintf("failed to read connector with id %s: query result is empty", connectorId))
	})
}

func TestClientConnectorReadEmptyError(t *testing.T) {
	t.Run("Test Twingate Resource : Client Connector Read Error", func(t *testing.T) {

		// response JSON
		readConnectorFAIL := `{}`

		client := newHTTPMockClient()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("POST", client.GraphqlServerURL,
			httpmock.NewStringResponder(200, readConnectorFAIL))

		connectors, _ := client.readConnectors(context.Background())

		assert.Empty(t, connectors)
	})
}

func TestClientConnectorEmptyReadError(t *testing.T) {
	t.Run("Test Twingate Resource : Client Connector Empty Read Error", func(t *testing.T) {

		// response JSON
		readConnectorOkJson := `{}`

		client := newHTTPMockClient()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("POST", client.GraphqlServerURL,
			httpmock.NewStringResponder(200, readConnectorOkJson))

		connector, err := client.readConnector(context.Background(), "")

		assert.Nil(t, connector)
		assert.EqualError(t, err, "failed to read connector: id is empty")
	})
}

func TestClientConnectorEmptyDeleteError(t *testing.T) {
	t.Run("Test Twingate Resource : Client Connector Empty Delete Error", func(t *testing.T) {

		// response JSON
		deleteConnectorOkJson := `{}`

		client := newHTTPMockClient()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("POST", client.GraphqlServerURL,
			httpmock.NewStringResponder(200, deleteConnectorOkJson))

		err := client.deleteConnector(context.Background(), "")

		assert.EqualError(t, err, "failed to delete connector: id is empty")
	})
}

func TestClientConnectorReadAllOk(t *testing.T) {
	t.Run("Test Twingate Resource : Read All Client Connectors", func(t *testing.T) {

		// response JSON
		readConnectorsOkJson := `{
	  "data": {
		"connectors": {
		  "edges": [
			{
			  "node": {
				"id": "connector1",
				"name": "tf-acc-connector1"
			  }
			},
			{
			  "node": {
				"id": "connector2",
				"name": "connector2"
			  }
			},
			{
			  "node": {
				"id": "connector3",
				"name": "tf-acc-connector3"
			  }
			}
		  ]
		}
	  }
	}`

		client := newHTTPMockClient()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("POST", client.GraphqlServerURL,
			httpmock.NewStringResponder(200, readConnectorsOkJson))

		connector, err := client.readConnectors(context.Background())
		assert.NoError(t, err)

		r0 := &Connectors{
			ID:   "connector1",
			Name: "tf-acc-connector1",
		}
		r1 := &Connectors{
			ID:   "connector2",
			Name: "connector2",
		}
		r2 := &Connectors{
			ID:   "connector3",
			Name: "tf-acc-connector3",
		}
		mockMap := make(map[int]*Connectors)

		mockMap[0] = r0
		mockMap[1] = r1
		mockMap[2] = r2

		counter := 0
		for _, elem := range connector {
			for _, i := range mockMap {
				if elem.Name == i.Name && elem.ID == i.ID {
					counter++
				}
			}
		}

		if len(mockMap) != counter {
			t.Errorf("Expected map not equal to origin!")
		}
		assert.EqualValues(t, len(mockMap), counter)
	})
}

func TestClientConnectorUpdateRequestError(t *testing.T) {
	t.Run("Test Twingate Resource : Client Connector Update Request Error", func(t *testing.T) {

		client := newHTTPMockClient()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("POST", client.GraphqlServerURL,
			httpmock.NewErrorResponder(errors.New("error_1")))
		connectorId := "test"

		err := client.updateConnector(context.Background(), connectorId, "new name")

		assert.EqualError(t, err, fmt.Sprintf(`failed to update connector with id %s: Post "%s": error_1`, connectorId, client.GraphqlServerURL))
	})
}

func TestClientConnectorDeleteRequestError(t *testing.T) {
	t.Run("Test Twingate Resource : Client Connector Delete Request Error", func(t *testing.T) {

		client := newHTTPMockClient()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("POST", client.GraphqlServerURL,
			httpmock.NewErrorResponder(errors.New("error_1")))
		connectorId := "test"

		err := client.deleteConnector(context.Background(), connectorId)

		assert.EqualError(t, err, fmt.Sprintf(`failed to delete connector with id %s: Post "%s": error_1`, connectorId, client.GraphqlServerURL))
	})
}

func TestClientConnectorCreateRequestError(t *testing.T) {
	t.Run("Test Twingate Resource : Client Connector Create Request Error", func(t *testing.T) {

		client := newHTTPMockClient()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("POST", client.GraphqlServerURL,
			httpmock.NewErrorResponder(errors.New("error_1")))

		remoteNetwork, err := client.createConnector(context.Background(), "test")

		assert.EqualError(t, err, fmt.Sprintf(`failed to create connector: Post "%s": error_1`, client.GraphqlServerURL))
		assert.Nil(t, remoteNetwork)
	})
}

func TestClientConnectorReadRequestError(t *testing.T) {
	t.Run("Test Twingate Resource : Client Connector Read Request Error", func(t *testing.T) {

		client := newHTTPMockClient()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("POST", client.GraphqlServerURL,
			httpmock.NewErrorResponder(errors.New("error_1")))
		connectorId := "test"

		connector, err := client.readConnector(context.Background(), connectorId)

		assert.Nil(t, connector)
		assert.EqualError(t, err, fmt.Sprintf(`failed to read connector with id %s: Post "%s": error_1`, connectorId, client.GraphqlServerURL))
	})
}
