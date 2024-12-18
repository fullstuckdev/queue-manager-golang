package main

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"

    "queue-manager/internal/client"
)

var (
    apiClient *client.APIClient
    baseURL   string
)

var rootCmd = &cobra.Command{
    Use:   "tenant-cli",
    Short: "CLI for managing tenants",
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        apiClient = client.NewAPIClient(baseURL)
    },
}

var createCmd = &cobra.Command{
    Use:   "create",
    Short: "Create a new tenant",
    Run: func(cmd *cobra.Command, args []string) {
        clientID, _ := cmd.Flags().GetString("client-id")
        name, _ := cmd.Flags().GetString("name")
        
        fmt.Printf("Creating tenant with client ID: %s, name: %s\n", clientID, name)
        if err := apiClient.CreateTenant(clientID, name); err != nil {
            fmt.Printf("Error: %v\n", err)
            os.Exit(1)
        }
        fmt.Println("Tenant created successfully")
    },
}

var deleteCmd = &cobra.Command{
    Use:   "delete",
    Short: "Delete a tenant",
    Run: func(cmd *cobra.Command, args []string) {
        clientID, _ := cmd.Flags().GetString("client-id")
        
        fmt.Printf("Deleting tenant with client ID: %s\n", clientID)
        if err := apiClient.DeleteTenant(clientID); err != nil {
            fmt.Printf("Error: %v\n", err)
            os.Exit(1)
        }
        fmt.Println("Tenant deleted successfully")
    },
}

var processCmd = &cobra.Command{
    Use:   "process",
    Short: "Process a payload for a tenant",
    Run: func(cmd *cobra.Command, args []string) {
        clientID, _ := cmd.Flags().GetString("client-id")
        payload, _ := cmd.Flags().GetString("payload")
        
        fmt.Printf("Processing payload for client ID: %s\n", clientID)
        if err := apiClient.ProcessPayload(clientID, payload); err != nil {
            fmt.Printf("Error: %v\n", err)
            os.Exit(1)
        }
        fmt.Println("Payload processed successfully")
    },
}

func init() {
    rootCmd.PersistentFlags().StringVar(&baseURL, "api-url", "http://localhost:8080", "API base URL")
    
    createCmd.Flags().String("client-id", "", "Client ID for the tenant")
    createCmd.Flags().String("name", "", "Name of the tenant")
    createCmd.MarkFlagRequired("client-id")
    createCmd.MarkFlagRequired("name")

    deleteCmd.Flags().String("client-id", "", "Client ID of the tenant to delete")
    deleteCmd.MarkFlagRequired("client-id")

    processCmd.Flags().String("client-id", "", "Client ID of the tenant")
    processCmd.Flags().String("payload", "", "JSON payload to process")
    processCmd.MarkFlagRequired("client-id")
    processCmd.MarkFlagRequired("payload")

    rootCmd.AddCommand(createCmd, deleteCmd, processCmd)

    viper.AutomaticEnv()
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
} 