import { api } from "../core/api"
import { TextButton } from "../component/input"
import { Theme } from "@core/theme"
import { useNavigate } from "react-router-dom"

function Account() {

  const nav = useNavigate()

  return (<div className="auth-account-root theme-surface">
    <div className="auth-account-logout theme-focus">
      <TextButton label="Logout" theme={Theme.Primary} onClick={e => {
        api.post("/auth/logout").then(rsp => {
          localStorage.removeItem("csrf_token")
          nav('auth/login')
        })
      }} />
    </div>
  </div>)
}

export default Account