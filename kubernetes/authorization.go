package main

// import (
// 	"context"

// 	"k8s.io/controller-manager/controller"
// 	"k8s.io/kubernetes/pkg/controller/clusterroleaggregation"
// )

// func startClusterRoleAggregrationController(ctx context.Context, controllerContext ControllerContext) (controller.Interface, bool, error) {
// 	go clusterroleaggregation.NewClusterRoleAggregation(
// 		controllerContext.InformerFactory.Rbac().V1().ClusterRoles(),
// 		controllerContext.ClientBuilder.ClientOrDie("clusterrole-aggregation-controller").RbacV1(),
// 	).Run(ctx, 5)
// 	return nil, true, nil
// }
