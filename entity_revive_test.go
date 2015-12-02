package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/context"
	"github.com/stretchr/testify/assert"
)

func TestReviveNoId(t *testing.T) {
	req := deleteRequest("/Shot/")

	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200, `foo`)
	defer server.Close()

	context.Set(req, "sgConn", *client)
	router(config).ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRevivePermissionsIssue(t *testing.T) {
	req := postRequest("/Project/75/revive", "")

	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200,
		`{"exception":true,"message":"API delete() CRUD ERROR #4.1: Entity Project 75 can not be `+
			`deleted by this user. Rule: API Admin -- PermissionRule 315: retire_entity_condition FOR `+
			`entity_type => Project.  RULE: {\"path\":\"name\", \"relation\":\"is_not\",\"values\":`+
			`[\"Template Project\"]}","error_code":104}`)
	defer server.Close()

	context.Set(req, "sgConn", *client)
	router(config).ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestReviveError(t *testing.T) {
	req := postRequest("/Project/75/revive", "")

	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200,
		`{"exception":true,"message":"Somer Error message","error_code":104}`)
	defer server.Close()

	context.Set(req, "sgConn", *client)
	router(config).ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestReviveSuccess(t *testing.T) {
	req := postRequest("/Project/75/revive", "")

	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200, `{"results": true}`)
	defer server.Close()

	context.Set(req, "sgConn", *client)
	router(config).ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestReviveMissing(t *testing.T) {
	req := postRequest("/Shot/1000000/revive", "")

	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200,
		`{"exception":true,"message":"API delete() CRUD ERROR #3: Entity of type [Shot] with id=1000000 does not exist.","error_code":104}`)
	defer server.Close()

	context.Set(req, "sgConn", *client)
	router(config).ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestReviveBadJsonResponse(t *testing.T) {
	req := postRequest("/Project/75/revive", "")

	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200, `foo`)
	defer server.Close()

	context.Set(req, "sgConn", *client)
	router(config).ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadGateway, w.Code)
}
