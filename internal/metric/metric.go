package metric

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

var (
	operationsTotalCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "operations_total",
	}, []string{"operation_type"})
)

func OperationsTotalCounterInc(operationType model.OperationType) {
	operationsTotalCounter.WithLabelValues(string(operationType)).Inc()
}

func MustRegister(registry *prometheus.Registry) {
	registry.MustRegister(
		operationsTotalCounter,
	)

	operationsTotalCounter.WithLabelValues(string(model.OperationTypeDebit)).Add(0)
	operationsTotalCounter.WithLabelValues(string(model.OperationTypeCredit)).Add(0)
}
