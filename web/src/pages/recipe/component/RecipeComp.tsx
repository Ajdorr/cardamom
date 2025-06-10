import "./component.css"
import { useEffect, useRef, useState } from "react"
import { InputTextBox } from "../../../component/input"
import { IngredientModel, ModifierDividerRegex, Units, UpdateIngredient } from "../schema"

const DragToDeleteTolerance = 60

type IngredientProps = {
  className?: string
  model: IngredientModel
  isDraggable: boolean
  isInInventory: boolean
  placeholder?: string
  onChange: (v: UpdateIngredient) => void
  onReorderMove: (s: number) => void
  onReorderComplete: (s: number) => void
  onDelete: () => void
}

export function RecipeIngredient(props: IngredientProps) {

  const [initX, setInitX] = useState(0)
  const [initY, setInitY] = useState(0)
  const [deltaX, setDeltaX] = useState(0)
  const [deltaY, setDeltaY] = useState(0)

  const root = useRef<HTMLDivElement>(null)

  const [quantityError, setQuantityError] = useState(false)
  const quantity = (typeof props.model.quantity === "number") ? String(props.model.quantity) : props.model.quantity
  const quantityClass = quantityError ? "recipe-ingredient-quantity recipe-ingredient-quantity-error" :
    "recipe-ingredient-quantity"

  const getIndexDelta = (dY: number): number => {
    if (root.current) {
      const idxDelta = Math.floor(dY / root.current.clientHeight)
      return (idxDelta < 0) ? idxDelta + 1 : idxDelta
    } else {
      return 0
    }
  }

  let rootStyle = undefined
  if (deltaX !== 0) {
    rootStyle = { transform: `translateX(${deltaX}px)` }
  } else if (deltaY !== 0) {
    rootStyle = { transform: `translateY(${deltaY}px)`, opacity: 0.6 }
  }

  let rootClasses = ["recipe-ingredient-root",
    props.isInInventory ? "recipe-ingredient-root-own" : "recipe-ingredient-root-lack"]
  if (props.className) {
    rootClasses.push(props.className)
  }

  return (<div ref={root} className={rootClasses.join(" ")} style={rootStyle}>
    {deltaX > 0 ?
      <div style={{ width: `${deltaX}px`, transform: `translateX(-${deltaX}px)` }} className="recipe-ingredient-delete-indicator">
        {deltaX > 40 ? <img alt="delete indicator" src="/icons/delete.svg" /> : null}
      </div> : null
    }
    <span className="recipe-ingredient-marker"
      onTouchStart={e => {
        if (!props.isDraggable) {
          return
        }

        setInitX(e.touches[0].clientX); setInitY(e.touches[0].clientY);
        if (document.activeElement instanceof HTMLElement) {
          document.activeElement.blur()
        }
      }}
      onTouchMove={e => {
        if (!props.isDraggable) {
          return
        }

        const dY = e.touches[0].clientY - initY;
        const dX = e.touches[0].clientX - initX;

        if (Math.abs(dY) > Math.abs(dX)) {
          setDeltaX(0); setDeltaY(dY);
          props.onReorderMove(getIndexDelta(dY));
        } else if (dX > 0) {
          setDeltaX(dX); setDeltaY(0);
        }
      }}
      onTouchEnd={e => {
        if (!props.isDraggable) {
          return
        }

        if (Math.abs(deltaY) > Math.abs(deltaX)) {
          props.onReorderComplete(getIndexDelta(deltaY));
        } else if (Math.abs(deltaX) > DragToDeleteTolerance) {
          props.onDelete()
        }
        setInitX(0); setDeltaX(0);
        setInitY(0); setDeltaY(0);
      }}
    >
      <img alt="draggable" src="/icons/drag-indicator.svg" style={{ visibility: props.isDraggable ? "visible" : "hidden" }} />
    </span>

    <InputTextBox value={quantity} className={quantityClass} placeholder={props.placeholder}
      onChange={s => {
        if (s.match(/^\d+$/) || s.match(/^\d+\/\d+$/) || s.match(/^\d+\s+\d+\/\d+$/)) {
          props.onChange({ quantity: s })
          setQuantityError(false)
        } else {
          setQuantityError(true)
        }
      }} />

    <select className="recipe-ingredient-unit" value={props.model.unit ? props.model.unit : ""}
      onChange={e => props.onChange({ unit: e.target.value })}>{
        Units.map((u, i) => { return (<option key={i} value={u}>{u.length > 0 ? u : "none"}</option>) })
      }</select>

    <InputTextBox className="recipe-ingredient-item" placeholder={props.placeholder}
      value={props.model.modifier ? props.model.item + ", " + props.model.modifier : props.model.item}
      inputAttrs={{ autoCapitalize: "none" }} onChange={e => {
        const itemAndMod = e.split(ModifierDividerRegex, 2)
        if (itemAndMod.length === 1) {
          props.onChange({ item: e.trim(), modifier: "" })
        } else {
          props.onChange({ item: itemAndMod[0].trim(), modifier: itemAndMod[1].trim() })
        }
      }} />

    <SingleCheckbox label="optional" className="recipe-ingredient-optional"
      value={props.model.optional} onChange={b => { props.onChange({ optional: b }) }} />
  </div>
  )
}

type InstructionProps = {
  value: string
  onChange: (s: string) => void
}

export function RecipeInstruction(props: InstructionProps) {

  const onChangeTimer = useRef<number | null>(null)
  const inputElement = useRef<HTMLTextAreaElement>(null)
  const [isEditMode, setEditMode] = useState(false)
  const [value, setValue] = useState("")
  const displayValue = isEditMode ? value : props.value

  useEffect(() => {
    if (isEditMode) {
      if (inputElement.current) { inputElement.current.focus() }
    }
  }, [isEditMode])

  const checkClearTimer = () => { if (onChangeTimer.current) window.clearTimeout(onChangeTimer.current) }
  useEffect(() => { return checkClearTimer }, [])
  const changeTimer = () => { if (props.value !== value) { props.onChange(value) } }

  return (<div className="recipe-instruction-root">
    <ol style={{ visibility: isEditMode ? "hidden" : "visible" }} className="recipe-instruction-list"
      onClick={e => { setEditMode(true); setValue(props.value) }}>{
        displayValue.split("\n").map((v, i) => { return (<li key={i}>{v}</li>) })
      }</ol>
    <div style={{ visibility: isEditMode ? "visible" : "hidden" }} className="recipe-instruction-list" >
      <textarea ref={inputElement} value={value}
        onBlur={e => {
          if (props.value !== value) {
            checkClearTimer();
            props.onChange(e.target.value.trim());
          }
          setEditMode(false)
        }}
        onChange={e => {
          checkClearTimer();
          setValue(e.target.value);
          onChangeTimer.current = window.setTimeout(changeTimer, 5000)
        }} />
    </div>
  </div>
  )
}

type SingleCheckboxProps = {
  label: string
  value: boolean
  className?: string
  onChange: (s: boolean) => void
}

function SingleCheckbox(props: SingleCheckboxProps) {

  const rootClass = props.className ? "input-single-checkbox-root " + props.className : "input-single-checkbox-root"
  return (
    <div className={rootClass}>
      <input type="checkbox" className="input-single-checkbox-input" checked={props.value}
        onChange={e => { props.onChange(e.target.checked) }} />
      <span className="input-single-checkbox-label format-font-subscript">{props.label}</span>
    </div>
  )
}