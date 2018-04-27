package istio

import (
	"testing"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/openshift/api/apps/v1"
	"log"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/kube/inject"
	"path/filepath"
	"os"
	"github.com/magiconair/properties/assert"
	"fmt"
	"gopkg.in/yaml.v2"
)

func TestIntoResourceFile(t *testing.T)  {

/*	c := &resource{
		in: filepath.Join(os.TempDir(), "in.yaml"),
	}*/
	path := filepath.Join(os.TempDir(), "in.yaml")
	var dc v1.DeploymentConfig
	 err :=IntoResource(path, &dc)
	log.Println(dc)
	assert.Equal(t, nil, err)

}

func TestIntoObject(t *testing.T) {
	debugMode := true
	mesh := model.DefaultMeshConfig()
	params := &inject.Params{
		InitImage:           inject.InitImageName(Hub, Tag, debugMode),
		ProxyImage:          inject.ProxyImageName(Hub, Tag, debugMode),
		ImagePullPolicy:     "IfNotPresent",
		Verbosity:           DefaultVerbosity,
		SidecarProxyUID:     DefaultSidecarProxyUID,
		Version:             "12345678",
		EnableCoreDump:      false,
		Mesh:                &mesh,
		DebugMode:           debugMode,
	}

	sidecarTemplate, err := inject.GenerateTemplateFromParams(params)
	if err != nil {
		t.Fatalf("GenerateTemplateFromParams(%v) failed: %v", params, err)
	}
	name := "foo"
	cfg := &v1.DeploymentConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"app": name,
			},
		},
		Spec: v1.DeploymentConfigSpec{
			Replicas: 1,

			Selector: map[string]string{
				"app": name,
			},

			Strategy: v1.DeploymentStrategy{
				Type: v1.DeploymentStrategyTypeRolling,
			},

			Template: &corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: name,
					Labels: map[string]string{
						"app":  name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							//Env:             e,
							Image:           " ",
							ImagePullPolicy: corev1.PullAlways,
							Name:            name,
							//Ports:           p,
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command : []string{
											"curl",
											"--silent",
											"--show-error",
											"--fail",
											"http://localhost:8080/health",
										},
									},
								},
								InitialDelaySeconds: 10,
								TimeoutSeconds:      1,
								PeriodSeconds:       5,
							},
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command : []string{
											"curl",
											"--silent",
											"--show-error",
											"--fail",
											"http://localhost:8080/health",
										},
									},
								},
								InitialDelaySeconds: 20,
								TimeoutSeconds:      1,
								PeriodSeconds:       5,
							},
						},
					},
					DNSPolicy:     corev1.DNSClusterFirst,
					RestartPolicy: corev1.RestartPolicyAlways,
					SchedulerName: "default-scheduler",
				},
			},
			Test: false,
			Triggers: v1.DeploymentTriggerPolicies{
				{
					Type: v1.DeploymentTriggerOnImageChange,
					ImageChangeParams: &v1.DeploymentTriggerImageChangeParams{
						Automatic: true,
						ContainerNames: []string{
							name,
						},
						From: corev1.ObjectReference{
							Kind:      "ImageStreamTag",
							Name:      name + ":" + "latest",
							Namespace: "demo-dev",
						},
					},
				},
			},
		},
	}
	log.Println(cfg)

	out, err := inject.IntoObject(sidecarTemplate, &mesh, nil)
	d := out.(*v1.DeploymentConfig)

	// patch to fix the issue that injection failed on OpenShift
	privileged := true
	d.Spec.Template.Spec.InitContainers[0].SecurityContext.Privileged = &privileged

	log.Print(d)
}


var data = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`

// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.
type T struct {
	A string
	B struct {
		RenamedC int   `yaml:"c"`
		D        []int `yaml:",flow"`
	}
}

func TestYaml(t1 *testing.T) {
	t := T{}
	in := []byte(data)
	err := yaml.Unmarshal(in, &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t:\n%v\n\n", t)

	d, err := yaml.Marshal(&t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t dump:\n%s\n\n", string(d))

	m := make(map[interface{}]interface{})

	err = yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- m:\n%v\n\n", m)

	d, err = yaml.Marshal(&m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- m dump:\n%s\n\n", string(d))
}