import { useState } from "react"
import { ImageButton, InputTextBox } from "../../component/input"

type ModifiableDropDownProps = {
  options: string[]
  value: string
  id?: string
  className?: string
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
  return (
    <div id={props.id} onMouseLeave={e => setVisible(false)} className={clazz}>
      <div className="modifiable-drop-down-workspace">
        <InputTextBox value={value} className="modifiable-drop-down-input"
          placeholder="Add or select store"
          onChange={s => save(s)} />

        <ImageButton alt="expand" src="icons/drop-down.svg" className="modifiable-drop-down-list"
          onClick={e => setVisible(!isVisible)} />

        {
          props.displayClear ? <ImageButton
            src="icons/backspace.svg" alt="clear selected store" className="modifiable-drop-down-clear"
            onClick={e => save("")} /> : null
        }

      </div>
      <div style={{ display: displayList ? "block" : "none" }} className="modifiable-drop-down-overlay theme-focus">{
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