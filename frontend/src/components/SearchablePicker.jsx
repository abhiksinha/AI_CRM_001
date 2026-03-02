import { useEffect, useMemo, useRef, useState } from "react";

export function SearchablePicker({
  label,
  placeholder,
  options,
  selectedValue,
  inputValue,
  onInputValueChange,
  onSelectValue,
  formatOption,
  searchOption,
  emptyText = "No options found",
  required = false,
  disabled = false
}) {
  const [isOpen, setIsOpen] = useState(false);
  const pickerRef = useRef(null);
  const isLockedSelection = Boolean(selectedValue);

  const filteredOptions = useMemo(() => {
    if (isLockedSelection) return options;
    const query = inputValue.trim().toLowerCase();
    if (!query) return options;
    return options.filter((option) => searchOption(option).toLowerCase().includes(query));
  }, [options, inputValue, searchOption, isLockedSelection]);

  const handleOptionSelect = (option) => {
    onSelectValue(option.id);
    onInputValueChange(formatOption(option));
    setIsOpen(false);
  };

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (!pickerRef.current) return;
      if (!pickerRef.current.contains(event.target)) {
        setIsOpen(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  return (
    <label>
      {label}
      <div className="contact-picker" ref={pickerRef}>
        <input
          type="text"
          className="contact-search"
          placeholder={placeholder}
          required={required}
          disabled={disabled}
          readOnly={isLockedSelection}
          value={inputValue}
          onFocus={() => setIsOpen(true)}
          onChange={(event) => {
            const value = event.target.value;
            onInputValueChange(value);
            onSelectValue("");
            setIsOpen(true);
          }}
        />
        {isOpen && !disabled ? (
          <div className="contact-suggestions">
            {filteredOptions.length === 0 ? (
              <div className="contact-empty">{emptyText}</div>
            ) : (
              filteredOptions.slice(0, 100).map((option) => {
                const optionValue = option.id;
                return (
                  <button
                    key={optionValue}
                    type="button"
                    className={`contact-option ${selectedValue === optionValue ? "selected" : ""}`}
                    onMouseDown={(event) => {
                      event.preventDefault();
                      handleOptionSelect(option);
                    }}
                  >
                    {formatOption(option)}
                  </button>
                );
              })
            )}
          </div>
        ) : null}
      </div>
    </label>
  );
}
