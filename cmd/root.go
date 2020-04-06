package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/danniel1205/sample-webhook-server/pkg/webhook"
	"github.com/spf13/cobra"

	"k8s.io/klog"
)

var CmdWebhook = &cobra.Command{
	Use:   "",
	Short: "Starts a HTTP server, useful for testing MutatingAdmissionWebhook and ValidatingAdmissionWebhook",
	Long: `Starts a HTTP server, useful for testing MutatingAdmissionWebhook and ValidatingAdmissionWebhook.
After deploying it to Kubernetes cluster, the Administrator needs to create a ValidatingWebhookConfiguration
in the Kubernetes cluster to register remote webhook admission controllers.`,
	Args: cobra.MaximumNArgs(0),
	Run:  webhook.RunCmd,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := CmdWebhook.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	loggingFlags := &flag.FlagSet{}
	klog.InitFlags(loggingFlags)

	CmdWebhook.PersistentFlags().AddGoFlagSet(loggingFlags)
	CmdWebhook.Flags().String("tls-cert-file", "",
		"File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated after server cert).")
	CmdWebhook.Flags().String("tls-private-key-file", "",
		"File containing the default x509 private key matching --tls-cert-file.")
	CmdWebhook.Flags().Int("port", 443,
		"Secure port that the webhook listens on")
	CmdWebhook.Flags().String("sidecar-image", "",
		"Image to be used as the injected sidecar")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
