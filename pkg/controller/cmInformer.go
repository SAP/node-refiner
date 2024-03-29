package controller

import (
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

// AddConfigMapEventHandler subscribes and routes the different events of interest to the ConfigMap informer
func (c *WorkloadsController) AddConfigMapEventHandler() {
	c.cmInformer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.addConfigMap,
			UpdateFunc: c.updateConfigMap,
			DeleteFunc: c.deleteConfigMap,
		})
}

// addConfigMap notifies informer that a config map is added
// checks if the configMap added is for the controller settings
// updates settings if found
func (c *WorkloadsController) addConfigMap(obj interface{}) {
	cm := obj.(*corev1.ConfigMap)

	if cm.Name == "node-refiner-cm" {
		zap.S().Info("ConfigMap add event, initiating an update to the drainer settings")
		err := c.d.UpdateSettings(cm)
		if err != nil {
			zap.S().Warnw("Couldn't update the drainer settings using ConfigMap", "error", err)
		}

	}
}

// updateConfigMap notifies informer that a config map is updated
// checks if the configMap updated is for the controller settings
// updates settings if found
func (c *WorkloadsController) updateConfigMap(old, new interface{}) {
	// Cast the obj as ConfigMap
	cmNew := new.(*corev1.ConfigMap)
	if cmNew.Name == "node-refiner-cm" {
		zap.S().Info("ConfigMap update event, initiating an update to the drainer settings")
		err := c.d.UpdateSettings(cmNew)
		if err != nil {
			zap.S().Warnw("Couldn't update the drainer settings using ConfigMap", "error", err)
		}
	}
}

// deleteConfigMap notifies informer that a config map is deleted
func (c *WorkloadsController) deleteConfigMap(obj interface{}) {
	// Cast the obj as Node
	cm := obj.(*corev1.ConfigMap)
	zap.S().Infow("Deleted a config map", "name", cm.Name)
}
