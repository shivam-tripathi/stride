package controllers

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"quizizz.com/domain"
	qerror "quizizz.com/errors"
	"quizizz.com/logger"
)

// ParsedRequest holds the request metadata extracted from gin context
type ParsedRequest struct {
	Method    string `json:"method"`
	Route     string `json:"route"`
	IP        string `json:"ip,omitempty"`
	Protocol  string `json:"protocol"`
	Referer   string `json:"referer,omitempty"`
	UserAgent string `json:"userAgent"`
}

// ParsedResponse holds the data extracted from the response
type ParsedResponse struct {
	Size       int     `json:"size"`
	StatusCode int     `json:"statusCode"`
	TimeTaken  float64 `json:"timeTaken"`
}

// ParsedError holds the data extracted from error
type ParsedError struct {
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// parseRequest helps parse the request
func parseRequest(c *gin.Context) ParsedRequest {
	return ParsedRequest{
		Method:    c.Request.Method,
		Route:     c.Request.RequestURI,
		IP:        c.ClientIP(),
		Protocol:  c.Request.Proto,
		Referer:   c.Request.Header.Get("Referer"),
		UserAgent: c.Request.Header.Get("User-Agent"),
	}
}

// parseResponse helps parse the response
func parseResponse(c *gin.Context) ParsedResponse {
	startTimeEpocNS := c.GetInt64("x-start-time")
	endTimeEpocNS := time.Now().UnixMicro()
	return ParsedResponse{
		Size:       c.Writer.Size(),
		StatusCode: c.Writer.Status(),
		TimeTaken:  (float64(endTimeEpocNS) - float64(startTimeEpocNS)) / 1000,
	}
}

// parseError helps parse the error
func parseError(err qerror.BaseError) ParsedError {
	return ParsedError{
		Message: err.GetMsg(),
		Details: err.GetErrorContext(),
	}
}

// HTTPRequestData holds the data for the request
type HTTPRequestData struct {
	Req ParsedRequest `json:"req"`
}

// HTTPResponseData holds the data for the response
type HTTPResponseData struct {
	Req ParsedRequest  `json:"req"`
	Res ParsedResponse `json:"res"`
}

// HTTPResponseWithErrorData holds the data for the response with error
type HTTPResponseWithErrorData struct {
	Req ParsedRequest  `json:"req"`
	Res ParsedResponse `json:"res"`
	Err ParsedError    `json:"err"`
}

// LogRequest helps log the request data tied to a gin HTTP request
func LogRequest(c *gin.Context) {
	logger.GetLogger().Log.Info("http", zap.String("type", "http-request"), zap.Any("data", HTTPRequestData{
		Req: parseRequest(c),
	}))
}

// LogResponse helps log the response tied to a gin HTTP request
func LogResponse(c *gin.Context) {
	logger.GetLogger().Log.Info("http", zap.String("type", "http-response"), zap.Any("data", HTTPResponseData{
		Req: parseRequest(c),
		Res: parseResponse(c),
	}))
}

// LogError helps log the given error and associated metadata for the gin HTTP request
func LogError(c *gin.Context, err qerror.BaseError) {
	logger.GetLogger().Log.Error("http", zap.String("type", "http-response"), zap.Any("data", HTTPResponseWithErrorData{
		Req: parseRequest(c),
		Res: parseResponse(c),
		Err: parseError(err),
	}))
}

// HTTPController is a wrapper over controllers to convert them into a HTTP controller
func HTTPController(controller Controller) func(c *gin.Context) {
	return func(c *gin.Context) {
		var request *domain.Request

		defer func() {
			if err := recover(); err != nil {
				HandleHTTPError(c, qerror.NewInternalServerError(
					qerror.WithError(err),
					qerror.WithStackTrace(string(debug.Stack())),
					qerror.WithErrorContext(request),
				))
			}
		}()

		LogRequest(c)

		request = controller.NewRequest()

		if request.Body != nil {
			if err := c.ShouldBindBodyWith(request.Body, binding.JSON); err != nil {
				if err.Error() == "EOF" {
					HandleHTTPError(c, qerror.NewBadRequestError(qerror.WithMsg("invalid empty body"), qerror.WithErrorContext("invalid empty body")))
					return
				}

				payload := map[string]interface{}{} // Only handles JSON which is a map at the moment
				rawErr := c.ShouldBindBodyWith(&payload, binding.JSON)

				if rawErr != nil {
					HandleHTTPError(c, qerror.NewBadRequestError(qerror.WithMsg(err.Error()), qerror.WithErrorContext(rawErr)))
				} else {
					HandleHTTPError(c, qerror.NewBadRequestError(qerror.WithMsg(err.Error()), qerror.WithErrorContext(payload)))
				}
				return
			}
		}

		for key := range request.Query {
			request.Query[key] = c.Query(key)
		}

		for key := range request.Params {
			request.Params[key] = c.Param(key)
		}

		if request.Headers != nil && len(request.Headers) > 0 {
			for key := range request.Headers {
				if _, ok := c.Request.Header[key]; ok {
					if header := c.Request.Header.Get(key); header != "" && header != "null" {
						request.Headers[key] = &header
					}
				}
			}
		}

		logger.GetLogger().Log.Debug("http-request-debug", zap.Any("request", request))

		if err := controller.SanitizeRequest(c, request); err != nil {
			HandleHTTPError(c, qerror.NewBadRequestError(
				qerror.WithMsg(err.Error()),
				qerror.WithErrorContext(request)),
			)
			return
		}

		if err := controller.ValidateRequest(c, request); err != nil {
			HandleHTTPError(c, qerror.NewBadRequestError(
				qerror.WithMsg(err.Error()),
				qerror.WithErrorContext(request),
			))
			return
		}

		response, err := controller.Handler(c, request)

		if err != nil {
			HandleHTTPError(c, err)
			return
		}

		c.JSON(http.StatusOK, domain.NewOkResponse(response))
		LogResponse(c)
	}
}

// HandleHTTPError handles errors in http requests
func HandleHTTPError(c *gin.Context, err interface{}) {
	requestError := qerror.HandleError(err)
	LogError(c, requestError)
	c.JSON(requestError.GetStatusCode(), domain.NewErrorResponse(map[string]interface{}{
		"context": requestError.GetUserContext(),
		"message": requestError.GetMsg(),
	}))
}

// NoRoute handles the case when the route is not found
func NoRoute(router *gin.Engine) {
	router.NoRoute(func(c *gin.Context) {
		LogRequest(c)
		logger.GetLogger().Log.Debug("http-request-debug", zap.Any("request", c.Request), zap.Any("params", c.Params))
		HandleHTTPError(c, qerror.NewResourceNotFoundError(qerror.WithMsg("route not found")))
	})
}
