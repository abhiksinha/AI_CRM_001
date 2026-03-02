import React from "react";

export function ContactsPage({ contacts, contactForm, setContactForm, createContact }) {
  return (
    <div className="grid-main">
      <div className="panel">
        <h2>Contacts</h2>
        <p className="muted">Recently updated contacts and ownership context.</p>
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
          <button className="primary" type="submit">
            Create contact
          </button>
        </form>
      </div>
    </div>
  );
}
