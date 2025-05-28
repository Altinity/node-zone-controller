package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/altinity/node-zone-controller/pkg/controller"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var kubeconfig string
	var contextName string
	var labelKey string
	
	defaultLabelKey := "altinity.cloud/auto-zone"
	if envLabelKey := os.Getenv("LABEL_KEY"); envLabelKey != "" {
		defaultLabelKey = envLabelKey
	}

	flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	flag.StringVar(&contextName, "context", "", "name of the kubeconfig context to use")
	flag.StringVar(&labelKey, "label-key", defaultLabelKey, "label key to watch for zone assignment (can also be set via LABEL_KEY env var)")
	flag.Parse()

	config, err := buildConfig(kubeconfig, contextName)
	if err != nil {
		slog.Error("Error building kubeconfig", "error", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		slog.Error("Error building kubernetes clientset", "error", err)
		os.Exit(1)
	}

	ctrl := controller.NewController(clientset, labelKey)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-stopCh
		slog.Info("Received termination signal, shutting down...")
		cancel()
	}()

	slog.Info("Starting node label controller", "labelKey", labelKey)
	if err := ctrl.Run(ctx, 2); err != nil {
		slog.Error("Error running controller", "error", err)
		os.Exit(1)
	}
}

func buildConfig(kubeconfig string, contextName string) (*rest.Config, error) {
	if kubeconfig != "" {
		loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig}
		configOverrides := &clientcmd.ConfigOverrides{CurrentContext: contextName}
		return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides).ClientConfig()
	}
	return rest.InClusterConfig()
}
