package istio

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/system"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"github.com/openshift/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
)


const (
	application = "application"
	config      = "/config"
	yaml1        = "yaml"
)


func TestJsonToYaml(t *testing.T)  {
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
	out, err :=	yaml.Marshal(cfg)
	if err != nil {
		log.Print("test to yaml",err)
	}
	b := &Builder{
		Path:       os.TempDir()+"config",
		Name:       "application",
		FileType:   "yaml",
		Profile:    "dc",
		ConfigType: system.Configuration{},
	}
	b.WriterFile(out)
}

