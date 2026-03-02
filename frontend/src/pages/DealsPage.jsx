import React from "react";
import { SearchablePicker } from "../components/SearchablePicker";

export function DealsPage({
  deals,
  dealForm,
  setDealForm,
  createDeal,
  contactOptions,
  dealContactInput,
  setDealContactInput,
  formatContactOption,
  searchContactOption
}) {
  return (
    <div className="grid-main">
      <div className="panel">
        <h2>Deals</h2>
        <p className="muted">Pipeline momentum and expected close dates.</p>
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
          <SearchablePicker
            label="Contact ID"
            placeholder="Click to view contacts, or type to search"
            options={contactOptions}
            selectedValue={dealForm.contact_id}
            inputValue={dealContactInput}
            onInputValueChange={setDealContactInput}
            onSelectValue={(value) =>
              setDealForm((prev) => ({
                ...prev,
                contact_id: value
              }))
            }
            formatOption={formatContactOption}
            searchOption={searchContactOption}
          />
          <button className="primary" type="submit">
            Create deal
          </button>
        </form>
      </div>
    </div>
  );
}
