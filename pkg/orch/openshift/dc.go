// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.


package openshift


import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	"github.com/openshift/api/apps/v1"
	"github.com/hidevopsio/hicicd/pkg/orch/k8s"
	"github.com/hidevopsio/hiboot/pkg/log"
	"k8s.io/apimachinery/pkg/api/errors"
	"fmt"
	"github.com/jinzhu/copier"
	"os"
	"github.com/hidevopsio/hiboot/pkg/system"
	"gopkg.in/yaml.v2"
	"github.com/hidevopsio/hicicd/pkg/orch/istio"
	"path/filepath"
)

type DeploymentConfig struct {
	Name      string
	Namespace string
	FullName  string
	Version   string
	Interface appsv1.DeploymentConfigInterface
}

func NewDeploymentConfig(name, namespace, version string) (*DeploymentConfig, error) {
	log.Debug("NewDeploymentConfig()")

	clientSet, err := appsv1.NewForConfig(k8s.Config)
	if err != nil {
		return nil, err
	}
	fullName := name + "-" + version

	return &DeploymentConfig{
		Name:      name,
		Namespace: namespace,
		FullName:  fullName,
		Version: version,
		Interface: clientSet.DeploymentConfigs(namespace),
	}, nil
}

func (dc *DeploymentConfig) Create(env interface{}, ports interface{}, replicas int32, force bool, healthEndPoint string,
		cb func(cfg interface{}) error) error {
	log.Debug("DeploymentConfig.Create()")
	//privileged := true
	// env
	e := make([]corev1.EnvVar, 0)
	copier.Copy(&e, env)

	p := make([]corev1.ContainerPort, 0)
	copier.Copy(&p, ports)

	cfg := &v1.DeploymentConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: dc.FullName,
			Labels: map[string]string{
				"app": dc.Name,
				"version": dc.Version,
			},
		},
		Spec: v1.DeploymentConfigSpec{
			Replicas: replicas,

			Selector: map[string]string{
				"app": dc.Name,
				"version": dc.Version,
			},

			Strategy: v1.DeploymentStrategy{
				Type: v1.DeploymentStrategyTypeRolling,
			},

			Template: &corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: dc.Name,
					Labels: map[string]string{
						"app":  dc.Name,
						"version": dc.Version,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Env:             e,
							Image:           " ",
							ImagePullPolicy: corev1.PullAlways,
							Name:            dc.Name,
							Ports:           p,
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command : []string{
											"curl",
											"--silent",
											"--show-error",
											"--fail",
											healthEndPoint,
										},
									},
								},
								InitialDelaySeconds: 60,
								TimeoutSeconds:      1,
								PeriodSeconds:       60,
							},
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command : []string{
											"curl",
											"--silent",
											"--show-error",
											"--fail",
											healthEndPoint,
										},
									},
								},
								InitialDelaySeconds: 20,
								TimeoutSeconds:      1,
								PeriodSeconds:       10,
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
							dc.Name,
						},
						From: corev1.ObjectReference{
							Kind:      "ImageStreamTag",
							Name:      dc.Name + ":" + dc.Version,
							Namespace: dc.Namespace,
						},
					},
				},
			},
		},
	}
<<<<<<< HEAD
	// inject side car here
	err := dc.InjectSideCar(cfg)
	if err != nil {
		return err
	}
	result, err := dc.Interface.Get(dc.Name, metav1.GetOptions{})

=======
/*	err := cb(cfg)
	if err != nil {
		return err
	}*/
	// inject side car here
	result, err := dc.Interface.Get(dc.FullName, metav1.GetOptions{})

	err = cb(result)
>>>>>>> master

	switch {
	case err == nil:
		// select update or patch according to the user's request
		if force {
			cfg.ObjectMeta.ResourceVersion = result.ResourceVersion
			result, err = dc.Interface.Update(cfg)
			if err == nil {
				log.Infof("Updated DeploymentConfig %v.", result.Name)
			} else {
				log.Error("Updated  error %v", err)
				return err
			}
		}
	case errors.IsNotFound(err) :
		d, err := dc.Interface.Create(cfg)
		if err != nil {
			return err
		}
		log.Infof("Created DeploymentConfig %v.\n", d.Name)
	default:
		return fmt.Errorf("failed to create DeploymentConfig: %s", err)
	}

	return nil
}


func (dc *DeploymentConfig) Get() (*v1.DeploymentConfig, error) {
	log.Debug("DeploymentConfig.Get()")
	return dc.Interface.Get(dc.Name, metav1.GetOptions{})
}

func (dc *DeploymentConfig) Delete() error {
	log.Debug("DeploymentConfig.Delete()")
	return dc.Interface.Delete(dc.Name, &metav1.DeleteOptions{})
}

func (dc *DeploymentConfig) Instantiate() (*v1.DeploymentConfig, error)  {
	log.Debug("DeploymentConfig.Instantiate()")

	request := &v1.DeploymentRequest{
		TypeMeta: metav1.TypeMeta{
			Kind: "DeploymentRequest",
			APIVersion: "v1",
		},
		Name: dc.Name,
		Force: true,
		Latest: true,
	}

	d, err := dc.Interface.Instantiate(dc.Name, request)
	if nil == err {
		log.Infof("Instantiated Build %v", d.Name)
	}

	return d, err
}

func (dc *DeploymentConfig) InjectSideCar(cfg *v1.DeploymentConfig)  (error){
	in, err :=	yaml.Marshal(cfg)
	if err != nil {
		log.Print("test to yml",err)
		return err
	}
	b := &istio.Builder{
		Path:       os.TempDir(),
		Name:       dc.FullName,
		FileType:   "yml",
		ConfigType: system.Configuration{},
	}
	err = b.WriterFile(in)
	if err != nil {
		log.Print("test to yml",err)
		return err
	}
	path := filepath.Join(os.TempDir(), dc.FullName + ".yml")
	 err = istio.IntoResource(path, cfg)
	if err != nil {
		log.Print("test to yml",err)
	}
	return nil
}


