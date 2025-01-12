package fusionbrain_api

const (
	fusionbrainApiHost      = "api-key.fusionbrain.ai"
	fusionbrainSecretKeyEnv = "FUSIONBRAIN_SECRET_KEY"
	fusionbrainApiKeyEnv    = "FUSIONBRAIN_API_KEY"
)

// Пути
const (
	fusionbrainAvailabilityPath = "/key/api/v1/text2image/availability"
	fusionbrainGeneratePath     = "/key/api/v1/text2image/run"
	fusionbrainCheckStatusPath  = "/key/api/v1/text2image/status/"
	fusionbrainModelsPath       = "/key/api/v1/models"
)

// Пути к стилям.
const (
	fusionbrainStylesWebUrl = "https://cdn.fusionbrain.ai/static/styles/web"
	fusionbrainStylesApiUrl = "https://cdn.fusionbrain.ai/static/styles/api"
)

// статусы генерации
const (
	FusionbrainGenerateStatusINITIAL    = "INITIAL"    //запрос получен, находится в очереди на обработку
	FusionbrainGenerateStatusPROCESSING = "PROCESSING" // запрос находится в процессе обработки
	FusionbrainGenerateDONE             = "DONE"       // задание выполнено
	FusionbrainGenerateFAIL             = "FAIL"       // задание не удалось выполнить.
)

// Доступность сервиса
const (
	FusionbrainModelStatusDisableByQueue = "DISABLED_BY_QUEUE"
)
