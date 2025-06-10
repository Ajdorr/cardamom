import "./input.css"
import { Component, MouseEventHandler, useState } from "react";
import { Theme } from "@core/theme";
import { SaveOnEnter } from "./util";


type NumberProps = {
  value: number
  placeholder?: string
  id?: string
  inputAttrs?: React.HtmlHTMLAttributes<HTMLInputElement>
  className?: string
  clearOnChange?: boolean
  onChange: (s: number) => void
}

export function InputNumberBox(props: NumberProps) {
  const clazz = props.className ? "input-number-box " + props.className : "input-number-box"
  return <InputBox value={props.value.toString()} className={clazz} id={props.id}
    inputType="number" inputAttrs={props.inputAttrs}
    clearOnChange={props.clearOnChange}
    placeholder={props.placeholder}
    onChange={s => props.onChange(parseInt(s))}
  />
}

type TextProps = {
  value: string
  placeholder?: string
  id?: string
  inputAttrs?: React.HtmlHTMLAttributes<HTMLInputElement>
  className?: string
  clearOnChange?: boolean
  onChange: (s: string) => void
}

export function InputTextBox(props: TextProps) {
  const clazz = props.className ? "input-text-box " + props.className : "input-text-box"
  return <InputBox value={props.value} className={clazz} id={props.id}
    inputType="text" inputAttrs={props.inputAttrs}
    clearOnChange={props.clearOnChange}
    placeholder={props.placeholder}
    onChange={props.onChange}
  />
}

type PasswordProps = {
  value: string
  placeholder?: string
  id?: string
  inputAttrs?: React.HtmlHTMLAttributes<HTMLInputElement>
  className?: string
  clearOnChange?: boolean
  onChange: (s: string) => void
}

export function InputPasswordBox(props: PasswordProps) {
  const clazz = props.className ? "input-password-box " + props.className : "input-password-box"
  return <InputBox value={props.value} className={clazz} id={props.id}
    inputType="password" inputAttrs={props.inputAttrs}
    clearOnChange={props.clearOnChange}
    placeholder={props.placeholder}
    onChange={props.onChange}
  />
}

type InputBoxProps = {
  value: string
  inputType: string
  placeholder?: string
  id?: string
  inputAttrs?: React.HtmlHTMLAttributes<HTMLInputElement>
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

  return (<div id={props.id} className={props.className}>
    <input className="input-box"
      type={props.inputType} {...props.inputAttrs}
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
  id?: string
  className?: string
  label: string
  theme: Theme
  onClick: MouseEventHandler
}

export class TextButton extends Component<TextButtonProps> {

  render() {
    const clazz = this.props.className ? "input-text-button " + this.props.className : "input-text-button"

    return (<input type="button"
      id={this.props.id}
      value={this.props.label}
      onClick={this.props.onClick}
      className={[clazz, this.props.theme.valueOf()].join(" ")}
    />)
  }

}

type ImageButtonProps = {
  id?: string
  src: string
  alt: string
  disabled?: boolean
  className?: string
  style?: React.CSSProperties
  onClick: MouseEventHandler
}

export function ImageButton(props: ImageButtonProps) {

  let rootClasses = ["input-img-button"]
  if (props.disabled) rootClasses.push("input-disabled")
  if (props.className) rootClasses.push(props.className)

  return (<div id={props.id} style={props.style} className={rootClasses.join(" ")}>
    <img src={props.src} alt={props.alt} style={{ cursor: "pointer", opacity: props.disabled ? "0.4" : "" }}
      onClick={props.onClick} width="100%" height="100%" />
  </div>)

}
