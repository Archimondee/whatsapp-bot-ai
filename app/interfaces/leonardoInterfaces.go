package interfaces

type Request struct {
	Input any `json:"input"`
}

type Model struct {
	Input map[string]interface{} `json:"input"`
}

type Response struct {
	SdGenerationJob struct {
		GenerationId string `json:"generationId"`
	} `json:"sdGenerationJob"`
}

type ResponseGenerate struct {
	GenerationsByPK GenerationsByPK `json:"generations_by_pk"`
}

type ResponseInitImage struct {
	UploadInitImage UploadData `json:"uploadInitImage"`
}

func (m *Model) CreateRequest() *Request {
	return &Request{
		Input: m.Input,
	}
}

func NewModel(prompt string, modelId string) (model *Model) {
	model = &Model{
		Input: map[string]interface{}{
			"prompt":     prompt,
			"modelId":    modelId,
			"width":      1024,
			"height":     1024,
			"sd_version": "v2",
			"num_images": 1,
		},
	}
	return
}

type GeneratedImage struct {
	URL                             string     `json:"url"`
	NSFW                            bool       `json:"nsfw"`
	ID                              string     `json:"id"`
	LikeCount                       int        `json:"likeCount"`
	GeneratedImageVariationGenerics []struct{} `json:"generated_image_variation_generics"`
}

type GenerationsByPK struct {
	GeneratedImages []GeneratedImage `json:"generated_images"`
	ModelID         string           `json:"modelId"`
	Prompt          string           `json:"prompt"`
	NegativePrompt  string           `json:"negativePrompt"`
	ImageHeight     int              `json:"imageHeight"`
	ImageWidth      int              `json:"imageWidth"`
	InferenceSteps  int              `json:"inferenceSteps"`
	Seed            int              `json:"seed"`
	Public          bool             `json:"public"`
	Scheduler       string           `json:"scheduler"`
	SDVersion       string           `json:"sdVersion"`
	Status          string           `json:"status"`
	PresetStyle     interface{}      `json:"presetStyle"`
	InitStrength    interface{}      `json:"initStrength"`
	GuidanceScale   int              `json:"guidanceScale"`
	ID              string           `json:"id"`
	CreatedAt       any              `json:"createdAt"`
}

type UploadData struct {
	Fields string `json:"fields"`
	URL    string `json:"url"`
	Key    string `json:"key"`
	Id     string `json:"id"`
}

type FieldData struct {
	ContentType       string `json:"Content-Type"`
	Bucket            string `json:"bucket"`
	XAmzAlgorithm     string `json:"X-Amz-Algorithm"`
	XAmzCredential    string `json:"X-Amz-Credential"`
	XAmzDate          string `json:"X-Amz-Date"`
	XAmzSecurityToken string `json:"X-Amz-Security-Token"`
	XAmzSignature     string `json:"X-Amz-Signature"`
	Key               string `json:"key"`
	Policy            string `json:"Policy"`
}
