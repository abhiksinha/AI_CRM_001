import React from "react";

export function Dashboard({ overview }) {
  return (
    <section className="overview">
      <div className="card">
        <h3>Contacts</h3>
        <p className="metric">{overview.contacts}</p>
        <span className="muted">Total tracked in CRM</span>
      </div>
      <div className="card">
        <h3>Deals</h3>
        <p className="metric">{overview.deals}</p>
        <span className="muted">Active pipeline items</span>
      </div>
      <div className="card">
        <h3>Tasks</h3>
        <p className="metric">{overview.tasks}</p>
        <span className="muted">For selected deal</span>
      </div>
    </section>
  );
}
