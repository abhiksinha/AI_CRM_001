import React from "react";
import { SearchablePicker } from "../components/SearchablePicker";

export function InsightsPage({
  insightForm,
  setInsightForm,
  runLeadScore,
  runChurnScore,
  runClv,
  contactOptions,
  insightLeadInput,
  setInsightLeadInput,
  insightContactInput,
  setInsightContactInput,
  formatContactOption,
  searchContactOption
}) {
  return (
    <div className="grid-actions">
      <div className="panel">
        <h2>Insight Engine</h2>
        <p className="muted">
          Trigger scoring requests routed through the edge service.
        </p>
        <form onSubmit={runLeadScore} className="form-stack">
          <SearchablePicker
            label="Lead ID"
            placeholder="Click to view contacts, or type to search"
            options={contactOptions}
            selectedValue={insightForm.lead_id}
            inputValue={insightLeadInput}
            onInputValueChange={setInsightLeadInput}
            onSelectValue={(value) =>
              setInsightForm((prev) => ({
                ...prev,
                lead_id: value
              }))
            }
            formatOption={formatContactOption}
            searchOption={searchContactOption}
            required
          />
          <button className="primary" type="submit">
            Run lead score
          </button>
        </form>
        <div className="divider" />
        <form onSubmit={runChurnScore} className="form-stack">
          <SearchablePicker
            label="Contact ID"
            placeholder="Click to view contacts, or type to search"
            options={contactOptions}
            selectedValue={insightForm.contact_id}
            inputValue={insightContactInput}
            onInputValueChange={setInsightContactInput}
            onSelectValue={(value) =>
              setInsightForm((prev) => ({
                ...prev,
                contact_id: value
              }))
            }
            formatOption={formatContactOption}
            searchOption={searchContactOption}
            required
          />
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
    </div>
  );
}
