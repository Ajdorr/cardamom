import "./widget.css"
import { useState } from "react";
import { ImageButton, InputNumberBox } from "./input"

type IntegerInputWidgetProps = {
  value: number
  minValue: number
  onChange: (v: number) => void
}

export function IntegerInputerWidget(props: IntegerInputWidgetProps) {
  const [value, setValue] = useState(props.value)
  const onUpdate = function (v: number) {
    setValue(v)
    props.onChange(v)
  }

  if (value <= props.minValue) {
    return (<ImageButton alt="add" src="/icons/plus.svg"
      onClick={e => onUpdate(props.minValue + 1)} />)
  }

  return (<span className="widget-interget-input-root">
    <ImageButton alt="add" src="/icons/plus.svg" onClick={e => onUpdate(value + 1)} />
    <InputNumberBox value={props.value} onChange={props.onChange} />
    <ImageButton alt="subtract" src="/icons/minus.svg" onClick={e => onUpdate(value - 1)} />
  </span>)
}