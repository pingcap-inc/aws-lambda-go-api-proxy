package httpadapter

import (
	"context"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
)

type HandlerAdapterALB struct {
	core.RequestAccessor
	handler http.Handler
}

func NewALB(handler http.Handler) *HandlerAdapterALB {
	return &HandlerAdapterALB{
		handler: handler,
	}
}

// Proxy receives an API Gateway proxy event, transforms it into an http.Request
// object, and sends it to the http.HandlerFunc for routing.
// It returns a proxy response object generated from the http.Handler.
func (h *HandlerAdapterALB) Proxy(event events.APIGatewayProxyRequest) (events.ALBTargetGroupResponse, error) {
	req, err := h.ProxyEventToHTTPRequest(event)
	return h.proxyInternal(req, err)
}

// ProxyWithContext receives context and an API Gateway proxy event,
// transforms them into an http.Request object, and sends it to the http.Handler for routing.
// It returns a proxy response object generated from the http.ResponseWriter.
func (h *HandlerAdapterALB) ProxyWithContext(ctx context.Context, event events.APIGatewayProxyRequest) (events.ALBTargetGroupResponse, error) {
	req, err := h.EventToRequestWithContext(ctx, event)
	return h.proxyInternal(req, err)
}

func (h *HandlerAdapterALB) proxyInternal(req *http.Request, err error) (events.ALBTargetGroupResponse, error) {
	if err != nil {
		return events.ALBTargetGroupResponse{StatusCode: http.StatusGatewayTimeout}, core.NewLoggedError("Could not convert proxy event to request: %v", err)
	}

	w := core.NewProxyResponseWriter()
	h.handler.ServeHTTP(http.ResponseWriter(w), req)

	resp, err := w.GetProxyResponse()
	if err != nil {
		return events.ALBTargetGroupResponse{StatusCode: http.StatusGatewayTimeout}, core.NewLoggedError("Error while generating proxy response: %v", err)
	}
	hd := make(map[string]string)
	for k, s := range resp.MultiValueHeaders {
		hd[k] = strings.Join(s, ";")
	}
	return events.ALBTargetGroupResponse{
		StatusCode:        resp.StatusCode,
		StatusDescription: http.StatusText(resp.StatusCode),
		Headers:           hd,
		MultiValueHeaders: resp.MultiValueHeaders,
		Body:              resp.Body,
		IsBase64Encoded:   resp.IsBase64Encoded,
	}, nil
}
