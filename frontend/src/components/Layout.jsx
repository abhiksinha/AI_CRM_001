import React from "react";

export function Layout({ auth, onLogout, currentPage, onNavigate, children }) {
  const menuItems = [
    { id: "dashboard", label: "Dashboard", icon: "📊" },
    { id: "contacts", label: "Contacts", icon: "👥" },
    { id: "deals", label: "Deals", icon: "🤝" },
    { id: "tasks", label: "Tasks & Notes", icon: "✅" },
    { id: "insights", label: "Insights", icon: "🧠" },
    { id: "settings", label: "Settings", icon: "⚙️" },
  ];

  return (
    <div className="app-shell layout-container">
      <aside className="sidebar">
        <div className="sidebar-header">
          <p className="tag">Atlas CRM</p>
        </div>
        <nav className="sidebar-nav">
          {menuItems.map((item) => (
            <button
              key={item.id}
              className={`nav-item ${currentPage === item.id ? "active" : ""}`}
              onClick={() => onNavigate(item.id)}
            >
              <span className="nav-icon">{item.icon}</span>
              {item.label}
            </button>
          ))}
        </nav>
        <div className="sidebar-footer">
          <div className="identity">
            <div className="avatar">{auth.user?.first_name?.[0] || "U"}</div>
            <div className="user-info">
              <span className="name">{auth.user?.first_name}</span>
              <span className="role">{auth.user?.role}</span>
            </div>
          </div>
          <button className="ghost logout-btn" onClick={onLogout}>
            Logout
          </button>
        </div>
      </aside>

      <main className="main-content">
        <header className="topbar">
          <h1>{menuItems.find((i) => i.id === currentPage)?.label}</h1>
        </header>
        <div className="content-scroll">{children}</div>
      </main>
    </div>
  );
}
