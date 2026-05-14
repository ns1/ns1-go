package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	api "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

func boolPtr(b bool) *bool {
	return &b
}

// This example demonstrates how to use the API key secret renewal endpoints.
// These endpoints allow you to:
// 1. View secret details (with manage_apikeys permission)
// 2. View your own secret details (self-service, no special permission)
// 3. Renew secrets to create new ones (with manage_apikeys permission)
// 4. Renew your own secret (self-service, no special permission)

func main() {
	// Create a new NS1 API client
	// The API key used here should have manage_apikeys permission for admin operations
	client := api.NewClient(nil, api.SetAPIKey("your-api-key-here"))

	// Example 1: Get details of a specific secret (requires manage_apikeys permission)
	fmt.Println("=== Example 1: Get Secret Details ===")
	secretID := "69a6c116a3edc300131c5fc0"
	secret, resp, err := client.APIKeys.GetSecret(secretID)
	if err != nil {
		log.Fatalf("Failed to get secret: %v", err)
	}
	fmt.Printf("Secret ID: %s\n", secret.ID)
	fmt.Printf("Enabled: %t\n", secret.Enabled)
	fmt.Printf("Expires At: %s\n", secret.ExpiresAt)
	if secret.LastAccess != "" {
		fmt.Printf("Last Access: %s\n", secret.LastAccess)
	}
	fmt.Printf("Response Status: %d\n\n", resp.StatusCode)

	// Example 2: Get details of the current secret (self-service, no special permission)
	fmt.Println("=== Example 2: Get Current Secret Details (Self) ===")
	currentSecret, resp, err := client.APIKeys.GetSecretSelf()
	if err != nil {
		log.Fatalf("Failed to get current secret: %v", err)
	}
	fmt.Printf("Current Secret ID: %s\n", currentSecret.ID)
	fmt.Printf("Enabled: %t\n", currentSecret.Enabled)
	fmt.Printf("Expires At: %s\n", currentSecret.ExpiresAt)

	// Check if secret is expiring soon (within 7 days)
	expiresAt, _ := time.Parse(time.RFC3339, currentSecret.ExpiresAt)
	daysUntilExpiry := time.Until(expiresAt).Hours() / 24
	fmt.Printf("Days until expiry: %.1f\n", daysUntilExpiry)
	if daysUntilExpiry < 7 {
		fmt.Println("WARNING: Secret is expiring soon! Consider renewing.")
	}
	fmt.Printf("Response Status: %d\n\n", resp.StatusCode)

	// Example 3: Renew a specific secret (requires manage_apikeys permission)
	fmt.Println("=== Example 3: Renew Secret ===")
	newSecret, resp, err := client.APIKeys.RenewSecret(secretID)
	if err != nil {
		log.Fatalf("Failed to renew secret: %v", err)
	}
	fmt.Printf("New Secret ID: %s\n", newSecret.ID)
	fmt.Printf("New Secret Key: %s\n", newSecret.Key) // SAVE THIS! Only visible once
	fmt.Printf("Enabled: %t\n", newSecret.Enabled)
	fmt.Printf("Expires At: %s\n", newSecret.ExpiresAt)
	fmt.Printf("Response Status: %d\n\n", resp.StatusCode)

	// IMPORTANT: Save the new secret key immediately!
	// This is the only time you'll see the plaintext secret value.
	fmt.Println("IMPORTANT: Save the new secret key immediately!")
	fmt.Println("    You can now use this new key for authentication.")
	fmt.Println("    The old secret remains active until it expires or is disabled.\n")

	// Example 4: Renew the current secret (self-service, no special permission)
	fmt.Println("=== Example 4: Renew Current Secret (Self) ===")
	renewedSecret, resp, err := client.APIKeys.RenewSecretSelf()
	if err != nil {
		log.Fatalf("Failed to renew current secret: %v", err)
	}
	fmt.Printf("Renewed Secret ID: %s\n", renewedSecret.ID)
	fmt.Printf("Renewed Secret Key: %s\n", renewedSecret.Key) // SAVE THIS! Only visible once
	fmt.Printf("Enabled: %t\n", renewedSecret.Enabled)
	fmt.Printf("Expires At: %s\n", renewedSecret.ExpiresAt)
	fmt.Printf("Response Status: %d\n\n", resp.StatusCode)

	// Example 5: Automated rotation workflow
	fmt.Println("=== Example 5: Automated Rotation Workflow ===")
	demonstrateAutomatedRotation(client)
}

