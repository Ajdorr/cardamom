import { useRef, useState } from "react"
import { ImageButton, InputTextBox } from "../../component/input"
import { Units } from "../schema"

type IngredientProps = {
  className?: string
  quantity: number | string
  unit: string
  value: string
  placeholder?: string
  onQuantityChange: (s: string) => void
  onUnitChange: (s: string) => void
  onValueChange: (s: string) => void
  onMove: (s: number) => void
  onReorder: (s: number) => void
  onDelete?: () => void
}

export function RecipeIngredient(props: IngredientProps) {

  const [initY, setInitY] = useState(0)
  const [deltaY, setDeltaY] = useState(0)
  const root = useRef<HTMLDivElement>(null)

  const [quantityError, setQuantityError] = useState(false)
  const quantity = (typeof props.quantity === "number") ? String(props.quantity) : props.quantity
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

  const cssStyle = (deltaY !== 0) ? {
    transform: `translateY(${deltaY}px)`,
    opacity: 0.6
  } : undefined
  const rootClass = (props.className) ? "recipe-ingredient-root " + props.className : "recipe-ingredient-root"
  return (<div className={rootClass} style={cssStyle} ref={root} >
    <span className="recipe-ingredient-marker"
      onTouchStart={e => { setInitY(e.touches[0].clientY) }}
      onTouchMove={e => {
        e.preventDefault()
        const dY = e.touches[0].clientY - initY;
        setDeltaY(dY);
        props.onMove(getIndexDelta(dY));
      }}
      onTouchEnd={e => { props.onReorder(getIndexDelta(deltaY)); setDeltaY(0); setInitY(0); }} >
      <img alt="draggable" src="/icons/drag-indicator.svg" />
    </span>

    <InputTextBox value={quantity} className={quantityClass} placeholder={props.placeholder}
      onChange={s => {
        if (s.match(/^\d+$/) || s.match(/^\d+\/\d+$/) || s.match(/^\d+\s+\d+\/\d+$/)) {
          props.onQuantityChange(s)
          setQuantityError(false)
        } else {
          setQuantityError(true)
        }
      }} />

    <select className="recipe-ingredient-unit" value={props.unit} onChange={e => props.onUnitChange(e.target.value)}>{
      Units.map((u, i) => { return (<option key={i} value={u}>{u.length > 0 ? u : "none"}</option>) })
    }</select>

    <InputTextBox value={props.value} className="recipe-ingredient-item" placeholder={props.placeholder}
      onChange={props.onValueChange} />

    {!props.onDelete ? null :
      <ImageButton alt="Delete ingredient" src="/icons/delete.svg"
        className="recipe-ingredient-delete" onClick={props.onDelete} />}
  </div>
  )
}

type InstructionProps = {
  order: string
  value: string
  placeholder?: string
  clearOnChange: boolean
  onChange: (s: string) => void
  onDelete?: () => void
}

export function RecipeInstruction(props: InstructionProps) {

  const save = function (newValue: string) {
    if (newValue.length > 0 && !/^\s+$/.test(newValue)) {
      props.onChange(newValue)
    }
  }

  return (<div className="recipe-instruction-root">
    <span className="recipe-instruction-marker">{props.order}</span>
    <InputTextBox
      value={props.value} className="recipe-instruction-input" placeholder={props.placeholder}
      clearOnChange={props.clearOnChange} onChange={s => { save(s) }}
    />
    {!props.onDelete ? null :
      <ImageButton alt="Delete instruction" src="/icons/delete.svg"
        className="recipe-instruction-delete" onClick={props.onDelete} />}
  </div>
  )
}
