package transformers

type Transformer struct {
	Data interface{} `json:"data"`
}

type CollectionTransformer struct {
	Data []interface{} `json:"data"`
}
