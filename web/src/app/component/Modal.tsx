import { ImageButton } from "./input";
import "./modal.css"

export interface ModalProps {
  children?: React.ReactNode;
  onClose: () => void
}

export default function ModalPanel(props: ModalProps) {

  return (<div className="modal-panel-root theme-background">
    <div className="modal-panel-workspace theme-focus">
      <div className="modal-panel-close theme-primary">
        <ImageButton alt="Close" src="/icons/close.svg" className="modal-panel-close-btn"
          onClick={e => { props.onClose() }} />
      </div>
      <div className="modal-panel-content">
        {props.children}
      </div>
    </div>
  </div>)

}