import { KeyboardEvent } from "react"

export function SaveOnEnter(e: KeyboardEvent, onSave: () => void) {
  if (e.key === "Enter") {
    e.preventDefault()
    onSave()
  }
}
