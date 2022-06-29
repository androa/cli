package doctor

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/heroku/color"
	nais_io_v1alpha1 "github.com/nais/liberator/pkg/apis/nais.io/v1alpha1"
	"github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	iconError = "✗"
	iconOK    = "✓"
	iconDot   = "•"
	iconSkip  = "-"
)

type Manager struct {
	log           *logrus.Logger
	k8sClient     kubernetes.Interface
	dynamicClient dynamic.Interface
	app           *nais_io_v1alpha1.Application
	out           io.Writer
}

func New(log *logrus.Logger, cfg *rest.Config) (*Manager, error) {
	k8sClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	dynamicClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &Manager{
		log:           log,
		k8sClient:     k8sClient,
		dynamicClient: dynamicClient,
		out:           os.Stdout,
	}, nil
}

func (m *Manager) Init(ctx context.Context, namespace, appName string) error {
	kapp, err := m.dynamicClient.Resource(schema.GroupVersionResource{
		Group:    "nais.io",
		Version:  "v1alpha1",
		Resource: "applications",
	}).Namespace(namespace).Get(ctx, appName, metav1.GetOptions{})
	if err != nil {
		if kerrors.IsNotFound(err) {
			return fmt.Errorf("application %v not found", appName)
		}
		return err
	}

	m.app = &nais_io_v1alpha1.Application{}
	return runtime.DefaultUnstructuredConverter.FromUnstructured(kapp.Object, m.app)
}

func (m *Manager) SetOutput(w io.Writer) {
	m.out = w
}

func (m *Manager) Run(ctx context.Context, verbose bool) error {
	hasError := false

	fmt.Fprintln(m.out, "Running checks:")
	for _, check := range checks {
		cfg := &Config{
			Application:   m.app.DeepCopy(),
			K8sClient:     m.k8sClient,
			DynamicClient: m.dynamicClient,
			Log:           m.log.WithField("check", check.Name()),
			Out:           m.out,
		}
		fmt.Fprint(m.out, "  "+iconDot+" ", check.Name())
		if verbose {
			fmt.Fprintln(m.out)
		}
		err := check.Check(ctx, cfg)
		if err != nil {
			if errors.Is(err, ErrSkip) {
				if !verbose {
					fmt.Fprintln(m.out, " "+color.YellowString(iconSkip))
				}
				continue
			}
			hasError = true
			if !verbose {
				fmt.Fprintln(m.out, " "+iconError)
			}
			fmt.Fprintln(m.out, color.RedString(err.Error()))
		} else if !verbose {
			fmt.Fprintln(m.out, " "+color.GreenString(iconOK))
		}
	}

	if hasError {
		return fmt.Errorf("some checks failed")
	}
	return nil
}

func List(w io.Writer) {
	checks := checks[:]
	sort.Slice(checks, func(i, j int) bool {
		return checks[i].Name() < checks[j].Name()
	})
	for _, check := range checks {
		fmt.Fprintf(w, "  %v: %v\n", check.Name(), check.Help())
	}
}
