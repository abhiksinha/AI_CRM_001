import React from "react";

export function SettingsPage({ createApiKey, apiKeyResponse, formatTimestamp }) {
  return (
    <div className="grid-actions">
      <div className="panel">
        <h2>Security & API Keys</h2>
        <p className="muted">
          Create a backend API key to use `/v1/token` for service-to-service
          sessions.
        </p>
        <button className="primary" onClick={createApiKey}>
          Create API key
        </button>
        {apiKeyResponse ? (
          <div className="key-result">
            <div>
              <strong>API Key ID:</strong> {apiKeyResponse.id}
            </div>
            <div>
              <strong>User ID:</strong> {apiKeyResponse.user_id}
            </div>
            <div>
              <strong>Created:</strong> {formatTimestamp(apiKeyResponse.created_at)}
            </div>
          </div>
        ) : null}
      </div>
    </div>
  );
}
