package controllers

import (
	"backend/models"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/greenhandatsjtu/ISP_exp_platform/backend/database"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	typedAppsV1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	typedCoreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var (
	clientset        *kubernetes.Clientset
	metricsClientSet *metricsv.Clientset
)

// upload yaml files
func PostUploadYaml(c *gin.Context) {
	id := c.Params.ByName("eno")
	var experiment models.Experiment
	if err = database.Db.Where(id).First(&experiment).Error; err != nil {
		NotFound(c, "Experiment not found.")
		return
	}
	deployYaml, err := c.FormFile("deploy")
	if err != nil {
		badRequest(c)
		return
	}
	serviceYaml, err := c.FormFile("service")
	if err != nil {
		badRequest(c)
		return
	}
	mkdirIfNotExists(filepath.Join("uploads", id, "yaml"))
	if filepath.Ext(deployYaml.Filename) != ".yaml" || filepath.Ext(serviceYaml.Filename) != ".yaml" {
		badRequest(c)
		return
	}
	if err = c.SaveUploadedFile(deployYaml, filepath.Join("uploads", id, "yaml", id+".yaml")); err != nil {
		log.Println(err)
		return
	}
	if err = c.SaveUploadedFile(serviceYaml, filepath.Join("uploads", id, "yaml", id+"_svc.yaml")); err != nil {
		log.Println(err)
		return
	}
	database.Db.Model(&experiment).Update("upload", true)
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Code:    http.StatusCreated,
		Message: "Yaml uploaded.",
		Data:    nil,
	})
}

func StartExperiment(c *gin.Context) {
	eno := c.Params.ByName("eno")
	var experiment models.Experiment
	if err := database.Db.Where(eno).First(&experiment).Error; err != nil {
		NotFound(c, "Experiment not found.")
		return
	}
	user, _ := c.Get("user")
	var (
		enoInt int
		port   int32
	)
	if enoInt, err = strconv.Atoi(eno); err != nil {
		log.Println(err)
		badRequest(c)
		return
	}
	if port, err = startExperiment(eno, "experiment"+eno+"-"+strconv.Itoa(int(user.(models.User).ID))); err != nil {
		log.Println(err)
		if errors.IsNotFound(err) {
			NotFound(c, "Yaml not found.")
		} else {
			badRequest(c)
		}
		return
	}
	database.Db.Create(&models.UserResource{
		User:         user.(models.User),
		ExperimentID: uint(enoInt),
		Port:         port,
	})
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Code:    http.StatusCreated,
		Message: "",
		Data:    gin.H{"port": port},
	})
}

func EndExperiment(c *gin.Context) {
	eno := c.Params.ByName("eno")
	user, _ := c.Get("user")
	var resource models.UserResource
	if err = database.Db.Where(map[string]interface{}{"experiment_id": eno, "user_id": user.(models.User).ID}).First(&resource).Error; err != nil {
		log.Println(err)
		NotFound(c, "resource not found")
		return
	}
	if err = deleteExperiment("experiment" + eno + "-" + strconv.Itoa(int(user.(models.User).ID))); err != nil {
		log.Println(err)
		badRequest(c)
		return
	}
	database.Db.Delete(&resource)
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "Delete experiment successfully",
		Data:    nil,
	})
}

func AdminEndExperiment(c *gin.Context) {
	eno := c.Params.ByName("eno")
	id := c.Params.ByName("id")
	var resource models.UserResource
	if err = database.Db.Where(map[string]interface{}{"experiment_id": eno, "user_id": id}).First(&resource).Error; err != nil {
		log.Println(err)
		NotFound(c, "resource not found")
		return
	}
	if err = deleteExperiment("experiment" + eno + "-" + id); err != nil {
		log.Println(err)
		badRequest(c)
		return
	}
	database.Db.Delete(&resource)
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "Delete experiment successfully",
		Data:    nil,
	})
}

//获取同学占用实验资源情况
func GetAllResourceStatus(c *gin.Context) {
	var resources []models.UserResource
	if err = database.Db.Preload("Experiment").Preload("User").Find(&resources).Error; err != nil {
		log.Println(err)
		NotFound(c, "Resources not found")
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    resources,
	})
}

