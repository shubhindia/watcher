package main

import (
	"fmt"
	"os"

	"github.com/shubhindia/watcher/config"
	"github.com/shubhindia/watcher/controllers"
	"k8s.io/apimachinery/pkg/runtime"

	corev1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

func main() {

	filename := "config.yaml"
	config, err := config.LoadConfig(filename)
	if err != nil {
		fmt.Println("Error reading config file", err)
		os.Exit(1)
	}

	scheme := runtime.NewScheme()
	utilruntime.Must(corev1.AddToScheme(scheme))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
	})
	if err != nil {
		fmt.Println("Unable to create manager", err)
		os.Exit(1)
	}

	if err := (&controllers.PodWatcherReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Config: config,
	}).SetupWithManager(mgr); err != nil {
		fmt.Println("Unable to create controller", err)
		os.Exit(1)
	}

	// Start the Manager
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		fmt.Println("Unable to start manager", err)
		os.Exit(1)
	}

}
