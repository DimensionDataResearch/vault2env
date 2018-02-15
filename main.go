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
	format     string
)

func main() {
	flag.BoolVar(&showHelp, "help", false, "Show usage information.")
	flag.StringVar(&secretPath, "secret-path", "", "The path, in vault, of the target secret.")
	flag.StringVar(&prefix, "prefix", "", "A prefix to add to each environment variable's name.")
	flag.StringVar(&prefix, "prefix", "", "A prefix to add to each environment variable's name.")
	flag.StringVar(&format, "format", "bash", "The environment variable format (bash, powershell, or powershell-env).")
	flag.Parse()

	if showHelp {
		flag.Usage()

		os.Exit(0)
	}

	var variableFormat string

	switch format {
	case "bash":
		variableFormat = "export %s%s='%s'"

	case "powershell":
		variableFormat = "$%s%s='%s'"

	case "powershell-env":
		variableFormat = "$env:%s%s='%s'"

	default:
		fmt.Printf("Unsupported format: '%s'.", format)

		os.Exit(4)
	}

	variableFormat += "\n"

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

	safeNameReplacer := strings.NewReplacer(
		"-", "_",
		".", "_",
		" ", "_",
	)
	for key, value := range secret.Data {
		safeName := strings.ToUpper(
			safeNameReplacer.Replace(key),
		)

		fmt.Printf(variableFormat, prefix, safeName, value)
	}
}