//获取单个实验的资源占用情况
func GetExperimentResourceStatus(c *gin.Context) {
	var resources []models.UserResource
	if eno, ok := c.Params.Get("eno"); ok {
		if err = database.Db.Where(map[string]interface{}{"experiment_id": eno}).Preload("Experiment").Preload("User").Find(&resources).Error; err != nil {
			log.Println(err)
			NotFound(c, "resource not found")
			return
		}
	} else {
		badRequest(c)
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    resources,
	})
}

//获取单个课程的资源占用情况
func GetCourseResourceStatus(c *gin.Context) {
	var resources []models.UserResource
	if cno, ok := c.Params.Get("cno"); ok {
		var course models.Course
		if err = database.Db.Where(cno).First(&course).Error; err != nil {
			log.Println(err)
			NotFound(c, "Course not found")
			return
		}
		experiments := course.GetExperiments()
		experiment_ids := make([]uint, len(experiments))
		for _, experiment := range experiments {
			experiment_ids = append(experiment_ids, experiment.ID)
		}
		if err = database.Db.Where("experiment_id IN (?)", experiment_ids).Preload("Experiment").Preload("User").Find(&resources).Error; err != nil {
			log.Println(err)
			NotFound(c, "resource not found")
			return
		}
	} else {
		badRequest(c)
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    resources,
	})
}

//获取某用户占用实验资源情况
func GetResourceStatus(c *gin.Context) {
	eno := c.Params.ByName("eno")
	user, _ := c.Get("user")
	var resource models.UserResource
	database.Db.Where(map[string]interface{}{"user_id": user.(models.User).ID, "experiment_id": eno}).First(&resource)
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    resource,
	})
}

func InitClient() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	metricsClientSet, err = metricsv.NewForConfig(config)
	if err != nil {
		panic(err)
	}
}

func startExperiment(config string, name string) (int32, error) {
	var (
		deployment appsv1.Deployment
		service    corev1.Service
		port       int32
	)
	deploymentConfigFile := filepath.Join("uploads", config, "yaml", config+".yaml")
	// yaml not exists
	if _, err = os.Stat(deploymentConfigFile); err != nil {
		return 0, err
	}
	serviceConfigFile := filepath.Join("uploads", config, "yaml", config+"_svc.yaml")
	if _, err = os.Stat(serviceConfigFile); err != nil {
		return 0, err
	}

	deploymentClient := clientset.AppsV1().Deployments(metav1.NamespaceDefault)
	serviceClient := clientset.CoreV1().Services(metav1.NamespaceDefault)
	if deployment, err = generateDeploymentFromFile(deploymentConfigFile); err != nil {
		return 0, err
	}
	deployment = renameDeployment(deployment, name)
	if err = createDeployment(deploymentClient, deployment); err != nil {
		return 0, err
	}

	if service, err = generateServiceFromFile(serviceConfigFile); err != nil {
		return 0, err
	}
	service = renameService(service, name)
	if port, err = createService(serviceClient, service); err != nil {
		return 0, err
	}
	return port, nil
}

func generateDeploymentFromFile(filename string) (appsv1.Deployment, error) {
	deployYaml, err := ioutil.ReadFile(filename)
	if err != nil {
		return appsv1.Deployment{}, err
	}
	deployJson, err := yaml.ToJSON(deployYaml)
	if err != nil {
		return appsv1.Deployment{}, err
	}
	var deployment appsv1.Deployment
	if err := json.Unmarshal(deployJson, &deployment); err != nil {
		return appsv1.Deployment{}, err
	}
	return deployment, nil
}

func generateServiceFromFile(filename string) (corev1.Service, error) {
	service := corev1.Service{}

	serviceYaml, err := ioutil.ReadFile(filename)
	if err != nil {
		return corev1.Service{}, err
	}

	serviceJson, err := yaml.ToJSON(serviceYaml)
	if err != nil {
		return corev1.Service{}, nil
	}

	if err = json.Unmarshal(serviceJson, &service); err != nil {
		return corev1.Service{}, nil
	}
	return service, nil
}

func renameDeployment(deployment appsv1.Deployment, name string) appsv1.Deployment {
	deployment.ObjectMeta.Name = name
	deployment.Spec.Selector.MatchLabels = map[string]string{"app": name}
	deployment.Spec.Template.Labels = map[string]string{"app": name}
	deployment.Spec.Template.Spec.Containers[0].Name = name
	//deployment.Spec.Template.ObjectMeta.Labels = map[string]string{"app": name}
	return deployment
}

