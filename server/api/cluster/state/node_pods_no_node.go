package state

import (
	"net/http"

	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
)

type GetPodsOnNodeRequest struct {
	GetNodesRequest
	NodeName string `uri:"node"`
}

func GetPodsOnNode(c *gin.Context) {
	var request GetPodsOnNodeRequest
	c.ShouldBindUri(&request)

	shellReq := command.AnsibleCommandsRequest{
		Name:    "nodes",
		Command: "kubectl get pods --all-namespaces --field-selector spec.nodeName=" + request.NodeName,
	}
	shellResult, err := command.ExecuteShellCommandsAbortOnFirstSuccess("cluster", request.ClusterName, "kube_control_plane[0]", []command.AnsibleCommandsRequest{shellReq})

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed", err)
		return
	}

	c.JSON(http.StatusOK, common.KuboardSprayResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    shellResult[0],
	})
}
