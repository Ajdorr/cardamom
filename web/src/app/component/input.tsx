import "./input.css"
import { Component, MouseEventHandler, useState } from "react";
import { joinWithClass, Theme } from "../theme";
import { SaveOnEnter } from "./util";


type NumberProps = {
  value: number
  placeholder?: string
  className?: string
  clearOnChange?: boolean
  onChange: (s: number) => void
}

export function InputNumberBox(props: NumberProps) {
  const clazz = props.className ? "input-number-box " + props.className : "input-number-box"
  return <InputBox value={props.value.toString()} className={clazz}
    inputType="number"
    clearOnChange={props.clearOnChange}
    placeholder={props.placeholder}
    onChange={s => props.onChange(parseInt(s))}
  />
}

type TextProps = {
  value: string
  placeholder?: string
  className?: string
  clearOnChange?: boolean
  onChange: (s: string) => void
}

export function InputTextBox(props: TextProps) {
  const clazz = props.className ? "input-text-box " + props.className : "input-text-box"
  return <InputBox value={props.value} className={clazz}
    inputType="text"
    clearOnChange={props.clearOnChange}
    placeholder={props.placeholder}
    onChange={props.onChange}
  />
}

type PasswordProps = {
  value: string
  placeholder?: string
  className?: string
  clearOnChange?: boolean
  onChange: (s: string) => void
}

export function InputPasswordBox(props: PasswordProps) {
  const clazz = props.className ? "input-password-box " + props.className : "input-password-box"
  return <InputBox value={props.value} className={clazz}
    inputType="password"
    clearOnChange={props.clearOnChange}
    placeholder={props.placeholder}
    onChange={props.onChange}
  />
}

type InputBoxProps = {
  value: string
  inputType: string
  placeholder?: string
  className?: string
  clearOnChange?: boolean
  onChange: (s: string) => void
}

function InputBox(props: InputBoxProps) {
  const [isFocused, setFocus] = useState(false)
  const [value, setValue] = useState(props.value)

  const save = function (newValue: string) {
    if (newValue === props.value) {
      return
    }

    if (props.clearOnChange) {
      setValue("")
      props.onChange(newValue)
    } else {
      setValue(newValue)
      props.onChange(newValue)
    }
  }

  return (<div className={props.className}>
    <input className="input-box"
      type={props.inputType}
      placeholder={props.placeholder}
      value={isFocused ? value : props.value}
      onKeyDown={e => SaveOnEnter(e, () => { save(value) })}
      onChange={e => { setValue(e.target.value) }}
      onBlur={e => { setFocus(false); save(value) }}
      onFocus={e => { setFocus(true); setValue(props.value) }}
    />
  </div>
  )
}

type TextButtonProps = {
  label: string
  theme: Theme
  onClick: MouseEventHandler
}

export class TextButton extends Component<TextButtonProps> {

  render() {
    return (<input type="button"
      value={this.props.label}
      onClick={this.props.onClick}
      className={joinWithClass("input-text-button", this.props.theme)}
    />)
  }

}

type ImageButtonProps = {
  src: string
  alt: string
  className?: string
  onClick: MouseEventHandler
}

export function ImageButton(props: ImageButtonProps) {

  return (<div className={props.className ? "input-img-button " + props.className : "input-img-button"}>
    <img src={props.src} alt={props.alt} style={{ cursor: "pointer" }}
      onClick={props.onClick} width="100%" height="100%" />
  </div>)

}
