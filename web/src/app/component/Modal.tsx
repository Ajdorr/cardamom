import { ImageButton } from "./input";
import "./modal.css"

export interface ModalProps {
  children?: React.ReactNode;
  onClose: () => void
}

export default function ModalPanel(props: ModalProps) {

  const boundaryMargin = 15
  const maxContentHeight = window.innerHeight - (2 * boundaryMargin) - 40 // 40 is the line height

  return (<div className="modal-panel-root">
    <div style={{ margin: `${boundaryMargin}px` }} className="modal-panel-boundary">
      <div className="modal-panel-workspace theme-focus">
        <div className="modal-panel-close theme-primary">
          <ImageButton alt="Close" src="/icons/close.svg" className="modal-panel-close-btn"
            onClick={e => { props.onClose() }} />
        </div>
        <div style={{ maxHeight: `${maxContentHeight}px` }} className="modal-panel-content">
          {props.children}
        </div>
      </div>
    </div>
  </div>)

}