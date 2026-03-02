import { useEffect, useMemo, useState } from "react";

const EDGE_BASE_URL =
  import.meta.env.VITE_EDGE_BASE_URL || "http://127.0.0.1:8081";

const COOKIE_TOKEN = "crm_auth_token";
const COOKIE_USER = "crm_user_id";

const defaultSignup = {
  first_name: "",
  last_name: "",
  email: "",
  password: "",
  role: "user"
};

const defaultContact = {
  first_name: "",
  last_name: "",
  email: "",
  phone: "",
  owner_id: ""
};

const defaultDeal = {
  name: "",
  stage: "prospect",
  value: "",
  expected_close_date: "",
  contact_id: "",
  owner_id: ""
};

const defaultTask = {
  deal_id: "",
  title: "",
  due_date: "",
  is_completed: false,
  assigned_to_id: ""
};

const defaultNote = {
  contact_id: "",
  content: ""
};

const defaultInsight = {
  lead_id: "",
  contact_id: ""
};

function getCookie(name) {
  return document.cookie
    .split(";")
    .map((cookie) => cookie.trim())
    .filter((cookie) => cookie.startsWith(`${name}=`))
    .map((cookie) => decodeURIComponent(cookie.split("=")[1]))[0];
}

function setCookie(name, value, maxAgeSeconds) {
  const maxAge = maxAgeSeconds ? `; Max-Age=${maxAgeSeconds}` : "";
  document.cookie = `${name}=${encodeURIComponent(
    value
  )}${maxAge}; Path=/; SameSite=Lax`;
}

function clearCookie(name) {
  document.cookie = `${name}=; Max-Age=0; Path=/; SameSite=Lax`;
}

