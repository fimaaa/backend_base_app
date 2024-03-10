package helper

import (
	"backend_base_app/shared/helper/str"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

const (
	textError             string = `error`
	textOk                string = `ok`
	codeSuccess           int    = 200
	codeBadRequestError   int    = 400
	codeUnauthorizedError int    = 401
	codeDatabaseError     int    = 402
	codeValidationError   int    = 403
	codeForbiddenError    int    = 403
	codeNotFound          int    = 404
)

// ResponseHelper ...
type ResponseHelper struct {
	C        *gin.Context
	Status   string
	Message  string
	Data     interface{}
	Code     int // not the http code
	CodeType string
	TraceID  string
}

// HTTPHelper ...
type HTTPHelper struct {
	Validate   *validator.Validate
	Translator ut.Translator
}

func (u *HTTPHelper) getTypeData(i interface{}) string {
	v := reflect.ValueOf(i)
	v = reflect.Indirect(v)

	return v.Type().String()
}

// GetStatusCode ...
func (u *HTTPHelper) GetStatusCode(err error) int {
	statusCode := http.StatusOK
	if err != nil {
		switch u.getTypeData(err) {
		case "models.ErrorUnauthorized":
			statusCode = http.StatusUnauthorized
		case "models.ErrorNotFound":
			statusCode = http.StatusNotFound
		case "models.ErrorConflict":
			statusCode = http.StatusConflict
		case "models.ErrorInternalServer":
			statusCode = http.StatusInternalServerError
		default:
			statusCode = http.StatusInternalServerError
		}
	}

	return statusCode
}

// SetResponse ...
// Set response data.
func (u *HTTPHelper) SetResponse(c *gin.Context, status string, message string, data interface{}, code int, codeType string, traceId string) ResponseHelper {
	return ResponseHelper{c, status, message, data, code, codeType, traceId}
}

// SendError ...
// Send error response to consumers.
func (u *HTTPHelper) SendError(c *gin.Context, message string, data interface{}, code int, codeType string, traceId string) error {
	res := u.SetResponse(c, `error`, message, data, code, codeType, traceId)

	return u.SendResponse(res)
}

// SendBadRequest ...
// Send bad request response to consumers.
func (u *HTTPHelper) SendBadRequest(c *gin.Context, message string, data interface{}, traceId string) error {
	res := u.SetResponse(c, `error`, message, data, codeBadRequestError, `badRequest`, traceId)

	return u.SendResponse(res)
}

// SendValidationError ...
// Send validation error response to consumers.
func (u *HTTPHelper) SendValidationError(c gin.Context, validationErrors validator.ValidationErrors) error {
	errorResponse := map[string][]string{}
	errorTranslation := validationErrors.Translate(u.Translator)
	for _, err := range validationErrors {
		errKey := str.Underscore(err.StructField())
		errorResponse[errKey] = append(errorResponse[errKey], errorTranslation[err.Namespace()])
	}

	c.JSON(200, map[string]interface{}{
		"code":         codeValidationError,
		"code_type":    "[Gateway] validationError",
		"code_message": errorResponse,
		"data":         u.EmptyJsonMap(),
	})

	return nil
}

// SendDatabaseError ...
// Send database error response to consumers.
func (u *HTTPHelper) SendDatabaseError(c *gin.Context, message string, data interface{}, traceId string) error {
	return u.SendError(c, message, data, codeDatabaseError, `databaseError`, traceId)
}

// SendUnauthorizedError ...
// Send unauthorized response to consumers.
func (u *HTTPHelper) SendUnauthorizedError(c *gin.Context, message string, data interface{}, traceId string) error {
	return u.SendError(c, message, data, codeUnauthorizedError, `unAuthorized`, traceId)
}

// SendForbiddenError ...
// Send forbidden response to consumers.
func (u *HTTPHelper) SendForbiddenError(c *gin.Context, message string, data interface{}, traceId string) error {
	return u.SendError(c, message, data, codeForbiddenError, `forbiddenAccess`, traceId)
}

// SendNotFoundError ...
// Send not found response to consumers.
func (u *HTTPHelper) SendNotFoundError(c *gin.Context, message string, data interface{}, traceId string) error {
	return u.SendError(c, message, data, codeNotFound, `notFound`, traceId)
}

// SendSuccess ...
// Send success response to consumers.
func (u *HTTPHelper) SendSuccess(c *gin.Context, message string, data interface{}, traceId string) error {
	res := u.SetResponse(c, `ok`, message, data, codeSuccess, `success`, traceId)

	return u.SendResponse(res)
}

// SendResponse ...
// Send response
func (u *HTTPHelper) SendResponse(res ResponseHelper) error {
	if len(res.Message) == 0 {
		res.Message = `success`
	}

	var resCode int
	if res.Code != 200 {
		resCode = http.StatusBadRequest
	} else {
		resCode = http.StatusOK
	}

	res.C.JSON(resCode, map[string]interface{}{
		"code":         res.Code,
		"code_type":    res.CodeType,
		"code_message": res.Message,
		"data":         res.Data,
		"trace_id":     res.TraceID,
	})
	return nil
}

func (u *HTTPHelper) EmptyJsonMap() map[string]interface{} {
	return make(map[string]interface{})
}