func renameService(svc corev1.Service, name string) corev1.Service {
	svc.Spec.Selector = map[string]string{"app": name}
	name += "-svc"
	svc.ObjectMeta.Name = name
	return svc
}

func createDeployment(deploymentClient typedAppsV1.DeploymentInterface, deployment appsv1.Deployment) error {
	if _, err := deploymentClient.Get(context.TODO(), deployment.Name, metav1.GetOptions{}); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		fmt.Printf("Start creating deployment...\n")
		result, err := deploymentClient.Create(context.TODO(), &deployment, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		fmt.Printf("Deployment %s created!\n", result.Name)
	} else {
		fmt.Printf("Start updating deployment...\n")
		if _, err = deploymentClient.Update(context.TODO(), &deployment, metav1.UpdateOptions{}); err != nil {
			return err
		}
		fmt.Printf("Deployment %s updated!\n", deployment.Name)
	}
	return nil
}

func createService(serviceClient typedCoreV1.ServiceInterface, service corev1.Service) (int32, error) {
	if _, err := serviceClient.Get(context.TODO(), service.Name, metav1.GetOptions{}); err != nil {
		if !errors.IsNotFound(err) {
			return 0, err
		}
		fmt.Printf("Start creating service...\n")
		result, err := serviceClient.Create(context.TODO(), &service, metav1.CreateOptions{})
		if err != nil {
			return 0, err
		}
		fmt.Printf("Service %s created!\n", result.GetObjectMeta().GetName())
		return result.Spec.Ports[0].NodePort, nil
	} else {
		fmt.Printf("Start updating service...\n")
		var result *corev1.Service
		if result, err = serviceClient.Update(context.TODO(), &service, metav1.UpdateOptions{}); err != nil {
			return 0, err
		}
		fmt.Printf("Service %s updated!\n", result.Name)
		return result.Spec.Ports[0].NodePort, nil
	}
}

func deleteExperiment(name string) error {
	deploymentClient := clientset.AppsV1().Deployments(corev1.NamespaceDefault)
	serviceClient := clientset.CoreV1().Services(corev1.NamespaceDefault)

	if err = deleteDeployment(deploymentClient, name); err != nil {
		return err
	}
	if err = deleteService(serviceClient, name+"-svc"); err != nil {
		return err
	}
	return nil
}

func deleteDeployment(deploymentClient typedAppsV1.DeploymentInterface, deploymentName string) error {
	fmt.Printf("Start deleting deployment %s...\n", deploymentName)
	err := deploymentClient.Delete(context.TODO(), deploymentName, metav1.DeleteOptions{})
	return err
}

func deleteService(serviceClient typedCoreV1.ServiceInterface, serviceName string) error {
	fmt.Printf("Start deleting service %s...\n", serviceName)
	err := serviceClient.Delete(context.TODO(), serviceName, metav1.DeleteOptions{})
	return err
}

type Resource struct {
	Usage       string  `json:"usage"`
	Allocatable string  `json:"allocatable"`
	Percent     float64 `json:"percent"`
}

type Metric struct {
	Name   string   `json:"name"`
	Cpu    Resource `json:"cpu"`
	Memory Resource `json:"memory"`
}

//获取各节点资源占用情况
func GetNodeMetrics(c *gin.Context) {
	var metrics []Metric
	nodeClient := clientset.CoreV1().Nodes()
	list, _ := nodeClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return
	}
	nodeMetricsList, err := metricsClientSet.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return
	}
	for _, item := range list.Items {
		node, _ := nodeClient.Get(context.TODO(), item.Name, metav1.GetOptions{})
		for _, metric := range nodeMetricsList.Items {
			if node.Name == metric.Name {
				metrics = append(metrics, Metric{
					Name: node.Name,
					Cpu: Resource{
						Usage:       metric.Usage.Cpu().String(),
						Allocatable: node.Status.Allocatable.Cpu().String(),
						Percent:     float64(metric.Usage.Cpu().MilliValue()) / float64(node.Status.Allocatable.Cpu().MilliValue()),
					},
					Memory: Resource{
						Usage:       metric.Usage.Memory().String(),
						Allocatable: node.Status.Allocatable.Memory().String(),
						Percent:     float64(metric.Usage.Memory().Value()) / float64(node.Status.Allocatable.Memory().Value()),
					},
				})
				break
			}
		}
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "",
		Data:    metrics,
	})
}
