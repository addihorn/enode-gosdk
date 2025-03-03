package users_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/addihorn/enode-gosdk/pkg/auth"
	"github.com/addihorn/enode-gosdk/pkg/session"
	"github.com/addihorn/enode-gosdk/pkg/users"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/log"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

const usersResourcePath = "/users"

func TestUsersPact_GetUsers(t *testing.T) {

	log.SetLogLevel("INFO")
	mockProvider, err := consumer.NewV4Pact(consumer.MockHTTPProviderConfig{
		Consumer: "enode-go-sdk",
		Provider: "ENODE",
		LogDir:   os.Getenv("LOG_DIR"),
		PactDir:  os.Getenv("PACT_DIR"),
	})

	assert.NoError(t, err)

	t.Run("Users exist", func(t *testing.T) {
		// Given the Service does return a list of users
		expectedUsers := 1
		request := mockProvider.
			AddInteraction().Given("at least one user exists").
			UponReceiving("A request to get all users").
			WithRequest("GET", usersResourcePath, func(b *consumer.V4RequestBuilder) {
				b.Header("Authorization", matchers.Regex("Bearer abc123", "Bearer .+"))
			}).
			WillRespondWith(http.StatusOK, func(b *consumer.V4ResponseBuilder) {
				b.
					JSONBody(matchers.MapMatcher{
						"data": matchers.EachLike(matchers.MapMatcher{
							"id":        matchers.UUID(),
							"createdAt": matchers.DateTimeGenerated("2022-01-01T00:00:00Z", "yyyy-MM-dd'T'HH:mm:ss'Z'"),
						}, expectedUsers),
					}).
					Header("Content-Type", matchers.Term("application/json", `application\/json`))
			})

		result := request.ExecuteTest(t, func(config consumer.MockServerConfig) error {

			sess := &session.Session{
				Authentication: &auth.Authentication{
					Environment:  fmt.Sprintf("http://%s:%d", config.Host, config.Port),
					Access_token: "test-access-token",
				},
			}

			// Call the GetUsers function and capture the result
			usr, err := users.ListUsers(sess)
			assert.NoError(t, err)

			assert.EqualValues(t, expectedUsers, len(usr))

			return err

		})

		assert.NoError(t, result)
	})

	t.Run("No users exist", func(t *testing.T) {
		expectedUsers := 0

		request := mockProvider.
			AddInteraction().Given("no user exists").
			UponReceiving("A request to get all users").
			WithRequest("GET", usersResourcePath, func(b *consumer.V4RequestBuilder) {
				// b.Header("Authorization", matchers.Regex("Bearer abc123", "Bearer .+"))
			}).
			WillRespondWith(http.StatusOK, func(b *consumer.V4ResponseBuilder) {
				b.
					JSONBody(matchers.MapMatcher{
						"data": matchers.ArrayMaxLike(matchers.MapMatcher{
							"id":        matchers.UUID(),
							"createdAt": matchers.DateTimeGenerated("2022-01-01T00:00:00Z", "yyyy-MM-dd'T'HH:mm:ss'Z'"),
						}, expectedUsers),
					}).
					Header("Content-Type", matchers.Term("application/json", `application\/json`))
			})

		result := request.ExecuteTest(t, func(config consumer.MockServerConfig) error {

			sess := &session.Session{
				Authentication: &auth.Authentication{
					Environment:  fmt.Sprintf("http://%s:%d", config.Host, config.Port),
					Access_token: "test-access-token",
				},
			}

			// Call the GetUsers function and capture the result
			usr, err := users.ListUsers(sess)
			assert.NoError(t, err)

			assert.EqualValues(t, expectedUsers, len(usr))

			return err

		})

		assert.NoError(t, result)
	})

}
