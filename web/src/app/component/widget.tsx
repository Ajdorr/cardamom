import "./widget.css"
import { useState } from "react";
import { ImageButton, InputNumberBox } from "./input"

type SwipeIndicatorWidgetProps = {
  height: number
  deltaX: number
  iconSrc: string
  className?: string
  isAlignedRight?: boolean
}

export function SwipeIndicatorWidget(props: SwipeIndicatorWidgetProps) {
  const icon = Math.abs(props.deltaX) >= props.height ?
    <img alt="swipe indicator" src={props.iconSrc} style={{
      height: props.height, right: props.isAlignedRight ? "0" : ""
    }} />
    : null

  return (
    <div style={{
      height: props.height,
      width: `${Math.abs(props.deltaX)}px`,
      transform: `translateX(${(-props.deltaX)}px)`,
      right: props.isAlignedRight ? "0" : "",
    }}
      className={props.className ? "widget-slide-indicator " + props.className : "widget-slide-indicator"}>
      {icon}
    </div>
  )
}

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