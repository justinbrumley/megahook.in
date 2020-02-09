package main

import (
	"fmt"
	"testing"
)

func TestDBCreateNamespace(t *testing.T) {
	err := initDB()
	if err != nil {
		t.Errorf("Error initializing db: %v\n", err)
	}

	ns, err := createNamespace("12345", "my_test_namespace")
	if err != nil {
		t.Errorf("Error creating namespace: %v\n", err)
	}

	if ns == nil {
		t.Errorf("Namespace not returned from creation")
	}
}

func TestDBLookupNamespace(t *testing.T) {
	err := initDB()
	if err != nil {
		t.Errorf("Error initializing db: %v\n", err)
	}

	ns, err := lookupNamespace("my_test_namespace")
	if err != nil {
		t.Errorf("Error looking up namespace: %v\n", err)
	}

	if ns == nil {
		t.Errorf("Failed to look up namespace: Namespace not found.")
	}
}

func TestDBGetTokenNamespace(t *testing.T) {
	err := initDB()
	if err != nil {
		t.Errorf("Error initializing db: %v\n", err)
	}

	ns, err := getTokenNamespace("12345")
	if err != nil {
		t.Errorf("Error looking up token: %v\n", err)
	}

	if ns == nil {
		t.Errorf("Failed to look up token: Token namespace not found.")
	}
}

func TestDBDeleteNamespace(t *testing.T) {
	err := initDB()
	if err != nil {
		t.Errorf("Error initializing db: %v\n", err)
	}

	err = deleteNamespace("12345")
	if err != nil {
		t.Errorf("Error deleting namespace: %v\n", err)
	}
}
