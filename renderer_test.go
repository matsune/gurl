package gurl

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/fatih/color"
)

func Test_colorForStatus(t *testing.T) {
	tests := []struct {
		name  string
		codes []int
		want  color.Attribute
	}{
		{
			name: "Green",
			codes: []int{
				http.StatusOK,
				http.StatusCreated,
				http.StatusAccepted,
				http.StatusNonAuthoritativeInfo,
				http.StatusNoContent,
				http.StatusResetContent,
				http.StatusPartialContent,
				http.StatusMultiStatus,
				http.StatusAlreadyReported,
				http.StatusIMUsed,
			},
			want: color.FgGreen,
		},
		{
			name: "Cyan",
			codes: []int{
				http.StatusMultipleChoices,
				http.StatusMovedPermanently,
				http.StatusFound,
				http.StatusSeeOther,
				http.StatusNotModified,
				http.StatusUseProxy,
				306,
				http.StatusTemporaryRedirect,
				http.StatusPermanentRedirect,
			},
			want: color.FgCyan,
		},
		{
			name: "yellow",
			codes: []int{
				http.StatusBadRequest,
				http.StatusUnauthorized,
				http.StatusPaymentRequired,
				http.StatusForbidden,
				http.StatusNotFound,
				http.StatusMethodNotAllowed,
				http.StatusNotAcceptable,
				http.StatusProxyAuthRequired,
				http.StatusRequestTimeout,
				http.StatusConflict,
				http.StatusGone,
				http.StatusLengthRequired,
				http.StatusPreconditionFailed,
				http.StatusRequestEntityTooLarge,
				http.StatusRequestURITooLong,
				http.StatusUnsupportedMediaType,
				http.StatusRequestedRangeNotSatisfiable,
				http.StatusExpectationFailed,
				http.StatusTeapot,
				http.StatusMisdirectedRequest,
				http.StatusUnprocessableEntity,
				http.StatusLocked,
				http.StatusFailedDependency,
				http.StatusUpgradeRequired,
				http.StatusPreconditionRequired,
				http.StatusTooManyRequests,
				http.StatusRequestHeaderFieldsTooLarge,
				http.StatusUnavailableForLegalReasons,
			},
			want: color.FgYellow,
		},
		{
			name: "Red",
			codes: []int{
				http.StatusInternalServerError,
				http.StatusNotImplemented,
				http.StatusBadGateway,
				http.StatusServiceUnavailable,
				http.StatusGatewayTimeout,
				http.StatusHTTPVersionNotSupported,
				http.StatusVariantAlsoNegotiates,
				http.StatusInsufficientStorage,
				http.StatusLoopDetected,
				http.StatusNotExtended,
				http.StatusNetworkAuthenticationRequired,
				999,
			},
			want: color.FgRed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, code := range tt.codes {
				if got := colorForStatus(code); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("colorForStatus() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
