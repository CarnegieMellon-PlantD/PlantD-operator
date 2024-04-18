package utils

import "fmt"

// GetServiceARecord returns the DNS A record for a Service.
func GetServiceARecord(serviceName, serviceNamespace string) string {
	return fmt.Sprintf("%s.%s.svc", serviceName, serviceNamespace)
}

// GetServiceSRVRecord returns the DNS SRV record for a Service.
func GetServiceSRVRecord(portName, portProtocol, serviceName, serviceNamespace string) string {
	return fmt.Sprintf("_%s._%s.%s.%s.svc", portName, portProtocol, serviceName, serviceNamespace)
}
