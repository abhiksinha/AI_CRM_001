import React from "react";
import { SearchablePicker } from "../components/SearchablePicker";

export function TasksPage({
  taskForm,
  setTaskForm,
  createTask,
  noteForm,
  setNoteForm,
  addNote,
  contactOptions,
  noteContactInput,
  setNoteContactInput,
  taskDealOptions,
  taskDealInput,
  setTaskDealInput,
  fetchTasks,
  activeDealTasks,
  formatContactOption,
  searchContactOption,
  formatDealOption,
  searchDealOption,
  formatTimestamp
}) {
  return (
    <div className="grid-main">
      <div className="panel">
        <h2>Tasks & Notes</h2>
        <form onSubmit={createTask} className="form-stack">
          <SearchablePicker
            label="Contact ID (for tasks & notes)"
            placeholder="Click to view contacts, or type to search"
            options={contactOptions}
            selectedValue={noteForm.contact_id}
            inputValue={noteContactInput}
            onInputValueChange={setNoteContactInput}
            onSelectValue={(value) => {
              setNoteForm((prev) => ({
                ...prev,
                contact_id: value
              }));
              setTaskForm((prev) => ({
                ...prev,
                deal_id: ""
              }));
              setTaskDealInput("");
            }}
            formatOption={formatContactOption}
            searchOption={searchContactOption}
          />
          <SearchablePicker
            label="Deal ID (for tasks)"
            placeholder={
              noteForm.contact_id
                ? "Click to view deals, or type to search"
                : "Select contact first"
            }
            options={taskDealOptions}
            selectedValue={taskForm.deal_id}
            inputValue={taskDealInput}
            onInputValueChange={setTaskDealInput}
            onSelectValue={(value) =>
              setTaskForm((prev) => ({
                ...prev,
                deal_id: value
              }))
            }
            formatOption={formatDealOption}
            searchOption={searchDealOption}
            emptyText={
              noteForm.contact_id
                ? "No deals found for this contact"
                : "Select contact first"
            }
            disabled={!noteForm.contact_id}
          />
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

      <div className="panel">
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
    </div>
  );
}
