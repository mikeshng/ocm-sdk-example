package main

import (
	"fmt"
	"log"
	"os"

	sdk "github.com/openshift-online/ocm-sdk-go"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

func main() {
	logger, err := sdk.NewGoLoggerBuilder().
		Debug(false).
		Build()
	if err != nil {
		log.Fatalf("can't build logger: %v\n", err)
	}

	connection, err := sdk.NewConnectionBuilder().
		Logger(logger).
		Tokens(os.Getenv("OFFLINE_ACCESS_TOKEN")).
		URL("http://localhost:9000").
		Insecure(true).
		Build()
	if err != nil {
		log.Fatalf("Can't build connection: %v", err)
	}
	defer connection.Close()

	clusterId := "REDACTED"

	autoscaler, err := cmv1.NewClusterAutoscaler().
		ResourceLimits(cmv1.NewAutoscalerResourceLimits().MaxNodesTotal(10)).
		Build()
	postResponse, err := connection.ClustersMgmt().V1().Clusters().Cluster(clusterId).Autoscaler().
		Post().Request(autoscaler).Send()

	fmt.Printf("1st Post response: %d \n", postResponse.Status())

	getResponse, err := connection.ClustersMgmt().V1().Clusters().Cluster(clusterId).Autoscaler().Get().Send()
	fmt.Printf("Get response: %d %d \n", getResponse.Status(), getResponse.Body().ResourceLimits().MaxNodesTotal())

	autoscaler, err = cmv1.NewClusterAutoscaler().
		ResourceLimits(cmv1.NewAutoscalerResourceLimits().MaxNodesTotal(10)).
		Build()
	postResponse, err = connection.ClustersMgmt().V1().Clusters().Cluster(clusterId).Autoscaler().
		Post().Request(autoscaler).Send()

	fmt.Printf("2nd Post response: %d %s \n", postResponse.Status(), postResponse.Error())

	patch, err := cmv1.NewClusterAutoscaler().
		ResourceLimits(cmv1.NewAutoscalerResourceLimits().MaxNodesTotal(20)).
		Build()
	updateResponse, err := connection.ClustersMgmt().V1().Clusters().Cluster(clusterId).Autoscaler().
		Update().Body(patch).Send()

	fmt.Printf("Update response: %d %d \n", updateResponse.Status(), updateResponse.Body().ResourceLimits().MaxNodesTotal())

	deleteResponse, err := connection.ClustersMgmt().V1().Clusters().Cluster(clusterId).Autoscaler().Delete().Send()
	fmt.Printf("Delete response: %d \n", deleteResponse.Status())
}
