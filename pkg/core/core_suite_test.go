package core

import (
	"log"
	"testing"
)

func TestCore(t *testing.T) {
	log.Println(prometheusResourceMemory)
	t.Log(prometheusResourceMemory)

	labels := map[string]string{
		kubeProxySelectorKey: kubeProxySelectorValue,
	}

	log.Println(labels)
}