// demonstrateAutomatedRotation shows a complete workflow for automated secret rotation
func demonstrateAutomatedRotation(client *api.Client) {
	// Step 1: Check current secret status
	currentSecret, _, err := client.APIKeys.GetSecretSelf()
	if err != nil {
		log.Printf("Failed to get current secret: %v", err)
		return
	}

	// Step 2: Parse expiration date
	expiresAt, err := time.Parse(time.RFC3339, currentSecret.ExpiresAt)
	if err != nil {
		log.Printf("Failed to parse expiration date: %v", err)
		return
	}

	// Step 3: Check if renewal is needed (e.g., within 7 days of expiry)
	daysUntilExpiry := time.Until(expiresAt).Hours() / 24
	fmt.Printf("Current secret expires in %.1f days\n", daysUntilExpiry)

	if daysUntilExpiry > 7 {
		fmt.Println("OK Secret is not expiring soon. No renewal needed.")
		return
	}

	fmt.Println("WARNING: Secret is expiring soon. Initiating renewal...")

	// Step 4: Renew the secret
	newSecret, _, err := client.APIKeys.RenewSecretSelf()
	if err != nil {
		log.Printf("Failed to renew secret: %v", err)
		return
	}

	fmt.Printf("OK New secret created: %s\n", newSecret.ID)
	fmt.Printf("  New secret key: %s\n", newSecret.Key)
	fmt.Printf("  Expires at: %s\n", newSecret.ExpiresAt)

	// Step 5: Save the new secret to your secure storage
	// In a real application, you would:
	// - Store the new secret in your secrets manager (e.g., AWS Secrets Manager, HashiCorp Vault)
	// - Update your application configuration
	// - Restart services to use the new secret
	fmt.Println("\nNext steps:")
	fmt.Println("   1. Save the new secret to your secrets manager")
	fmt.Println("   2. Update application configuration")
	fmt.Println("   3. Deploy/restart services with new secret")
	fmt.Println("   4. Monitor for successful authentication")
	fmt.Println("   5. Old secret will remain active until it expires")

	// Step 6: Optional - Disable the old secret after successful migration
	// (This would be done after verifying the new secret works)
	fmt.Println("\nNOTE: After verifying the new secret works:")
	fmt.Println("   - You can disable the old secret using UpdateSecret()")
	fmt.Println("   - Or let it expire naturally")
}

// Example error handling patterns
func demonstrateErrorHandling(client *api.Client) {
	fmt.Println("=== Error Handling Examples ===")

	// Handle missing secret
	_, _, err := client.APIKeys.GetSecret("non-existent-id")
	if err == api.ErrSecretMissing {
		fmt.Println("OK Correctly handled missing secret error")
	}

	// Handle invalid authentication for self-service endpoints
	_, _, err = client.APIKeys.GetSecretSelf()
	if err != nil && strings.Contains(err.Error(), "invalid authentication credentials") {
		fmt.Println("OK Correctly handled invalid authentication error")
	}

	// Handle max secrets limit
	_, _, err = client.APIKeys.RenewSecret("some-secret-id")
	if err != nil && err.Error() == "cannot renew secret: api key already has 2 active secrets" {
		fmt.Println("OK Correctly handled max secrets limit error")
		fmt.Println("  Solution: Disable or delete an old secret before renewing")
	}

	// Handle missing expiry_duration
	_, _, err = client.APIKeys.RenewSecretSelf()
	if err != nil && err.Error() == "apikey expiry_duration is required for secret renewal" {
		fmt.Println("OK Correctly handled missing expiry_duration error")
		fmt.Println("  Solution: API key must be created with expiry_duration set")
	}
}

// Example: Complete API key lifecycle with rotation
func demonstrateCompleteLifecycle(client *api.Client) {
	fmt.Println("=== Complete API Key Lifecycle ===")

	// Step 1: Create API key with rotation enabled
	newKey := &account.APIKey{
		Name:           "rotating-api-key",
		ExpiryDuration: "30d", // Secrets expire after 30 days
		Permissions: account.PermissionsMap{
			DNS: account.PermissionsDNS{
				ViewZones: true,
			},
		},
	}

	_, err := client.APIKeys.Create(newKey)
	if err != nil {
		log.Printf("Failed to create API key: %v", err)
		return
	}
	fmt.Printf("OK Created API key: %s\n", newKey.ID)
	fmt.Printf("  Initial secret: %s\n", newKey.Secrets[0].Key)
	fmt.Printf("  Expires at: %s\n", newKey.Secrets[0].ExpiresAt)

	// Step 2: Use the API key for some time...
	// (In production, this would be days/weeks of normal operation)

	// Step 3: Before expiry, renew the secret
	// Create a new client using the first secret
	rotatingClient := api.NewClient(nil, api.SetAPIKey(newKey.Secrets[0].Key))

	newSecret, _, err := rotatingClient.APIKeys.RenewSecretSelf()
	if err != nil {
		log.Printf("Failed to renew secret: %v", err)
		return
	}
	fmt.Printf("OK Renewed secret: %s\n", newSecret.ID)
	fmt.Printf("  New secret key: %s\n", newSecret.Key)

	// Step 4: Transition to new secret
	// Both secrets are now active - you have time to migrate
	fmt.Println("OK Both secrets are active during transition period")

	// Step 5: After migration, optionally disable old secret
	oldSecret := newKey.Secrets[0]
	oldSecret.Enabled = boolPtr(false)
	_, err = client.APIKeys.UpdateSecret(oldSecret)
	if err != nil {
		log.Printf("Failed to disable old secret: %v", err)
		return
	}
	fmt.Println("OK Disabled old secret after successful migration")

	// Step 6: Continue using new secret until next renewal
	fmt.Println("OK Lifecycle complete - ready for next rotation cycle")
}
