package urldelete_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/url/urldelete"
	"url-shortener/internal/http-server/handlers/url/urldelete/mocks"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
		},
		{
			name:      "Empty alias",
			alias:     "",
			respError: "invalid request",
		},
		{
			name:      "DeleteURL Error",
			alias:     "test_alias",
			respError: "internal error",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlDeleterMock := mocks.NewUrlDeleter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlDeleterMock.On("DeleteUrl", tc.alias).
					Return(tc.mockError).
					Once()
			}

			handler := urldelete.New(slogdiscard.NewDiscardLogger(), urlDeleterMock)

			uri := fmt.Sprintf("/url/{%s}", tc.alias)
			req, err := http.NewRequest(http.MethodDelete, uri, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("alias", tc.alias)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp urldelete.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			// TODO: add more checks
		})
	}
}