async function apiRequest(path, { method = "GET", body, token } = {}) {
  const headers = {
    "Content-Type": "application/json"
  };
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const response = await fetch(`${EDGE_BASE_URL}${path}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined
  });

  if (response.status === 204) {
    return null;
  }

  const data = await response.json().catch(() => null);
  if (!response.ok) {
    const message =
      data?.description ||
      data?.message ||
      data?.error ||
      response.statusText ||
      "Request failed";
    const error = new Error(message);
    error.status = response.status;
    error.payload = data;
    throw error;
  }
  return data;
}

function formatTimestamp(ts) {
  if (!ts) return "-";
  const date = new Date(ts * 1000);
  if (Number.isNaN(date.getTime())) return "-";
  return date.toLocaleDateString();
}

function App() {
  const [auth, setAuth] = useState({
    status: "checking",
    token: null,
    userId: null,
    user: null
  });
  const [authError, setAuthError] = useState("");

  const [loginForm, setLoginForm] = useState({
    username: "",
    password: ""
  });
  const [signupForm, setSignupForm] = useState(defaultSignup);

  const [contacts, setContacts] = useState([]);
  const [deals, setDeals] = useState([]);
  const [overview, setOverview] = useState({
    contacts: 0,
    deals: 0,
    tasks: 0
  });
  const [activeDealTasks, setActiveDealTasks] = useState([]);

  const [contactForm, setContactForm] = useState(defaultContact);
  const [dealForm, setDealForm] = useState(defaultDeal);
  const [taskForm, setTaskForm] = useState(defaultTask);
  const [noteForm, setNoteForm] = useState(defaultNote);
  const [insightForm, setInsightForm] = useState(defaultInsight);
  const [apiKeyResponse, setApiKeyResponse] = useState(null);

  const [actionStatus, setActionStatus] = useState({
    loading: false,
    message: "",
    tone: ""
  });

  const authHeader = useMemo(() => auth.token, [auth.token]);

  useEffect(() => {
    const existingToken = getCookie(COOKIE_TOKEN);
    const existingUserId = getCookie(COOKIE_USER);
    if (!existingToken || !existingUserId) {
      setAuth({ status: "logged_out", token: null, userId: null, user: null });
      return;
    }

    const init = async () => {
      try {
        const user = await apiRequest(`/v1/users/${existingUserId}`, {
          token: existingToken
        });
        setAuth({
          status: "logged_in",
          token: existingToken,
          userId: existingUserId,
          user
        });
      } catch (error) {
        clearCookie(COOKIE_TOKEN);
        clearCookie(COOKIE_USER);
        setAuth({ status: "logged_out", token: null, userId: null, user: null });
      }
    };

    init();
  }, []);

  useEffect(() => {
    if (auth.status !== "logged_in") return;
    setContactForm((prev) => ({ ...prev, owner_id: auth.userId || "" }));
    setDealForm((prev) => ({ ...prev, owner_id: auth.userId || "" }));
    setTaskForm((prev) => ({ ...prev, assigned_to_id: auth.userId || "" }));
  }, [auth.status, auth.userId]);

  useEffect(() => {
    if (auth.status !== "logged_in") return;
    refreshOverview();
  }, [auth.status]);

  const refreshOverview = async () => {
    try {
      const [contactsResponse, dealsResponse] = await Promise.all([
        apiRequest("/v1/contacts?page=1&page_size=6", { token: authHeader }),
        apiRequest("/v1/deals?page=1&page_size=6", { token: authHeader })
      ]);
      setContacts(contactsResponse?.data || []);
      setDeals(dealsResponse?.data || []);
      setOverview({
        contacts: contactsResponse?.total_count ?? 0,
        deals: dealsResponse?.total_count ?? 0,
        tasks: activeDealTasks.length
      });
    } catch (error) {
      setActionStatus({
        loading: false,
        message: error.message || "Failed to load data",
        tone: "danger"
      });
    }
  };

  const loginWithCredentials = async (username, password) => {
    const login = await apiRequest("/v1/login", {
      method: "POST",
      body: { username, password }
    });
    setCookie(COOKIE_TOKEN, login.api_token, login.expires_in);
    setCookie(COOKIE_USER, login.user_id, login.expires_in);

    const user = await apiRequest(`/v1/users/${login.user_id}`, {
      token: login.api_token
    });

    setAuth({
      status: "logged_in",
      token: login.api_token,
      userId: login.user_id,
      user
    });
  };

  const handleLogin = async (event) => {
    event.preventDefault();
    setAuthError("");
    try {
      await loginWithCredentials(loginForm.username, loginForm.password);
      setLoginForm({ username: "", password: "" });
    } catch (error) {
      setAuthError(error.message || "Login failed");
    }
  };

  const handleSignup = async (event) => {
    event.preventDefault();
    setAuthError("");
    try {
      await apiRequest("/v1/users", {
        method: "POST",
        body: signupForm
      });
      await loginWithCredentials(signupForm.email, signupForm.password);
    } catch (error) {
      setAuthError(error.message || "Signup failed");
    }
  };

  const handleLogout = async () => {
    try {
      if (auth.token) {
        await apiRequest("/v1/logout", {
          method: "POST",
          body: { api_token: auth.token }
        });
      }
    } catch (error) {
      // ignore logout errors
    }
    clearCookie(COOKIE_TOKEN);
    clearCookie(COOKIE_USER);
    setAuth({ status: "logged_out", token: null, userId: null, user: null });
  };

  const updateActionStatus = (message, tone = "neutral") => {
    setActionStatus({ loading: false, message, tone });
    setTimeout(() => {
      setActionStatus({ loading: false, message: "", tone: "" });
    }, 4500);
  };

  const runAction = async (action, successMessage) => {
    setActionStatus({ loading: true, message: "", tone: "" });
    try {
      const result = await action();
      updateActionStatus(successMessage, "success");
      await refreshOverview();
      return result;
    } catch (error) {
      updateActionStatus(error.message || "Action failed", "danger");
      return null;
    }
  };

  const createContact = async (event) => {
    event.preventDefault();
    await runAction(
      () =>
        apiRequest("/v1/contacts", {
          method: "POST",
          token: authHeader,
          body: contactForm
        }),
      "Contact created"
    );
    setContactForm((prev) => ({
      ...defaultContact,
      owner_id: prev.owner_id
    }));
  };

  const createDeal = async (event) => {
    event.preventDefault();
    const payload = {
      ...dealForm,
      value: Number(dealForm.value || 0)
    };
    await runAction(
      () =>
        apiRequest("/v1/deals", {
          method: "POST",
          token: authHeader,
          body: payload
        }),
      "Deal created"
    );
    setDealForm((prev) => ({
      ...defaultDeal,
      owner_id: prev.owner_id
    }));
  };

  const createTask = async (event) => {
    event.preventDefault();
    if (!taskForm.deal_id) {
      updateActionStatus("Select a deal ID first", "danger");
      return;
    }
    await runAction(
      () =>
        apiRequest(`/v1/deals/${taskForm.deal_id}/tasks`, {
          method: "POST",
          token: authHeader,
          body: {
            title: taskForm.title,
            due_date: taskForm.due_date,
            is_completed: taskForm.is_completed,
            assigned_to_id: taskForm.assigned_to_id
          }
        }),
      "Task created"
    );
    setTaskForm((prev) => ({
      ...defaultTask,
      assigned_to_id: prev.assigned_to_id
    }));
  };

  const addNote = async (event) => {
    event.preventDefault();
    if (!noteForm.contact_id) {
      updateActionStatus("Select a contact ID first", "danger");
      return;
    }
    await runAction(
      () =>
        apiRequest(`/v1/contacts/${noteForm.contact_id}/notes`, {
          method: "POST",
          token: authHeader,
          body: { content: noteForm.content }
        }),
      "Note added"
    );
    setNoteForm(defaultNote);
  };

  const fetchTasks = async (dealId) => {
    if (!dealId) {
      setActiveDealTasks([]);
      setOverview((prev) => ({ ...prev, tasks: 0 }));
      return;
    }
    await runAction(
      async () => {
        const payload = await apiRequest(`/v1/deals/${dealId}/tasks`, {
          token: authHeader
        });
        setActiveDealTasks(payload?.data || []);
        setOverview((prev) => ({
          ...prev,
          tasks: payload?.total_count ?? 0
        }));
        return payload;
      },
      "Tasks loaded"
    );
  };

  const runInsight = async (action, successBuilder) => {
    setActionStatus({ loading: true, message: "", tone: "" });
    try {
      const payload = await action();
      updateActionStatus(successBuilder(payload), "info");
    } catch (error) {
      updateActionStatus(error.message || "Insight failed", "danger");
    }
  };

  const runLeadScore = async (event) => {
    event.preventDefault();
    await runInsight(
      () =>
        apiRequest("/v1/lead-score", {
          method: "POST",
          token: authHeader,
          body: { lead_id: insightForm.lead_id }
        }),
      (payload) =>
        `Lead score: ${payload?.score ?? "-"} (model ${
          payload?.model_version || "-"
        })`
    );
  };

  const runChurnScore = async (event) => {
    event.preventDefault();
    await runInsight(
      () =>
        apiRequest("/v1/churn-score", {
          method: "POST",
          token: authHeader,
          body: { contact_id: insightForm.contact_id }
        }),
      (payload) =>
        `Churn risk: ${payload?.risk_score ?? "-"} (model ${
          payload?.model_version || "-"
        })`
    );
  };

  const runClv = async (event) => {
    event.preventDefault();
    await runInsight(
      () =>
        apiRequest("/v1/clv", {
          method: "POST",
          token: authHeader,
          body: { contact_id: insightForm.contact_id }
        }),
      (payload) =>
        `CLV: ${payload?.clv ?? "-"} (model ${
          payload?.model_version || "-"
        })`
    );
  };

  const createApiKey = async () => {
    await runAction(
      async () => {
        const payload = await apiRequest("/v1/users/api-keys", {
          method: "POST",
          token: authHeader,
          body: { user_id: auth.userId }
        });
        setApiKeyResponse(payload);
        return payload;
      },
      "API key created"
    );
  };

  if (auth.status === "checking") {
    return (
      <div className="app-shell">
        <div className="loading">Checking session...</div>
      </div>
    );
  }

  if (auth.status !== "logged_in") {
    return (
      <div className="app-shell">
        <header className="hero">
          <div>
            <p className="tag">Atlas CRM</p>
            <h1>Real-time CRM intelligence, now within reach.</h1>
            <p className="subtitle">
              Coordinate contacts, track pipeline momentum, and score leads with
              built-in insights.
            </p>
          </div>
        </header>
        <section className="auth-grid">
          <form className="panel" onSubmit={handleLogin}>
            <h2>Login</h2>
            <p className="muted">
              Use your email or user ID from the backend service.
            </p>
            <label>
              Username
              <input
                type="text"
                placeholder="email or user id"
                value={loginForm.username}
                onChange={(event) =>
                  setLoginForm((prev) => ({
                    ...prev,
                    username: event.target.value
                  }))
                }
                required
              />
            </label>
            <label>
              Password
              <input
                type="password"
                value={loginForm.password}
                onChange={(event) =>
                  setLoginForm((prev) => ({
                    ...prev,
                    password: event.target.value
                  }))
                }
                required
              />
            </label>
            <button className="primary" type="submit">
              Login
            </button>
          </form>

          <form className="panel" onSubmit={handleSignup}>
            <h2>Signup</h2>
            <p className="muted">Create a new internal CRM account.</p>
            <div className="grid-two">
              <label>
                First name
                <input
                  type="text"
                  value={signupForm.first_name}
                  onChange={(event) =>
                    setSignupForm((prev) => ({
                      ...prev,
                      first_name: event.target.value
                    }))
                  }
                  required
                />
              </label>
              <label>
                Last name
                <input
                  type="text"
                  value={signupForm.last_name}
                  onChange={(event) =>
                    setSignupForm((prev) => ({
                      ...prev,
                      last_name: event.target.value
                    }))
                  }
                  required
                />
              </label>
            </div>
            <label>
              Email
              <input
                type="email"
                value={signupForm.email}
                onChange={(event) =>
                  setSignupForm((prev) => ({
                    ...prev,
                    email: event.target.value
                  }))
                }
                required
              />
            </label>
            <label>
              Password
              <input
                type="password"
                value={signupForm.password}
                onChange={(event) =>
                  setSignupForm((prev) => ({
                    ...prev,
                    password: event.target.value
                  }))
                }
                required
              />
            </label>
            <label>
              Role
              <select
                value={signupForm.role}
                onChange={(event) =>
                  setSignupForm((prev) => ({
                    ...prev,
                    role: event.target.value
                  }))
                }
              >
                <option value="user">User</option>
                <option value="admin">Admin</option>
              </select>
            </label>
            <button className="primary" type="submit">
              Create account
            </button>
          </form>
        </section>
        {authError ? <div className="banner danger">{authError}</div> : null}
        <footer className="auth-footer">
          Edge base URL: <span>{EDGE_BASE_URL}</span>
        </footer>
      </div>
    );
  }

  return (
    <div className="app-shell">
      <header className="topbar">
        <div>
          <p className="tag">Atlas CRM</p>
          <h1>Welcome back, {auth.user?.first_name || "Operator"}.</h1>
          <p className="subtitle">
            Manage contacts, deals, and predictive insights from one workspace.
          </p>
        </div>
        <div className="topbar-actions">
          <div className="identity">
            <span>{auth.user?.email}</span>
            <span className="role">{auth.user?.role}</span>
          </div>
          <button className="ghost" onClick={handleLogout}>
            Logout
          </button>
        </div>
      </header>

      {actionStatus.message ? (
        <div className={`banner ${actionStatus.tone}`}>
          {actionStatus.message}
        </div>
      ) : null}

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

      <section className="grid-main">
        <div className="panel">
          <h2>Contacts</h2>
          <p className="muted">
            Recently updated contacts and ownership context.
          </p>
          <div className="table">
            <div className="table-row header">
              <span>Name</span>
              <span>Email</span>
              <span>Owner</span>
            </div>
            {contacts.length === 0 ? (
              <div className="table-row empty">No contacts yet.</div>
            ) : (
              contacts.map((contact) => (
                <div className="table-row" key={contact.id}>
                  <span>{contact.name}</span>
                  <span>{contact.email}</span>
                  <span>{contact.owner_name || contact.owner_id}</span>
                </div>
              ))
            )}
          </div>
        </div>

        <div className="panel">
          <h2>Deals</h2>
          <p className="muted">
            Pipeline momentum and expected close dates.
          </p>
          <div className="table">
            <div className="table-row header">
              <span>Deal</span>
              <span>Stage</span>
              <span>Close</span>
            </div>
            {deals.length === 0 ? (
              <div className="table-row empty">No deals yet.</div>
            ) : (
              deals.map((deal) => (
                <div className="table-row" key={deal.id}>
                  <span>{deal.name}</span>
                  <span>{deal.stage}</span>
                  <span>{deal.expected_close_date || "-"}</span>
                </div>
              ))
            )}
          </div>
        </div>
      </section>

      <section className="grid-actions">
        <div className="panel">
          <h2>New Contact</h2>
          <form onSubmit={createContact} className="form-stack">
            <div className="grid-two">
              <label>
                First name
                <input
                  type="text"
                  value={contactForm.first_name}
                  onChange={(event) =>
                    setContactForm((prev) => ({
                      ...prev,
                      first_name: event.target.value
                    }))
                  }
                  required
                />
              </label>
              <label>
                Last name
                <input
                  type="text"
                  value={contactForm.last_name}
                  onChange={(event) =>
                    setContactForm((prev) => ({
                      ...prev,
                      last_name: event.target.value
                    }))
                  }
                  required
                />
              </label>
            </div>
            <label>
              Email
              <input
                type="email"
                value={contactForm.email}
                onChange={(event) =>
                  setContactForm((prev) => ({
                    ...prev,
                    email: event.target.value
                  }))
                }
                required
              />
            </label>
            <label>
              Phone
              <input
                type="text"
                value={contactForm.phone}
                onChange={(event) =>
                  setContactForm((prev) => ({
                    ...prev,
                    phone: event.target.value
                  }))
                }
                required
              />
            </label>
            <label>
              Owner ID
              <input
                type="text"
                value={contactForm.owner_id}
                onChange={(event) =>
                  setContactForm((prev) => ({
                    ...prev,
                    owner_id: event.target.value
                  }))
                }
              />
            </label>
            <button className="primary" type="submit">
              Create contact
            </button>
          </form>
        </div>

        <div className="panel">
          <h2>New Deal</h2>
          <form onSubmit={createDeal} className="form-stack">
            <label>
              Deal name
              <input
                type="text"
                value={dealForm.name}
                onChange={(event) =>
                  setDealForm((prev) => ({
                    ...prev,
                    name: event.target.value
                  }))
                }
                required
              />
            </label>
            <div className="grid-two">
              <label>
                Stage
                <select
                  value={dealForm.stage}
                  onChange={(event) =>
                    setDealForm((prev) => ({
                      ...prev,
                      stage: event.target.value
                    }))
                  }
                >
                  <option value="prospect">Prospect</option>
                  <option value="qualified">Qualified</option>
                  <option value="proposal">Proposal</option>
                  <option value="negotiation">Negotiation</option>
                  <option value="won">Won</option>
                  <option value="lost">Lost</option>
                </select>
              </label>
              <label>
                Value
                <input
                  type="number"
                  value={dealForm.value}
                  onChange={(event) =>
                    setDealForm((prev) => ({
                      ...prev,
                      value: event.target.value
                    }))
                  }
                  min="0"
                  step="0.01"
                />
              </label>
            </div>
            <label>
              Expected close date
              <input
                type="date"
                value={dealForm.expected_close_date}
                onChange={(event) =>
                  setDealForm((prev) => ({
                    ...prev,
                    expected_close_date: event.target.value
                  }))
                }
              />
            </label>
            <label>
              Contact ID
              <input
                type="text"
                value={dealForm.contact_id}
                onChange={(event) =>
                  setDealForm((prev) => ({
                    ...prev,
                    contact_id: event.target.value
                  }))
                }
                required
              />
            </label>
            <label>
              Owner ID
              <input
                type="text"
                value={dealForm.owner_id}
                onChange={(event) =>
                  setDealForm((prev) => ({
                    ...prev,
                    owner_id: event.target.value
                  }))
                }
              />
            </label>
            <button className="primary" type="submit">
              Create deal
            </button>
          </form>
        </div>

        <div className="panel">
          <h2>Tasks & Notes</h2>
          <form onSubmit={createTask} className="form-stack">
            <label>
              Deal ID (for tasks)
              <input
                type="text"
                value={taskForm.deal_id}
                onChange={(event) =>
                  setTaskForm((prev) => ({
                    ...prev,
                    deal_id: event.target.value
                  }))
                }
              />
            </label>
            <button
              className="ghost"
              type="button"
              onClick={() => fetchTasks(taskForm.deal_id)}
            >
              Load tasks for deal
            </button>
            <label>
              Task title
              <input
                type="text"
                value={taskForm.title}
                onChange={(event) =>
                  setTaskForm((prev) => ({
                    ...prev,
                    title: event.target.value
                  }))
                }
                required
              />
            </label>
            <label>
              Due date
              <input
                type="date"
                value={taskForm.due_date}
                onChange={(event) =>
                  setTaskForm((prev) => ({
                    ...prev,
                    due_date: event.target.value
                  }))
                }
              />
            </label>
            <label>
              Assigned to (user ID)
              <input
                type="text"
                value={taskForm.assigned_to_id}
                onChange={(event) =>
                  setTaskForm((prev) => ({
                    ...prev,
                    assigned_to_id: event.target.value
                  }))
                }
              />
            </label>
            <label className="checkbox">
              <input
                type="checkbox"
                checked={taskForm.is_completed}
                onChange={(event) =>
                  setTaskForm((prev) => ({
                    ...prev,
                    is_completed: event.target.checked
                  }))
                }
              />
              Mark complete
            </label>
            <button className="primary" type="submit">
              Create task
            </button>
          </form>

          <div className="divider" />

          <form onSubmit={addNote} className="form-stack">
            <label>
              Contact ID (for notes)
              <input
                type="text"
                value={noteForm.contact_id}
                onChange={(event) =>
                  setNoteForm((prev) => ({
                    ...prev,
                    contact_id: event.target.value
                  }))
                }
              />
            </label>
            <label>
              Note content
              <textarea
                rows="3"
                value={noteForm.content}
                onChange={(event) =>
                  setNoteForm((prev) => ({
                    ...prev,
                    content: event.target.value
                  }))
                }
                required
              />
            </label>
            <button className="primary" type="submit">
              Add note
            </button>
          </form>
        </div>
      </section>

      <section className="grid-actions">
        <div className="panel">
          <h2>Insight Engine</h2>
          <p className="muted">
            Trigger scoring requests routed through the edge service.
          </p>
          <form onSubmit={runLeadScore} className="form-stack">
            <label>
              Lead ID
              <input
                type="text"
                value={insightForm.lead_id}
                onChange={(event) =>
                  setInsightForm((prev) => ({
                    ...prev,
                    lead_id: event.target.value
                  }))
                }
                required
              />
            </label>
            <button className="primary" type="submit">
              Run lead score
            </button>
          </form>
          <div className="divider" />
          <form onSubmit={runChurnScore} className="form-stack">
            <label>
              Contact ID
              <input
                type="text"
                value={insightForm.contact_id}
                onChange={(event) =>
                  setInsightForm((prev) => ({
                    ...prev,
                    contact_id: event.target.value
                  }))
                }
                required
              />
            </label>
            <button className="primary" type="submit">
              Run churn score
            </button>
          </form>
          <form onSubmit={runClv} className="form-stack">
            <button className="ghost" type="submit">
              Calculate CLV
            </button>
          </form>
        </div>

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
          <div className="divider" />
          <h3>Active tasks</h3>
          <div className="table">
            <div className="table-row header">
              <span>Title</span>
              <span>Due</span>
              <span>Status</span>
            </div>
            {activeDealTasks.length === 0 ? (
              <div className="table-row empty">No tasks loaded.</div>
            ) : (
              activeDealTasks.map((task) => (
                <div className="table-row" key={task.id}>
                  <span>{task.title}</span>
                  <span>{formatTimestamp(task.due_date)}</span>
                  <span>{task.is_completed ? "Done" : "Open"}</span>
                </div>
              ))
            )}
          </div>
        </div>
      </section>

      <footer className="footer">
        <span>Edge URL: {EDGE_BASE_URL}</span>
        <button className="ghost" onClick={refreshOverview}>
          Refresh data
        </button>
      </footer>
    </div>
  );
}

export default App;
