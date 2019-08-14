package controller

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/informers/apps/v1"
	listersv1 "k8s.io/client-go/listers/apps/v1"
	"time"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {
	kubeclientset     kubernetes.Interface
	deploymentLister  listersv1.DeploymentLister
	deploymentsSynced cache.InformerSynced
	workqueue         workqueue.RateLimitingInterface
}

func NewController(
	kubeclientset kubernetes.Interface,
	deploymentInformer v1.DeploymentInformer) *Controller {

	controller := &Controller{
		kubeclientset:     kubeclientset,
		deploymentLister:  deploymentInformer.Lister(),
		deploymentsSynced: deploymentInformer.Informer().HasSynced,
		workqueue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "deployment"),
	}

	glog.Info("Setting up event handlers")

	deploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueDeployment,
		UpdateFunc: func(old, new interface{}) {
			oldDeployment := old.(*appsv1.Deployment)
			newDeployment := new.(*appsv1.Deployment)
			if oldDeployment.ResourceVersion == newDeployment.ResourceVersion {
				return
			}
			controller.enqueueDeployment(new)
		},
		DeleteFunc: controller.enqueueDeploymentForDelete,
	})

	return controller
}

func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	glog.Info("Starting Deployment control loop")

	glog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.deploymentsSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	glog.Info("Starting workers")

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	glog.Info("Started workers")
	<-stopCh
	glog.Info("Shutting down workers")

	return nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	err := func(obj interface{}) error {

		defer c.workqueue.Done(obj)
		var key string
		var ok bool

		if key, ok = obj.(string); !ok {

			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}

		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}

		c.workqueue.Forget(obj)
		glog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

func (c *Controller) syncHandler(key string) error {

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	deployment, err := c.deploymentLister.Deployments(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			glog.Warningf("Deployment: %s/%s does not exist in local cache, will delete rs-pod ...",
				namespace, name)

			//TODO
			//这里就可以执行删除的业务逻辑
			glog.Infof("Deleting rs-pod: %s/%s ...", namespace, name)

			return nil
		}
		runtime.HandleError(fmt.Errorf("failed to list deployment by: %s/%s", namespace, name))
		return err
	}
	//TODO
	//这里通过比较期望状态和实际业务状态，执行创建或者更新的业务逻辑

	glog.Infof("Try to create rs-pod: %#v ...", deployment)

	return nil
}

func (c *Controller) enqueueDeployment(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}

func (c *Controller) enqueueDeploymentForDelete(obj interface{}) {
	var key string
	var err error
	key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}
