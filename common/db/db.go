package db

// GraphEndpoint graph.endpoint
type GraphEndpoint struct {
	ID       int64
	Endpoint string
}

// GraphTagEndpoint graph.tag_endpoint
type GraphTagEndpoint struct {
	ID         int64
	Tag        string
	EndpointID int64
}

// GraphEndpointCounter graph.endpoint_counter
type GraphEndpointCounter struct {
	ID         int64
	EndpointID int64
	Counter    string
}
