import "./store-drop-down.css"
import { useState } from "react"
import { InputTextBox } from "../../../component/input"

type StoreDropDownProps = {
  options: string[]
  value: string
  id?: string
  className?: string
  placeholder?: string
  dropDownButtonOnLeft?: boolean
  onChange: (s: string | null) => void
}

export function StoreDropDown(props: StoreDropDownProps) {

  const [value, setValue] = useState(props.value)

  function update(s: string) {
    const newValue = s.toLowerCase()
    setValue(newValue)
    props.onChange(newValue || null)
  }

  const rootClass = props.className ? "store-drop-down-root " + props.className : "store-drop-down-root"
  const defaultPlaceholder = props.options.length > 1 ? "Add or select store" : "Add store"
  return (
    <div id={props.id} className={rootClass}>
      <div className="store-drop-down-workspace" style={{
        flexDirection: props.dropDownButtonOnLeft ? "row-reverse" : "row"
      }}>

        <InputTextBox value={value} className="store-drop-down-input"
          placeholder={props.placeholder ? props.placeholder : defaultPlaceholder}
          onChange={s => update(s)} />

        <select data-testid="store-drop-down-select" className="store-drop-down-select" value={props.value} onChange={e => { update(e.target.value) }}>
          {props.options.map((opt, i) => <option key={i} value={opt}>{opt}</option>)}
        </select>
      </div>
    </div>
  )
}
