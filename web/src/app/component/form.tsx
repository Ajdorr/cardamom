import "./form.css"
import { useState } from "react";
import { InputPasswordBox, InputTextBox } from "./input";

type TextProps = {
  value: string
  label: string
  className?: string
  onChange: (s: string) => void
}

export function FormText(props: TextProps) {
  const clazz = props.className ? "form-text-root " + props.className : "form-text-root"
  return (<div className={clazz}>
    <div className="form-text-label">{props.label}</div>
    <InputTextBox value={props.value} onChange={props.onChange} />
  </div>)
}

type PasswordProps = {
  value: string
  className?: boolean
  label: string
  onChange: (s: string) => void
}

export function FormPassword(props: PasswordProps) {
  const clazz = props.className ? "form-password-root " + props.className : "form-password-root"
  return (<div className={clazz}>
    <div className="form-password-label">{props.label}</div>
    <InputPasswordBox value={props.value} onChange={props.onChange} />
  </div>)
}

type TextAreaProps = {
  value: string
  label: string
  rows?: number
  className?: string
  defaultValue?: string
  onChange: (s: string) => void
}

export function FormTextArea(props: TextAreaProps) {
  const [isFocused, setFocus] = useState(false)
  const [value, setValue] = useState(props.value)
  const save = function (newValue: string) {
    setValue(newValue)
    props.onChange(newValue)
  }

  const clazz = props.className ? "form-textarea-root " + props.className : "form-textarea-root"
  return (<div className={clazz}>
    <div className="form-textarea-label">{props.label}</div>
    <textarea className="form-textarea-input"
      placeholder={props.label}
      value={isFocused ? value : props.value}
      defaultValue={props.defaultValue} rows={props.rows} style={{ resize: "none" }}
      onChange={e => { setValue(e.target.value) }}
      onFocus={e => { setFocus(true); setValue(props.value) }}
      onBlur={e => { setFocus(false); save(value) }}
    />
  </div>
  )
}

type DropDownProps = {
  value: string
  label: string
  className?: string
  options: Map<string, string>
  onChange: (s: string) => void
}
export function FormDropDown(props: DropDownProps) {
  const optArr = Array.from(props.options.entries())

  const clazz = props.className ? "form-dropdown-root " + props.className : "form-dropdown-root"
  return (<div className={clazz}>
    <div className="form-dropdown-label">{props.label}</div>
    <select className="form-dropdown-select" value={props.value} onChange={e => props.onChange(e.target.value)}>
      {optArr.map(([k, v]) => <option key={k} value={k}>{v}</option>)}
    </select>
  </div>)
}