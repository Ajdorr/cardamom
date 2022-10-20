import { useState } from "react"
import { ImageButton, InputTextBox } from "../../component/input"

type ModifiableDropDownProps = {
  options: string[]
  value: string
  id?: string
  className?: string
  placeholder?: string
  dropDownButtonOnLeft?: boolean
  displayClear?: boolean
  onChange: (s: string) => void
}

export function ModifiableDropDown(props: ModifiableDropDownProps) {
  const [value, setValue] = useState(props.value)
  const [isVisible, setVisible] = useState(false)
  const displayList = isVisible && props.options.length > 0

  function save(s: string) {
    const newValue = s.toLowerCase()
    setValue(newValue)
    props.onChange(newValue)
  }

  const clazz = props.className ? "modifiable-drop-down-root " + props.className : "modifiable-drop-down-root"
  const defaultPlaceholder = props.options.length > 0 ? "Add or select store" : "Add store"
  return (
    <div id={props.id} onMouseLeave={e => setVisible(false)} className={clazz}>
      <div className="modifiable-drop-down-workspace" style={{
        flexDirection: props.dropDownButtonOnLeft ? "row-reverse" : "row"
      }}>
        <InputTextBox value={value} className="modifiable-drop-down-input"
          placeholder={props.placeholder ? props.placeholder : defaultPlaceholder}
          onChange={s => save(s)} />

        <ImageButton alt="expand" src="/icons/drop-down.svg"
          style={{ visibility: props.options.length > 0 ? "visible" : "hidden" }}
          className="modifiable-drop-down-list"
          onClick={e => setVisible(!isVisible)} />

        {
          props.displayClear ? <ImageButton
            src="/icons/backspace.svg" alt="clear selected store" className="modifiable-drop-down-clear"
            onClick={e => save("")} /> : null
        }

      </div>
      <div style={{ display: displayList ? "grid" : "none" }}
        className="modifiable-drop-down-overlay theme-focus format-font-medium">{
          props.options.map(o => { return (<DropDownElement key={o} value={o} onClick={() => { save(o); setVisible(false) }} />) })
        }</div>
    </div>
  )
}

type DropDownElementProps = {
  value: string,
  onClick: () => void
}

function DropDownElement(props: DropDownElementProps) {
  return (<div onClick={props.onClick}>{props.value}</div>)
}