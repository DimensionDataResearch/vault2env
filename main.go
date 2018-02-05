package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	vault "github.com/hashicorp/vault/api"
)

var (
	showHelp   bool
	secretPath string
	prefix     string
)

func main() {
	flag.BoolVar(&showHelp, "help", false, "Show usage information.")
	flag.StringVar(&secretPath, "secret-path", "", "The path, in vault, of the target secret.")
	flag.StringVar(&prefix, "prefix", "", "A prefix to add to each environment variable's name.")
	flag.Parse()

	if showHelp {
		flag.Usage()

		os.Exit(0)
	}

	if secretPath == "" {
		flag.Usage()

		os.Exit(1)
	}

	vaultAddress := os.Getenv(vault.EnvVaultAddress)
	if len(vaultAddress) == 0 {
		fmt.Printf("Must specify the Vault server address using the %s environment variable.\n", vault.EnvVaultAddress)

		os.Exit(1)
	}
	vaultToken := os.Getenv(vault.EnvVaultToken)
	if len(vaultToken) == 0 {
		fmt.Printf("Must specify the Vault access token using the %s environment variable.\n", vault.EnvVaultToken)

		os.Exit(1)
	}

	clientConfig := vault.DefaultConfig()
	clientConfig.Address = vaultAddress

	client, err := vault.NewClient(clientConfig)
	if err != nil {
		fmt.Println(err)

		os.Exit(2)
	}

	client.SetToken(vaultToken)

	secret, err := client.Logical().Read(secretPath)
	if err != nil {
		fmt.Println(err)

		os.Exit(2)
	}
	if secret == nil {
		fmt.Printf("Cannot find secret '%s'.", secretPath)

		os.Exit(3)
	}

	for key, value := range secret.Data {
		fmt.Printf("export %s%s='%s'\n",
			prefix,
			strings.ToUpper(key),
			value,
		)
	}
}
