package aml

import (
	"jonatak/aml/proto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AMLApi struct {
	grpcClient proto.PaymentsClient
}

func (a *AMLApi) PostTransaction(ctx *gin.Context) {
	var tx Transaction
	if err := ctx.ShouldBindJSON(&tx); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	txProto := &proto.TransactionRequest{}
	txProto.Transaction = tx.ToPbTransaction()
	txResponse, err := a.grpcClient.ApproveTransaction(ctx.Request.Context(), txProto)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var status ApiTransactionStatus
	var code int
	switch txResponse.Status {
	case proto.TransactionStatus_APPROVED:
		code, status = http.StatusAccepted, APPROVED
	case proto.TransactionStatus_INVALID_DATE:
		code, status = http.StatusForbidden, INVALID_DATE
	case proto.TransactionStatus_MAX_AMOUNT_REACH:
		code, status = http.StatusForbidden, MAX_AMOUNT_REACH
	}
	ctx.JSON(code, status)
}
