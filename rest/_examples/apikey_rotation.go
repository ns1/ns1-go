package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	api "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

var client *api.Client

func boolPtr(b bool) *bool {
	return &b
}

// Helper that initializes rest api client from environment variable.
func init() {
	k := os.Getenv("NS1_APIKEY")
	if k == "" {
		fmt.Println("NS1_APIKEY environment variable is not set, giving up")
	}

	httpClient := &http.Client{Timeout: time.Second * 10}
	// Adds logging to each http request.
	doer := api.Decorate(httpClient, api.Logging(log.New(os.Stdout, "", log.LstdFlags)))
	client = api.NewClient(doer, api.SetAPIKey(k))
}

func main() {
	// Example 1: Create an API key with automatic secret rotation
	fmt.Println("=== Creating API Key with Rotation ===")
	newKey := &account.APIKey{
		Name:           "rotating-api-key-example",
		ExpiryDuration: "30d", // Secrets expire every 30 days (options: "10d", "30d", "90d")
		Permissions: account.PermissionsMap{
			DNS: account.PermissionsDNS{
				ViewZones:           true,
				ManageZones:         true,
				ZonesAllowByDefault: true,
			},
		},
	}

	if _, err := client.APIKeys.Create(newKey); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created API key: %s (ID: %s)\n", newKey.Name, newKey.ID)
	fmt.Printf("Expiry duration: %s\n", newKey.ExpiryDuration)

	// The initial secret is returned in the Secrets array
	if len(newKey.Secrets) > 0 {
		fmt.Printf("Initial secret created:\n")
		fmt.Printf("  Secret ID: %s\n", newKey.Secrets[0].ID)
		fmt.Printf("  Secret Key: %s\n", newKey.Secrets[0].Key) // Only shown on creation!
		fmt.Printf("  Expires At: %s\n", newKey.Secrets[0].ExpiresAt)
		if newKey.Secrets[0].Enabled != nil {
			fmt.Printf("  Enabled: %t\n", *newKey.Secrets[0].Enabled)
		}
	}

	// Example 2: Get an API key and view its secrets
	fmt.Println("\n=== Getting API Key with Secrets ===")
	retrievedKey, _, err := client.APIKeys.Get(newKey.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("API Key: %s\n", retrievedKey.Name)
	fmt.Printf("Number of secrets: %d\n", len(retrievedKey.Secrets))
	for i, secret := range retrievedKey.Secrets {
		fmt.Printf("\nSecret %d:\n", i+1)
		fmt.Printf("  ID: %s\n", secret.ID)
		fmt.Printf("  Expires At: %s\n", secret.ExpiresAt)
		fmt.Printf("  Last Access: %s\n", secret.LastAccess)
		if secret.Enabled != nil {
			fmt.Printf("  Enabled: %t\n", *secret.Enabled)
		}
		// Note: The actual secret key is NOT returned after creation
	}

	// Example 3: List all API keys and show which have rotation enabled
	fmt.Println("\n=== Listing API Keys with Rotation Status ===")
	keys, _, err := client.APIKeys.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, k := range keys {
		if k.ExpiryDuration != "" {
			fmt.Printf("Key: %s (ID: %s)\n", k.Name, k.ID)
			fmt.Printf("  Rotation: %s\n", k.ExpiryDuration)
			fmt.Printf("  Active Secrets: %d\n", len(k.Secrets))
		}
	}

	// Example 4: Update a secret (disable it)
	if len(retrievedKey.Secrets) > 0 {
		fmt.Println("\n=== Disabling a Secret ===")
		secretToUpdate := retrievedKey.Secrets[0]

		// Modify the secret
		secretToUpdate.Enabled = boolPtr(false)

		// Update it
		if _, err := client.APIKeys.UpdateSecret(secretToUpdate); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Disabled secret: %s\n", secretToUpdate.ID)
		if secretToUpdate.Enabled != nil {
			fmt.Printf("  Enabled: %t\n", *secretToUpdate.Enabled)
		}
	}

	// Example 5: Update a secret's expiration date
	if len(retrievedKey.Secrets) > 0 {
		fmt.Println("\n=== Extending Secret Expiration ===")
		secretToExtend := retrievedKey.Secrets[0]

		// Extend expiration by setting a new date
		secretToExtend.ExpiresAt = "2026-12-31"

		// Update it
		if _, err := client.APIKeys.UpdateSecret(secretToExtend); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Extended secret expiration: %s\n", secretToExtend.ID)
		fmt.Printf("  New Expires At: %s\n", secretToExtend.ExpiresAt)
	}

	// Example 6: Re-enable a secret
	if len(retrievedKey.Secrets) > 0 {
		fmt.Println("\n=== Re-enabling a Secret ===")
		secretToEnable := retrievedKey.Secrets[0]

		// Re-enable the secret
		secretToEnable.Enabled = boolPtr(true)

		// Update it
		if _, err := client.APIKeys.UpdateSecret(secretToEnable); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Re-enabled secret: %s\n", secretToEnable.ID)
		if secretToEnable.Enabled != nil {
			fmt.Printf("  Enabled: %t\n", *secretToEnable.Enabled)
		}
	}

	// Example 7: Delete a secret
	// Note: You should have at least one active secret before deleting others
	if len(retrievedKey.Secrets) > 1 {
		fmt.Println("\n=== Deleting a Secret ===")
		secretToDelete := retrievedKey.Secrets[1].ID

		if _, err := client.APIKeys.DeleteSecret(secretToDelete); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Deleted secret: %s\n", secretToDelete)
	}

	// Example 8: Get secret details (admin operation)
	if len(retrievedKey.Secrets) > 0 {
		fmt.Println("\n=== Getting Secret Details ===")
		secretID := retrievedKey.Secrets[0].ID
		secretDetails, _, err := client.APIKeys.GetSecret(secretID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Secret ID: %s\n", secretDetails.ID)
		fmt.Printf("Enabled: %t\n", *secretDetails.Enabled)
		fmt.Printf("Expires At: %s\n", secretDetails.ExpiresAt)
		if secretDetails.LastAccess != "" {
			fmt.Printf("Last Access: %s\n", secretDetails.LastAccess)
		}
	}

	// Example 9: Renew a secret (admin operation)
	var newSecret *account.APIKeySecret
	if len(retrievedKey.Secrets) > 0 {
		fmt.Println("\n=== Renewing a Secret ===")
		secretID := retrievedKey.Secrets[0].ID
		newSecret, _, err = client.APIKeys.RenewSecret(secretID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("New Secret Created:\n")
		fmt.Printf("  Secret ID: %s\n", newSecret.ID)
		fmt.Printf("  Secret Key: %s\n", newSecret.Key) // Only visible once!
		fmt.Printf("  Expires At: %s\n", newSecret.ExpiresAt)
		fmt.Printf("  Enabled: %t\n", *newSecret.Enabled)
		fmt.Println("\nIMPORTANT: Save the secret key immediately - it won't be shown again!")
	}

	// Example 10: Self-service operations (using the renewed secret)
	// Note: These operations allow an API key to manage itself
	if newSecret != nil && newSecret.Key != "" {
		fmt.Println("\n=== Self-Service Operations ===")
		fmt.Println("(Demonstrating with the newly created secret)")

		// Create a new client using the renewed secret
		secretClient := api.NewClient(
			&http.Client{Timeout: time.Second * 10},
			api.SetAPIKey(newSecret.Key),
		)

		// Get own secret details
		selfSecret, _, err := secretClient.APIKeys.GetSecretSelf()
		if err != nil {
			fmt.Printf("GetSecretSelf error: %v\n", err)
		} else {
			fmt.Printf("Current Secret ID: %s\n", selfSecret.ID)
			fmt.Printf("Expires At: %s\n", selfSecret.ExpiresAt)
			if selfSecret.Enabled != nil {
				fmt.Printf("Enabled: %t\n", *selfSecret.Enabled)
			}
		}

		// Note: RenewSecretSelf would create another secret, but we'll skip it
		// to avoid creating too many secrets in this example
		fmt.Println("\n(Skipping RenewSecretSelf to avoid creating too many secrets)")
	}

	// Example 11: Pretty print the final state
	fmt.Println("\n=== Final API Key State ===")
	finalKey, _, err := client.APIKeys.Get(newKey.ID)
	if err != nil {
		log.Fatal(err)
	}

	b, _ := json.MarshalIndent(finalKey, "", "  ")
	fmt.Println(string(b))

	// Cleanup: Delete the example API key
	fmt.Println("\n=== Cleanup ===")
	if _, err := client.APIKeys.Delete(newKey.ID); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted API key: %s\n", newKey.ID)
}
